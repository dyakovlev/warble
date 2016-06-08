// stores the model of the entire project
// manages saving/loading/synchronizing a project

/* a project file looks like this:
 * {
 *	 id: <project GUID>,
 *   clips: {
 *     <id>: {
 *       id: <unique clip id>,
 *       startTime: <ms from 0 until start position>,
 *       duration: <clip duration in ms>,
 *       name: <clip name>,
 *       meta: <text blob for notes about clip>
 *
 *     }, ... },
 *   channels: [
 *     {
 *       name: <channel name>,
 *       meta: <text blob for notes about channel>
 *     }, ...],
 *   urls: {
 *     loadClip: <url fragment for where to grab clip buffers from>
 *     saveProject: <url fragment for where to send a serialized project file to>
 *     saveClip: <url fragment for where to send a clip buffer to>
 *     deleteClip: <url fragment for removing clips>
 *   },
 *   userConfig: {
 *     <user-settable config>
 *   },
 * }
 */

import { INFO, WARN, ERROR, DEBUG, TRACE } from "./utils/error";
import { ClipStore } from "./components/ClipStore";


export class Project {
	constructor(hub){
		this.dirty = false;
		this.project = null;
		this.emit = hub.emit;

		// check if we need to save once every ~5s
		this.saveTimer = window.setInterval(()=>{ this.save() }, 5000);

		// stores clip buffer arrays. kept separate to simplify project serialization.
		this.buffers = new Map();

		// stores clip metadata in time-indexed order; used for time-based retrieval
		this.clipStore = new ClipStore();

		hub.on('clipupdated', this.updateClip);
		hub.on('clipremoved', this.removeClip);
		hub.on('configupdated', this.updateConfig);
		hub.on('channelsmodified', this.updateChannels);
		hub.on('saveforced', () => { this.save(true); });
	}

	load(){
		this.emit('projectloadstarted');

		// a new, uninitialized pageload will have an empty, new project
		// serialized into it. if it is not, the user tried to access a project
		// they do not have access to, or something else went wrong.
		if (window.__project === undefined) {
			ERROR("no project defined at window.__project");
			this.emit("projectloadfailed", "no project definition found");
			return;
		}

		this.project = window.JSON.parse(window.__project);

		// CAVEAT if this misbehaves it'll be O(insertion sort) which sucks
		this.project.clips.forEach(clip => { this.clipStore.add(clip); });

		for (var clip of this.project.clips) {
			INFO(`downloading clip buffer for clip ${clip.id}`);

			// TODO use some deferred magic to pipeline this
			// TODO investigate using localStorage to cache these instead of storing remotely
			// TODO maybe it makes sense to bundle these somehow in the download?
			$.get(this.project.urls.loadClip.format(clip.id))
				.done(data => { this.buffers.set(clip.id) = data; })
				.fail(err => { ERROR(`failed to load buffer for clip ${clip.id}`); });
		}

		this.emit('loadedproject');
	}

	addClip({clip, buffer}){
		if (this.project.clips.hasOwnProperty(clip.id)){
			ERROR("tried to add a clip that is already stored", clip);
			return;
		}

		this.project.clips[clip.id] = clip;
		this.buffers[clip.id] = buffer;
		this.clipStore.add(clip);
		this.dirty = true;

		// upload to backend
		$.post(this.project.urls.saveClip, {id: clip.id, buffer: buffer})
			.fail((err)=>{ ERROR(`upload of clip ${clip.id} failed`, err) });

		this.emit("addedclip", clip);
	}

	updateClip({clip, buffer}){
		if (!this.project.clips.hasOwnProperty(clip.id)){
			ERROR("attempted to update clip we don't know about", clip)
			return;
		}

		this.project.clips[clip.id] = clip;
		this.buffers[clip.id] = buffer;
		this.clipStore.update(clip);
		this.dirty = true;

		// todo differentiate between updating clip data and buffer?

		this.emit("updatedclip", clip);
	}

	removeClip(clipId){
		if (!this.project.clips.hasOwnProperty(clip.id)){
			ERROR("tried to remove a clip we don't know about", clipId);
			return;
		}

		delete this.project.clips[clipId];
		delete this.buffers[clipId];
		this.clipStore.delete(clipId);
		this.dirty = true;

		$.post(this.project.urls.deleteClip, clipId)
			.fail((err)=>{ ERROR(`delete of clip ${clip.id} failed`, err) });

		this.emit("removedclip", clipId);
	}

	updateChannels(channels){
		this.project.channels = channels;
		this.dirty = true;
	}

	updateConfig(config){
		this.project.config = config;
		this.dirty = true;
	}

	save(force=false){
		if (force || this.dirty){
			$.post(this.project.urls.saveProject, JSON.serialize(this.project))
				.done(this.emit('savedproject'))
				.fail((err)=>{ ERROR("incremental project save failed", err) });
		}
	}
}



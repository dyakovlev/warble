"use strict";

// manage saving/loading/synchronizing a project, as well as action history, redo/undo

/* a project is an object that looks like this:
 * {
 *	 id: <project GUID>,
 *   clips: {
 *     <id>: {
 *       <clip meta>
 *     }, ... },
 *   channels: [
 *     {
 *       <channel meta>
 *     }, ...],
 *   urls: {
 *     loadClip: <url fragment for where to grab clip buffers from>
 *     saveProject: <url fragment for where to send a serialized project file to>
 *     saveClip: <url fragment for where to send a clip buffer to>
 *   },
 *   userConfig: {
 *     <user-settable config>
 *   },
 * }
 *
 *
 */

class ProjectManager {
	constructor(hub){
		this.emit = hub.ns('pm');

		this.project = {};

		hub.on('clipAdded', addClip);
		hub.on('clipUpdated', updateClip);
		hub.on('clipRemoved', removeClip);

		hub.on('configUpdated', updateConfig);
		hub.on('channelsModified', updateChannels);
	}

	load(){
		this.emit('loading');

		if (window._warbleProject === undefined) {
			WARNING("no project file found, starting new project");
			// some sort of new-project stuff here

		} else {
			this.project = window.JSON.parse(window._warbleProject);

			for (var clip of this.project.clips) {
				INFO(`loading clip ${clip.id}`);

				$.get(this.project.urls.loadClip.format(clip.id))
					.done(data => {
						clip.buffer = data;
						this.emit('gotclip', clip);
					});
			}
		}

		this.emit('doneloading');
	}

	updateClips(allClips){
		let clipsToUpload = [], clipsToDelete = [];

		// update self.project with stripped down clips for any modified ones
		for (let clip of allClips){
			if (! this.project.clips.hasOwnProperty(clip.id)){
				clipsToUpload.push(clip);
				this.project.clips[clip.id] = pick(clip, 
			} else if (clip.modified) 

			}
		}

		// find any clips that are no longer 

		$.post(this.project.urls.saveProject, this.project);

		for (let clip of clipsToUpload)
			// TODO pipeline this using Deferred (then)
			$.post(this.project.urls.saveClip, {id: clip.id, buffer: clip.buffer})
				.done()
				.fail();

		for (let clip of clipsToDelete)
			$.post(this.project.urls.softDeleteClip, clip.id)
				.done()
				.fail();
	}

	updateChannels(channels){
	}

	updateConfig(Config){
	}
}

"use strict";

let config = {
	sampleRate: 44100,
};

let w = class Warble {

	constructor(config){
		this.config = config;

		let hub = new EventHub();

		this.project = new ProjectManager(hub);
		this.workspace = new Workspace(hub);

		this.editor = new ClipEditor(hub);
		this.player = new Player(hub);
		this.recorder = new Recorder(hub);

		this.project.load();
	}

	// declare cross-component deps
	// TODO convert this to event-driven stuff
	makeDeps(){
		var eventHandlers = {};

		eventHandlers.recordStart = () => {
			INFO("starting recording clip");
			this.workspace.lock(this.workspace.overlays.RECORDING);
			this.recorder.start();
		};

		eventHandlers.recordStop = () => {
			INFO("stopping recording clip");
			var buffer = this.recorder.stop();
			var channel = this.workspace.getCurrentChannel();
			var start = this.workspace.getCursorTS();
			var newClip = new Clip(buffer, start, channel);

			DEBUG("new clip:", channel, start);

			this.store.add(newClip);
			this.editor.show(newClip);
			this.workspace.unlock(),
		};

		eventHandlers.playAll = () => {
			// TODO
		};

		eventHandlers.pause = () => {
			INFO("toggling playback pause");
			this.player.togglePause();
		};

		eventHandlers.stop = () => {
			INFO("stopping playback");
			this.player.stop();
		};

		eventHandlers.playOne = (clipId) => {
			INFO("playing clip {0}".format(clipId));
			this.player.schedule(this.store.get(clipId));
		};

		eventHandlers.edit = (clipId) => {
			INFO("showing editor for clip {0}".format(clipId));
			this.editor.show(this.store.get(clipId));
		};

		eventHandlers.delete = (clipId) => {
			INFO("deleting clip {0}".format(clipId));
			this.store.delete(clipId);
		};

		eventHandlers.updateClip = (clipId, new_params) => {
			INFO("updating clip {0}".format(clipId));
			DEBUG("clip update parameters:");
			DEBUG(new_params);

			let clip = this.store.get(clipId);
			$.each(new_params, (prop, val) => {clip[prop] = val});
		};

		eventHandlers.getClipFromId = (clipId) => {
			return this.store.get(clipId);
		};

		eventHandlers.drawClipFromId = (clipId) => {
			this.workspace.drawClip(clipId, this.store.get(clipId));
		};

		return eventHandlers;
	}

}(config);

// Warble is the main object and handles wiring together the UI and other components

var Warble = function(){
	this.initComponents();
	this.initUI();
};


Warble.prototype.initComponents = function(){
	// manages Clip storage and metadata management
	this.store = new ClipStore();

	// manages editing of a single clip
	this.editor = new ClipEditor();

	// manages playback of sound clips
	this.player = new Player();

	// manages microphone recording
	this.recorder = new Recorder();

	// darws the clip boxes and channels and time marks and such
	this.workspace = new Workspace();
};


// wire together UI events and app components
Warble.prototype.initUI = function(){
	var eventHandlers = {};

	eventHandlers.recordStart = function(){
		INFO("starting recording clip");
		this.workspace.lock(this.workspace.overlays.RECORDING);
		this.recorder.start();
	}.bind(this);

	eventHandlers.recordStop = function(){
		INFO("stopping recording clip");
		var buffer = this.recorder.stop();
		var channel = this.workspace.getCurrentChannel();
		var start = this.workspace.getCursorTS();
		var newClip = new Clip(buffer, start, channel);

		INFO("adding clip to store, showing editor..");
		this.store.add(newClip);
		this.editor.show(newClip);
		this.workspace.unlock(),
	}.bind(this);

	eventHandlers.playAll = function(){
		// TODO
	}.bind(this);

	eventHandlers.pause = function(){
		INFO("toggling playback pause");
		this.player.togglePause();
	}.bind(this);

	eventHandlers.stop = function(){
		INFO("stopping playback");
		this.player.stop();
	}.bind(this);

	eventHandlers.playOne = function(clipId){
		INFO("playing clip {0}".format(clipId));
		this.player.schedule(this.store.get(clipId));
	}.bind(this);

	eventHandlers.edit = function(clipId){
		INFO("showing editor for clip {0}".format(clipId));
		this.editor.show(this.store.get(clipId));
	}.bind(this);

	eventHandlers.delete = function(clipId){
		INFO("deleting clip {0}".format(clipId));
		this.store.delete(clipId);
	}.bind(this);

	eventHandlers.update = function(clipId, new_params){
		INFO("updating clip {0}".format(clipId));
		DEBUG("clip update parameters:");
		DEBUG(new_params);

		var clip = this.store.get(clipId);
		$.each(new_params, function(prop, val){clip[prop] = val});
	}.bind(this);

	eventHandlers.getClip = function(clipId){
		return this.store.get(clipId);
	}.bind(this);

	eventHandlers.drawClip = function(clipId){
		this.workspace.drawClip(clipId, this.store.get(clipId));
	}.bind(this);

	// workspace and editor share event handlers for common events
	this.workspace.initUI(eventHandlers);
	this.editor.initUI(eventHandlers);
};

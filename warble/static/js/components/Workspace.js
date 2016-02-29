// Workspace generates the visual part of the app

var Workspace = function(){
	this._elements = {
		overlay: $('.workspace .overlay'),
		pane: $('.workspace .pane'),
		controls: $('.workspace .controls'),
		play_control: $('.workspace .control.play'),
		pause_control: $('.workspace .control.pause'),
		stop_control: $('.workspace .control.stop'),
		record_control: $('.workspace .control.record'),
	};

	this._channels = [];
};


// attach UI elements to internal and external event handlers
Workspace.prototype.initUI = function(handlers){
	this._handlers = handlers;

	// record

	// play

	// pause

	// stop

	// delegate for clip actions
}

// templates to overlay when disabling UI for various reasons
Workspace.prototype.overlays = {
	RECORDING: '<span class="recording-overlay">RECORDING...</span>',
	EDITING: '<span class="editing-overlay">EDITING...</span>',
	RENDERING: '<span class="rendering-overlay">RENDERING...</span>',
};

// template for a clip segment
Workspace.prototype._clipTemplate = '<div style="margin-left:{0}px" class="clip"><span class="clip-title">{1}</span><span class="clip-controls"><button class="clip-control play" /><button class="clip-control edit" /></span></div>';


// disable interaction with UI elements
// overlay - HTML template to place in overlay block while UI is locked
Workspace.prototype.lock = function(overlay){
	if (overlay !== undefined) this._elements.overlay.html(overlay);
	this._elements.overlay.show();
	this._elements.controls.disabled = true;
};


// enable interaction with UI elements
Workspace.prototype.unlock = function(){
	this._elements.overlay.hide();
	this._elements.controls.disabled = false;
};


// generate a clip display box and place it in clip's appropriate location
// in the workspace.
Workspace.prototype.drawClip = function(clip){

	var leftOffsetInPX = timeToOffset(clip.start_ts);
	var lengthInPX = timeToLength(clip.length);
	var clipHTML = this._clipTemplate.format(
		leftOffsetInPX,
		lengthInPX,
		clip.channel,
		clip.title,
		clip.id
	);

	this._elements.pane.append($(clipHTML));
};

// place cursor at mouse location
Workspace.prototype.placeCursor = function(){
};

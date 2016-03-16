// Workspace deals with generating the app's views

import {INFO, WARN, ERROR, DEBUG, TRACE} from "./utils/misc";

export class Workspace {
	constructor(hub){
		// TODO this can only be done after domReady
		this.$pane = $('.workspace .pane');
		this.$allChannels = $('.workspace .channel');

		this.overlay = new Overlay(hub);
		this.controls = new Controls(hub);
		this.paneScroller = new PaneScroller(hub);
		this.recorderView = new RecorderView(hub);

		this.emit = hub.emit;
	}

	drawClip(clip){
		var leftOffsetInPX = timeToOffset(clip.start_ts);
		var lengthInPX = timeToLength(clip.length);
		var clipHTML = CLIPTEMPLATE(leftOffsetInPX, lengthInPX, clip.title, clip.id);

		this._elements.pane.append($(clipHTML));
	}
}

// template for a clip segment
const CLIPTEMPLATE = (left, width, title, id) => { `
		<div style="margin-left:${left}px; width:${width}px" class="clip">
			<span class="clip-title">${title}</span>
			<span class="clip-controls">
				<button class="clip-control play" />
				<button class="clip-control edit" />
			</span>
		</div>`
};


/*
 * manages showing and hiding non-interactive overlays
 */
class Overlay {
	constructor(hub){
		this.emit = hub.emit;
		this.$overlay: $('.workspace .overlay');

		// react to project loading behaviors
		hub.on('loadingproject', () => { this.show(OVERLAYS.WAIT) });
		hub.on('loadedproject', () => { this.hide() });
		hub.on('projectloadfailed', (err) => { this.show(OVERLAYS.ERROR, err) });

		// react to recording behaviors
		hub.on('recordpressed', () => { this.show(OVERLAYS.RECORDING) });

		// own UI
		this.$overlay.find('.close').click(this.hide);
	}

	show(overlay, ...args){
		let renderedOverlay = (typeof(overlay) == "function") ? overlay(...args) : overlay;

		this.$overlay.html(renderedOverlay);
		this.$overlay.show();
		this.emit('overlayshown');
	}

	hide(){
		this.$overlay.hide();
		this.emit('overlayhidden');
	}
}

// templates to overlay when disabling UI for various reasons
const OVERLAYS = {
	RECORDING: '<span class="recording-overlay">RECORDING...</span>',
	EDITING: '<span class="editing-overlay">EDITING...</span>',
	RENDERING: '<span class="rendering-overlay">RENDERING...</span>',
	ERROR: (err) => { `<span class="error-overlay">${err}</span>` },
};


/*
 * manages interacting with the overall play/pause/stop/record controls
 */
class Controls {
	constructor(hub){
		this.emit = hub.emit;

		this.$allControls = $('.workspace .controls');

		this.$play_control = $('.workspace .control.play');
		this.$pause_control = $('.workspace .control.pause');
		this.$stop_control = $('.workspace .control.stop');
		this.$record_control = $('.workspace .control.record');

		this.$play_control.click(evt => { this.clickControl(evt, 'playpressed') });
		this.$pause_control.click(evt => { this.clickControl(evt, 'pausepressed') });
		this.$stop_control.click(evt => { this.clickControl(evt, 'stoppressed') });
		this.$record_control.click(evt => { this.clickControl(evt, 'recordpressed') });

		hub.on('overlayshown', this.disable);
		hub.on('overlayhidden', this.enable);

		hub.on('loadingproject', this.disable);
		hub.on('loadedproject', this.enable);
	}

	disable(){
		this.$allControls.disabled = true;
		this.emit('controlsdisabled');
	}

	enable(){
		this.$allControls.disabled = false;
		this.emit('controlsenabled');
	}

	clickControl(button, evt){
		button.addClass('active');
		this.$allControls.not(button).removeClass('active');
		this.emit(evt);
	}
}


// deals with scrolling the view pane when playback is going
class PaneScroller {
	constructor(hub){
		this.emit = hub.emit;

		hub.on('playpressed', this.animatePlayback);
		hub.on('pausepressed', this.pausePlayback);
		hub.on('stoppressed', this.resetPlayback);
	}
}


// deals with the recording modal
class RecorderView {
	constructor(hub){
		this.emit = hub.emit;
	}
}

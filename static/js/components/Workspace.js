// Workspace deals with generating the app's views

import { INFO, WARN, ERROR, DEBUG, TRACE } from "../util/error";
import { draggable } from "../util/misc";
import { Recorder } from "../util/audio";

export class Workspace {
	constructor(hub){
		this.$pane = $('.workspace .pane');

		this.overlay = new Overlay(hub);
		this.controls = new Controls(hub);
		this.paneScroller = new PaneScroller(hub);
		this.recorder = new RecordingModal(hub);
		this.slider = new PlaybackPosition(hub);
	}
}

// manages showing and hiding non-interactive overlays
class Overlay {
	constructor(hub){
		this.emit = hub.emit;

		hub.on('projectloadstarted', () => { this.show(this.OVERLAYS.WAIT) });
		hub.on('loadedproject', this.hide);
		hub.on('projectloadfailed', (err) => { this.show(this.OVERLAYS.ERROR, err) });

		hub.on('recordpressed', () => { this.show(this.OVERLAYS.GREYEDOUT) });

		this.$overlay = $('.workspace .overlay');
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

Overlay.prototype.OVERLAYS = {
	GREYEDOUT: '<span class="greyed-out-overlay"></span>',
	EDITING: '<span class="editing-overlay">EDITING...</span>',
	RENDERING: '<span class="rendering-overlay">RENDERING...</span>',
	ERROR: (err) => { `<span class="error-overlay">${err}</span>` },
};


// manages interacting with the overall play/pause/stop/record controls
class Controls {
	constructor(hub){
		this.emit = hub.emit;

		this.$allControls = $('.workspace .controls');

		this.$allControls.find('.control.play').click(this.makeClickControl('playpressed'));
		this.$allControls.find('.control.pause').click(this.makeClickControl('pausepressed'));
		this.$allControls.find('.control.stop').click(this.makeClickControl('stoppressed'));
		this.$allControls.find('.control.record').click(this.makeClickControl('recordpressed'));

		hub.on('overlayshown', this.disable);
		hub.on('overlayhidden', this.enable);
	}

	disable(){
		this.$allControls.disabled = true;
	}

	enable(){
		this.$allControls.disabled = false;
	}

	makeClickControl(eventToEmit){
		return () => {
			button.addClass('active');
			this.$allControls.not(button).removeClass('active');
			this.emit(eventToEmit);
		};
	}
}


// deals with the physical and temporal position of the current-position slider
class PlaybackPosition{
	constructor(hub){
		this.emit = hub.emit;

		this.$slider = $('.workspace .slider');

		hub.on('playpressed', this.animate);
		hub.on('pausepressed', this.pause);
		hub.on('stoppressed', this.reset);

		draggable(this.$slider, {
			x: true,
			cursor: 'pointer',
			ondrag: () => { this.emit('slidermoved', this.getTime()) },
		});
	}

	getTime(){
	}

	getPos(){
	}

	animate(){
	}

	pause(){
	}

	reset(){
	}

	moveToTime(time){
	}

	moveToPos(px){
	}
}


// shows and sets up interactions with the recording modals
class RecordingModal {
	constructor(hub){
		this.$modal = $('.recording-modal');
		this.$countdown = $('.recording-modal .countdown');
		this.count = 3;
		this.recorder = new Recorder();

		this.prefs = {
			'backingTrack': false, // play recorded tracks while recording clip
		};

		this.$modal.on('changed', 'input.backing-track', this.setPref('backingTrack'));

		this.$modal.on('click', 'button.close', () => {
			this.$modal.hide();
			this.recorder.stop();
		});

		hub.on('recordpressed', this.show);
	}

	setPref(pref){
		return (evt) => {
			this.prefs[pref] = evt.checked;
		};
	}

	show(){
		this.$modal.show();
		// TODO switch 1-second interval to something off the metronome
		this.interval = window.setInterval(() => { this.countdown() }, 1000);
	}

	countdown(timeToGo){
		if (this.count == 0){
			window.clearInterval(this.interval);
			this.count = 3;
			this.recorder.start();
		} else {
			this.$countdown.html(this.count);
		}
	}

	hide(){
		this.$modal.hide();
		this.recorder.stop((buffer) => {
			this.emit('bufferrecorded', buffer);
		});
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



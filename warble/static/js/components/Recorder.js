/*
 * Recorder encapsulates a ScriptProcessorNode to grab and store raw audio samples whenever a
 * microphone input buffer is attached to it. Some housekeeping is done to clear the buffer
 * and prepare the output for consumption by the rest of the app.
 *
 * TODO when it's formalized, use MediaRecorder API instead of ScriptProcessor shim
 * TODO formalize config passing format
 */

import {INFO, WARN, ERROR, DEBUG, TRACE} from "./utils/misc";

// TODO figure out a better range for this
BUFFERLEN = 44100;

export class Recorder{
	constructor(hub){
		this._stream = null;
		this._mic = null;
		this._buffer = [];

		this._AC = new AudioContext();
		this._recorder = this.AC.createScriptProcessorNode(BUFFERLEN, 1, 1);
		this._recorder.onaudioprocess = evt => {
			DEBUG("Recorder received onaudioprocess event"); TRACE(evt);
			this._buffer.push(evt.inputBuffer.getChannelData[0]);
		};

		if(navigator.mediaDevices === undefined) navigator.mediaDevices = {};
		if(navigator.mediaDevices.getUserMedia === undefined){
			WARN("no mediaDevices.getUserMedia found, patching during init");
			navigator.mediaDevices.getUserMedia = generateGUMPatch();
		}

		hub.on('recordstart', this.record);
		hub.on('recordstop', this.stop);

		this.emit = hub.emit;
	}

	record(){
		navigator.getUserMedia({audio: true}, (stream) => {
			this._stream = stream;
			this._mic = this.AC.createMediaStreamSource(stream);
			this._mic.connect(this._recorder);
			this._start_ts = ts();

			INFO("started recording at {0}".format(this._start_ts));
		});
	}

	stop(){
		this._mic.disconnect();
		this._stream.getTracks()[0].stop();

		INFO("stopped recording");

		// return the recorded buffer
		var outBuffer = this._buffer;
		this._buffer = [];

		TRACE(outBuffer);

		emit('gotbuffer', outBuffer);
	}
}

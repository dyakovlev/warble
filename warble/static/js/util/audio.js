import { INFO, WARN, ERROR, DEBUG, TRACE } from "./utils/error";

/* Recorder encapsulates a ScriptProcessorNode to grab and store raw audio samples whenever a
 * microphone input buffer is attached to it.
 *
 * TODO allow use of MediaRecorder API instead of ScriptProcessor shim
 * TODO deal with buffer output/copy more effectively
 * TODO replace Array with typed array once we know what kind of thing is in the buffer
 *
 * usage:
 *   let r = new Recorder();
 *   r.start();
 *	 r.stop((buffer) => { <do things with recorded buffer> });
 */

const BUFFERLEN = 4096;

export class Recorder{
	constructor(){
		this.stream = null;
		this.mic = null;
		this.buffer = []; // TODO replace with byte array
		this.recording = false;

		polyfillGetUserMedia(); // fill in navigator.getUserMedia if it's nonstandard

		this.AC = makeAudioContext();

		this.recorder = this.AC.createScriptProcessorNode(BUFFERLEN, 1, 1);
		this.recorder.onaudioprocess = (evt) => {
			DEBUG("Recorder received onaudioprocess event"); TRACE(evt);
			this.buffer.push(evt.inputBuffer.getChannelData[0]);
		};
	}

	start(){
		if (this.recording) {
			ERROR(`tried to start recording, but recording is already in progress.`);
			return;
		}

		this.recording = true;
		INFO(`started recording at ${ts()}`);

		navigator.getUserMedia({audio: true}, (stream) => {
			this.stream = stream;
			this.mic = this.AC.createMediaStreamSource(stream);
			this.mic.connect(this.recorder);
		});
	}

	stop(cb){
		if (!this.recording) {
			ERROR(`tried to stop recording, but no recording is in progress.`);
			return;
		}

		this.mic.disconnect();
		this.stream.getTracks()[0].stop();

		this.recording = false;
		INFO(`stopped recording at ${ts()}`);

		TRACE(this.buffer);

		cb && cb(this.buffer);

		// clear buffer for future recording
		this.buffer = [];
	}
}


function makeAudioContext(){
	let ac = window.AudioContext || window.webkitAudioContext;
	return new ac();
}


function polyfillGetUserMedia(){
	if(navigator.mediaDevices === undefined) navigator.mediaDevices = {};
	if(navigator.mediaDevices.getUserMedia === undefined){
		WARN("mediaDevices.getUserMedia not found, patching during init");
		navigator.mediaDevices.getUserMedia = function(constraints, onSuccess, onError) {
			// Maybe we have a browser-prefixed one?
			let getUserMedia = (navigator.getUserMedia ||
								navigator.webkitGetUserMedia ||
								navigator.mozGetUserMedia ||
								navigator.msGetUserMedia);

			// Some browsers just don't implement it - return a rejected promise with an error
			if(!getUserMedia) {
				return Promise.reject(new Error('getUserMedia not implemented in this browser'));
			}

			// Otherwise, wrap getUserMedia in a promise to mimic modern API
			return new Promise(function(onSuccess, onError) {
				getUserMedia.call(navigator, constraints, onSuccess, onError);
			});
		}
	}
}

/*
 * Recorder encapsulates a ScriptProcessorNode to grab and store raw audio samples whenever a
 * microphone input buffer is attached to it. Some housekeeping is done to clear the buffer
 * and prepare the output for consumption by the rest of the app.
 *
 * TODO when it's formalized, use MediaRecorder API instead of ScriptProcessor shim
 * TODO formalize config passing format
 */

var Recorder = function(){
	// set up getUserMedia API polyfill
	if(navigator.mediaDevices === undefined) navigator.mediaDevices = {};
	if(navigator.mediaDevices.getUserMedia === undefined){
		INFO("getUserMedia patched during init");
		navigator.mediaDevices.getUserMedia = generateGUMPatch();
	}

	this.isRecording = false;

	this._stream = null;
	this._mic = null;
	this._buffer = [];
	this._AC = new AudioContext();
    this._recorder = this.AC.createScriptProcessorNode(config.bufferLen, 1, 1);

	this._recorder.onaudioprocess = function(evt){
		DEBUG("Recorder received onaudioprocess event");
		TRACE(evt);
		this._buffer.push(evt.inputBuffer.getChannelData[0]);
	}.bind(this);

	INFO("initialized Recorder component");
};

Recorder.prototype.start = function(){
	if (this.isRecording) throw OperationError("already recording, can't start");

	navigator.getUserMedia({audio: true}, function(stream){
		this._stream = stream;
	    this._mic = this.AC.createMediaStreamSource(stream);
	    this._mic.connect(this._recorder);
		this.isRecording = true;
		this._start_ts = ts();

		INFO("started recording at {0}".format(this._start_ts));
	}.bind(this));
};

Recorder.prototype.stop = function(){
	if (!this.isRecording) throw OperationError("not recording, nothing to stop");

	// detach the media stream from the recorder
	this._mic.disconnect();
	this._stream.getTracks()[0].stop();
	this.isRecording = false;

	INFO("stopped recording");

	// return the recorded buffer
	var outBuffer = this._buffer;
	this._buffer = [];

	TRACE(outBuffer);
	return outBuffer;
};

Recorder.prototype.harvest = function(){

};

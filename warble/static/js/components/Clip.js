"use strict";

// Clips encapsulate an audio recording and related metadata
class Clip {

	constructor(
		buffer = [],
		buffer_type = this.DEFAULT_BUFFER_TYPE,
		start_ts = 0,
		channel = this.DEFAULT_CHANNEL
	) {

		// clip's unique ID
		this.id = Symbol();

		// object to store audio
		this.buffer = buffer;

		// type of the buffer (uncompressed WAV, reference to another clip by symbol, draft)
		this.buffer_type = buffer_type;

		// milliseconds between start of project and start of clip
		this.start_ts = start_ts;

		// current channel (index into workspace's channel list)
		this.channel = channel;
	}

};


// put clips into the default channel if it's not specified during construction
// todo: clip channel should always be specified, default shouldn't be needed
Clip.prototype.DEFAULT_CHANNEL = 0;


// the different kinds of buffers a Clip can hold
// - wav is a uncompressed wav format buffer
// - mp3 is a compressed mp3 format buffer
// - ref is a Symbol reference to another clip's buffer
// - nil is an empty placeholder buffer
Clip.prototype.BUFFER_TYPES = {
	wav: Symbol('wav buffer'),
	mp3: Symbol('mp3 buffer'),
	ref: Symbol('reference buffer'),
	nil: Symbol('placeholder buffer'),
};


// todo: clip buffer type should always be specified, default shouldn't be needed
Clip.prototype.DEFAULT_BUFFER_TYPE = Clip.prototype.BUFFER_TYPES['wav'];


// length of clip in ms
Clip.prototype.length = () => { this.buffer[0].length / config.sampleRate };


// Clips encapsulate an audio recording and related metadata

import { INFO, WARN, ERROR, DEBUG, TRACE } from "./utils/error";

export class Clip {
	constructor(
		startTime,
		endTime,
		channel,
		buffer_type = BUFFER_TYPES['wav']
	){
		// clip's unique ID
		// TODO this needs to get converted into something unique on the backend
		this.id = Symbol();

		// type of the buffer (uncompressed WAV, reference to another clip by symbol, draft)
		this.buffer_type = buffer_type;

		// milliseconds between start of project and start of clip
		this.startTime = startTime;

		// milliseconds between start of project and end of clip
		this.endTime = endTime;

		// channel
		this.channel = channel;
	}
}


// the different kinds of buffers a Clip can hold
const BUFFER_TYPES = {
	wav: Symbol('wav buffer'),			// uncompressed wav format buffer
	mp3: Symbol('mp3 buffer'),			// compressed mp3 format buffer
	ref: Symbol('reference buffer'),    // Symbol reference to another clip's buffer
	nil: Symbol('placeholder buffer'),  // empty placeholder buffer
};



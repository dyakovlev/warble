// Clips encapsulate an audio recording and its metadata

var Clip = function(buffer, start_ts, channel){
	// clip's unique ID
	this._id = _id();

	// object to store uncompressed WAV samples, like {left: [...], right: [...]}
	this.buffer = buffer;

	// length of the clip in seconds
	this.length = buffer[0].length / config.sampleRate;

	// seconds between start of project and start of clip
	this.start_ts = start_ts;

	// current channel (index into workspace's channel list)
	this.channel = channel;
};


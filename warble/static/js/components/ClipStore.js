"use strict";

// ClipStore holds all the clips in sorted order to enable efficient search and ordered playback

// TODO: add generators to ClipStore for in-order traversal of clips

var ClipStore = function(){
	this._max_id = 0;

	// id -> Clip
	this._clip_map = {};

	// sorted by start time asc
	this._clip_list = [];

	// all the channels we're aware of (strict ordering)
	this._channels = [];

	INFO("initialized ClipStore component");
};

// stores a clip
ClipStore.prototype.add = function(clip){
	var id = this._id();
	DEBUG("adding clip {0} to datastore".format(id));

	this._clip_map[id] = clip;

	// insert handle into sorted list
};


// returns a clip
ClipStore.prototype.get = function(id){
	DEBUG("retrieving clip {0}".format(id));
	return this._clip_map[id];
};


// returns list of clips that start in provided time range
ClipStore.prototype.getSlice = function(start, end){
	// page to start time
	// step until end time
};


// removes clip
ClipStore.prototype.delete = function(id){
};


// make a clip handle
ClipStore.prototype._id = function(){
	return ++this._max_id;
};


var Player = function(){
	window.AudioContext = window.AudioContext || window.webkitAudioContext;
	this.ac = new AudioContext();
	var isPlaying = false;

	INFO("initialized Player component");
};

// takes array of Clips, schedules according to start offsets, sets up effects, plays all of them
Player.prototype.play = function(list_o_clips){
};


// stops playback
Player.prototype.stop = function(){
};



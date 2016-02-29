/* Misc Utils */

// logging
var _log_level = 1; // TODO move to config
var _log_error = function(lvl){return function(msg){_log_level >= lvl && console.log(msg)}}
var ERROR = _log_error(1)
var WARN = _log_error(2)
var INFO = _log_error(3)
var DEBUG = _log_error(4)
var TRACE = _log_error(5)

// math
function min(a, b){ return (a < b) ? a : b; }
function max(a, b){ return (a > b) ? a : b; }

// errors
function OperationError(msg){
	this.type = "OperationError";
	this.message = msg;
}

// string formatting
if (!String.prototype.format) {
    String.prototype.format = function() {
		var args = arguments;
		return this.replace(/{(\d+)}/g, function(match, number) {
		    return typeof args[number] != 'undefined' ? args[number] : match;
		});
    };
}

var __max_id = 0;
function _id(){
	return ++__max_id;
};

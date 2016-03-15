/* Misc Utils */

// logging
var _log_level = 1; // TODO move to config
var _log_error = function(lvl){
	return function(msg){
		if (_log_level >= lvl)
			for (let i = 0; i < arguments.length; i++)
				console.log(arguments[i]);
	}
};

var ERROR = _log_error(1)
var WARN = _log_error(2)
var INFO = _log_error(3)
var DEBUG = _log_error(4)
var TRACE = _log_error(5)


// math
var clamp = (a, b, c) => Math.min(Math.max(a, b, c))


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

// picking out object fields
function pick(o, ...fields) {
	return Object.assign({}, ...(for (p of fields) {[p]: o[p]}));
}

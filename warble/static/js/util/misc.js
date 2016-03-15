/* Misc Utils */

// logging
const _log_level = 1; // TODO move to config
const _log_error = function(lvl){
	return function(msg){
		if (_log_level >= lvl)
			for (let i = 0; i < arguments.length; i++)
				console.log(arguments[i]);
	}
};

export const ERROR = _log_error(1)
export const WARN = _log_error(2)
export const INFO = _log_error(3)
export const DEBUG = _log_error(4)
export const TRACE = _log_error(5)


// math
var clamp = (a, b, c) => Math.min(Math.max(a, b, c))


// errors
export function OperationError(msg){
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


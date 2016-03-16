// logging
const _log_level = 1; // TODO move to config, maybe find a way to set it per module
const _log_error = function(lvl){
	return (...args) => {
		if (_log_level >= lvl)
			args.forEach((arg) => { console.log(arg) });
	}
};

export const ERROR = _log_error(1)
export const WARN = _log_error(2)
export const INFO = _log_error(3)
export const DEBUG = _log_error(4)
export const TRACE = _log_error(5)


// math
var clamp = (a, b, c) => Math.min(Math.max(a, b), c)


// string formatting
// TODO this needs to be imported in places that use it for the side effect to take effect
// TODO reconsider if we actually need this
if (!String.prototype.format) {
    String.prototype.format = function() {
		var args = arguments;
		return this.replace(/{(\d+)}/g, function(match, number) {
		    return typeof args[number] != 'undefined' ? args[number] : match;
		});
    };
}


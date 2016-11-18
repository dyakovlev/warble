/* Error logging module
 *
 * TODO
 * - set LOG_LEVEL in a centralized config somewhere?
 */

const LOG_LEVEL = 1;

const log_error = function(lvl){
	return (...args) => {
		if (LOG_LEVEL >= lvl)
			args.forEach((arg) => { console.log(arg) });
	}
};

export const ERROR = log_error(1)
export const WARN = log_error(2)
export const INFO = log_error(3)
export const DEBUG = log_error(4)
export const TRACE = log_error(5)


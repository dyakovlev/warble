/*
 * an ES6 event hub
 *
 * usage:
 *	let h = new EventHub()
 *
 *	let handle = h.on('myevent', (args...) => { console.log(...args) });
 *
 *  hub.emit('myevent', 1, 2, 3); // > 1, 2, 3
 *  hub.off('myevent', handle); // callback is removed from myevent's list of callbacks
 *
 *  let handle = hub.once('myevent', () => { console.log('cb called') });
 *  hub.emit('myevent'); // callback is called, then immediately removed
 *  hub.off('myevent', handle); // error, callback is already removed
 *
 */

import { INFO, WARN, ERROR, DEBUG, TRACE } from "./utils/error";

export class EventHub {
	constructor() {
		this._callbacks = new Map();
	}

	on(evt, callback) {
		let cbSymbol = Symbol();

		this._callbacks.has(evt) || this._callbacks.set(evt, new Map());
		this._callbacks.get(evt).set(cbSymbol, callback);

		INFO(`callback added for event "${evt}"`, callback, cbSymbol);
		return cbSymbol;
	}

	once(evt, callback){
		let cbSymbol = Symbol();

		this._callbacks.has(evt) || this._callbacks.set(evt, new Map());
		this._callbacks.get(evt).set(cbSymbol, (...args) => {
			callback(...args);
			this.off(cbSymbol);
		});

		INFO(`one-time callback added for event "${evt}"`, callback, cbSymbol);
		return cbSymbol;
	}

	off(evt, cbSymbol){
		if (!this._callbacks.has(evt)){
			ERROR(`tried to remove callback ${cbSymbol} for event "${svt}" but no callbacks have been stored for that event`);
			return;
		}

		if (!this._callbacks.get(evt).delete(cbSymbol)){
			ERROR(`tried to remove callback ${cbSymbol} for event "${svt}" but could not find it in that event's list of callbacks`);
			return;
		}

		DEBUG(`callback removed for event "${evt}"`, cbSymbol);
	}

	emit(evt, ...args) {
		DEBUG(`"${evt}" event emitted`, data);

		if (!this._callbacks.has(evt)) {
			WARN(`no listeners are listening for "${evt}" events`);
			return;
		}

		this._callbacks.get(evt).forEach(cb => cb(...args));
	}
}

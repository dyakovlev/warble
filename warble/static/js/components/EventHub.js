"use strict";

/* a namespaced event hub
 *
 * usage:
 *   hub = new EventHub()
 *   emit = hub.ns('ns')
 *
 *   hub.emit('loaded') // emits a 'loaded' event
 *   emit('loaded') // emits a 'ns.loaded' event
 */

class EventHub {
	constructor() {
		this._callbacks = {};
	}

	on(evt, callback) {
		if (! this._callbacks.hasOwnProperty(evt)) this._callbacks[evt] = [];
		this._callbacks[evt].push(callback);
	}

	emit(evt, data) {
		DEBUG(`${evt} event emitted`, data);

		// 'ns.evt' listeners should be called before 'evt' ones

		if (this._callbacks.hasOwnProperty(evt)) {
			this._callbacks[evt].map(call, data);
		}

		[ns, shortEvt] = evt.split('.');

		if (this._callbacks.hasOwnProperty(shortEvt) && evt !== shortEvt) {
			this._callbacks[shortEvt].map(call, data);
		}
	}

	ns(name) {
		return (evt, data) => { this.emit(`${name}.${evt}`, data) }
	}
}


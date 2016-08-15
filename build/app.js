(function () {
  'use strict';

  /* Error logging module
   *
   * TODO
   * - set LOG_LEVEL in a centralized config somewhere?
   */

  var LOG_LEVEL = 1;

  var log_error = function log_error(lvl) {
    return function () {
      for (var _len = arguments.length, args = Array(_len), _key = 0; _key < _len; _key++) {
        args[_key] = arguments[_key];
      }

      if (LOG_LEVEL >= lvl) args.forEach(function (arg) {
        console.log(arg);
      });
    };
  };

  var ERROR = log_error(1);
  var WARN = log_error(2);
  var INFO = log_error(3);
  var DEBUG = log_error(4);
  var TRACE = log_error(5);

  var classCallCheck = function (instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  };

  var createClass = function () {
    function defineProperties(target, props) {
      for (var i = 0; i < props.length; i++) {
        var descriptor = props[i];
        descriptor.enumerable = descriptor.enumerable || false;
        descriptor.configurable = true;
        if ("value" in descriptor) descriptor.writable = true;
        Object.defineProperty(target, descriptor.key, descriptor);
      }
    }

    return function (Constructor, protoProps, staticProps) {
      if (protoProps) defineProperties(Constructor.prototype, protoProps);
      if (staticProps) defineProperties(Constructor, staticProps);
      return Constructor;
    };
  }();

  var EventHub = function () {
  	function EventHub() {
  		classCallCheck(this, EventHub);

  		this._callbacks = new Map();
  	}

  	createClass(EventHub, [{
  		key: "on",
  		value: function on(evt, callback) {
  			var cbSymbol = Symbol();

  			this._callbacks.has(evt) || this._callbacks.set(evt, new Map());
  			this._callbacks.get(evt).set(cbSymbol, callback);

  			INFO("callback added for event \"" + evt + "\"", callback, cbSymbol);
  			return cbSymbol;
  		}
  	}, {
  		key: "once",
  		value: function once(evt, callback) {
  			var _this = this;

  			var cbSymbol = Symbol();

  			this._callbacks.has(evt) || this._callbacks.set(evt, new Map());
  			this._callbacks.get(evt).set(cbSymbol, function () {
  				callback.apply(undefined, arguments);
  				_this.off(cbSymbol);
  			});

  			INFO("one-time callback added for event \"" + evt + "\"", callback, cbSymbol);
  			return cbSymbol;
  		}
  	}, {
  		key: "off",
  		value: function off(evt, cbSymbol) {
  			if (!this._callbacks.has(evt)) {
  				ERROR("tried to remove callback " + cbSymbol + " for event \"" + svt + "\" but no callbacks have been stored for that event");
  				return;
  			}

  			if (!this._callbacks.get(evt).delete(cbSymbol)) {
  				ERROR("tried to remove callback " + cbSymbol + " for event \"" + svt + "\" but could not find it in that event's list of callbacks");
  				return;
  			}

  			DEBUG("callback removed for event \"" + evt + "\"", cbSymbol);
  		}
  	}, {
  		key: "emit",
  		value: function emit(evt) {
  			for (var _len = arguments.length, args = Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
  				args[_key - 1] = arguments[_key];
  			}

  			DEBUG("\"" + evt + "\" event emitted", data);

  			if (!this._callbacks.has(evt)) {
  				WARN("no listeners are listening for \"" + evt + "\" events");
  				return;
  			}

  			this._callbacks.get(evt).forEach(function (cb) {
  				return cb.apply(undefined, args);
  			});
  		}
  	}]);
  	return EventHub;
  }();

  var ClipStore = function () {
  	function ClipStore() {
  		classCallCheck(this, ClipStore);

  		this._clipsByStartTime = new Array();
  		this._clipsByEndTime = new Array();
  	}

  	createClass(ClipStore, [{
  		key: "add",
  		value: function add(clip) {
  			INFO("adding clip " + clip.id + " to clip store");

  			var startInsertPos = bisect(this._clipsByStartTime, clip.startTime),
  			    endInsertPos = bisect(this._clipsByEndTime, clip.endTime);

  			this._clipsByStartTime.splice(startInsertPos, 0, [clip.startTime, clip.id]);
  			DEBUG("clip " + clip.id + " added to position " + startInsertPos + " in start-time list");

  			this._clipsByEndTime.splice(endInsertPos, 0, [clip.endTime, clip.id]);
  			DEBUG("clip " + clip.id + " added to position " + endInsertPos + " in end-time list");
  		}
  	}, {
  		key: "delete",
  		value: function _delete(id) {
  			INFO("deleting clip " + id + " from clip store");

  			var startClipDataIndex = this._clipsByStartTime.findIndex(function (item) {
  				item[1].id == id;
  			});
  			if (startClipDataIndex == -1) {
  				WARN("clip " + id + " not found in start-time list");
  			}
  			this._clipsByStartTime.splice(startClipDataIndex, 1);

  			var endClipDataIndex = this._clipsByEndTime.findIndex(function (item) {
  				item[1].id == id;
  			});
  			if (endClipDataIndex == -1) {
  				WARN("clip " + id + " not found in end-time list");
  			}
  			this._clipsByEndTime.splice(endClipDataIndex, 1);
  		}
  	}, {
  		key: "update",
  		value: function update(clip) {
  			INFO("updating clip " + clip.id + " in clip store");

  			this.delete(clip.id);
  			this.add(clip);
  		}
  	}, {
  		key: "get",
  		value: function get(time) {
  			// get all clip IDs that occur at specific time

  			// TODO figure out how this is used
  		}
  	}, {
  		key: "getRange",
  		value: function getRange(start, end) {
  			// get all clip IDs that start after start and end before end

  			// TODO figure out how this is used
  		}
  	}]);
  	return ClipStore;
  }();

  // binary search to find insertion index for pivot that keeps array sorted
  var bisect = function bisect(array, pivot) {
  	var low = 0,
  	    high = array.length,
  	    mid = void 0;

  	while (low < high) {
  		mid = low + high >> 1;
  		pivot < array[mid][0] ? high = mid : low = mid + 1;
  	}

  	return low;
  };

  var Project = function () {
  	function Project(hub) {
  		var _this = this;

  		classCallCheck(this, Project);

  		this.dirty = false;
  		this.project = null;
  		this.emit = hub.emit;

  		// check if we need to save once every ~5s
  		this.saveTimer = window.setInterval(function () {
  			_this.save();
  		}, 5000);

  		// stores clip buffer arrays. kept separate to simplify project serialization.
  		this.buffers = new Map();

  		// stores clip metadata in time-indexed order; used for time-based retrieval
  		this.clipStore = new ClipStore();

  		hub.on('clipupdated', this.updateClip);
  		hub.on('clipremoved', this.removeClip);
  		hub.on('configupdated', this.updateConfig);
  		hub.on('channelsmodified', this.updateChannels);
  		hub.on('saveforced', function () {
  			_this.save(true);
  		});
  	}

  	createClass(Project, [{
  		key: "load",
  		value: function load() {
  			var _this2 = this;

  			this.emit('projectloadstarted');

  			// a new, uninitialized pageload will have an empty, new project
  			// serialized into it. if it is not, the user tried to access a project
  			// they do not have access to, or something else went wrong.
  			if (window.__project === undefined) {
  				ERROR("no project defined at window.__project");
  				this.emit("projectloadfailed", "no project definition found");
  				return;
  			}

  			this.project = window.JSON.parse(window.__project);

  			// CAVEAT if this misbehaves it'll be O(insertion sort) which sucks
  			this.project.clips.forEach(function (clip) {
  				_this2.clipStore.add(clip);
  			});

  			var _iteratorNormalCompletion = true;
  			var _didIteratorError = false;
  			var _iteratorError = undefined;

  			try {
  				for (var _iterator = this.project.clips[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true) {
  					var clip = _step.value;

  					INFO("downloading clip buffer for clip " + clip.id);

  					// TODO use some deferred magic to pipeline this
  					// TODO investigate using localStorage to cache these instead of storing remotely
  					// TODO maybe it makes sense to bundle these somehow in the download?
  					$.get(this.project.urls.loadClip.format(clip.id)).done(function (data) {
  						_this2.buffers.set(clip.id, data);
  					}).fail(function (err) {
  						ERROR("failed to load buffer for clip " + clip.id);
  					});
  				}
  			} catch (err) {
  				_didIteratorError = true;
  				_iteratorError = err;
  			} finally {
  				try {
  					if (!_iteratorNormalCompletion && _iterator.return) {
  						_iterator.return();
  					}
  				} finally {
  					if (_didIteratorError) {
  						throw _iteratorError;
  					}
  				}
  			}

  			this.emit('loadedproject');
  		}
  	}, {
  		key: "addClip",
  		value: function addClip(_ref) {
  			var clip = _ref.clip;
  			var buffer = _ref.buffer;

  			if (this.project.clips.hasOwnProperty(clip.id)) {
  				ERROR("tried to add a clip that is already stored", clip);
  				return;
  			}

  			this.project.clips[clip.id] = clip;
  			this.buffers[clip.id] = buffer;
  			this.clipStore.add(clip);
  			this.dirty = true;

  			// upload to backend
  			$.post(this.project.urls.saveClip, { id: clip.id, buffer: buffer }).fail(function (err) {
  				ERROR("upload of clip " + clip.id + " failed", err);
  			});

  			this.emit("addedclip", clip);
  		}
  	}, {
  		key: "updateClip",
  		value: function updateClip(_ref2) {
  			var clip = _ref2.clip;
  			var buffer = _ref2.buffer;

  			if (!this.project.clips.hasOwnProperty(clip.id)) {
  				ERROR("attempted to update clip we don't know about", clip);
  				return;
  			}

  			this.project.clips[clip.id] = clip;
  			this.buffers[clip.id] = buffer;
  			this.clipStore.update(clip);
  			this.dirty = true;

  			// todo differentiate between updating clip data and buffer?

  			this.emit("updatedclip", clip);
  		}
  	}, {
  		key: "removeClip",
  		value: function removeClip(clipId) {
  			if (!this.project.clips.hasOwnProperty(clip.id)) {
  				ERROR("tried to remove a clip we don't know about", clipId);
  				return;
  			}

  			delete this.project.clips[clipId];
  			delete this.buffers[clipId];
  			this.clipStore.delete(clipId);
  			this.dirty = true;

  			$.post(this.project.urls.deleteClip, clipId).fail(function (err) {
  				ERROR("delete of clip " + clip.id + " failed", err);
  			});

  			this.emit("removedclip", clipId);
  		}
  	}, {
  		key: "updateChannels",
  		value: function updateChannels(channels) {
  			this.project.channels = channels;
  			this.dirty = true;
  		}
  	}, {
  		key: "updateConfig",
  		value: function updateConfig(config) {
  			this.project.config = config;
  			this.dirty = true;
  		}
  	}, {
  		key: "save",
  		value: function save() {
  			var force = arguments.length <= 0 || arguments[0] === undefined ? false : arguments[0];

  			if (force || this.dirty) {
  				$.post(this.project.urls.saveProject, JSON.serialize(this.project)).done(this.emit('savedproject')).fail(function (err) {
  					ERROR("incremental project save failed", err);
  				});
  			}
  		}
  	}]);
  	return Project;
  }();

  // let things be dragged
  function draggable(el, opts) {
  	var defaults = {
  		x: false, // do not allow drag in the x direction
  		y: false, // do not allow drag in the y direction
  		cursor: 'pointer', // change cursor to pointer when hovering over el
  		ondrag: function ondrag() {} };
  }

  /* Recorder encapsulates a ScriptProcessorNode to grab and store raw audio samples whenever a
   * microphone input buffer is attached to it.
   *
   * TODO allow use of MediaRecorder API instead of ScriptProcessor shim
   * TODO deal with buffer output/copy more effectively
   * TODO replace Array with typed array once we know what kind of thing is in the buffer
   *
   * usage:
   *   let r = new Recorder();
   *   r.start();
   *	 r.stop((buffer) => { <do things with recorded buffer> });
   */

  var BUFFERLEN = 4096;

  var Recorder = function () {
  	function Recorder() {
  		var _this = this;

  		classCallCheck(this, Recorder);

  		this.stream = null;
  		this.mic = null;
  		this.buffer = []; // TODO replace with byte array
  		this.recording = false;

  		polyfillGetUserMedia(); // fill in navigator.getUserMedia if it's nonstandard

  		this.AC = makeAudioContext();

  		this.recorder = this.AC.createScriptProcessorNode(BUFFERLEN, 1, 1);
  		this.recorder.onaudioprocess = function (evt) {
  			DEBUG("Recorder received onaudioprocess event");TRACE(evt);
  			_this.buffer.push(evt.inputBuffer.getChannelData[0]);
  		};
  	}

  	createClass(Recorder, [{
  		key: "start",
  		value: function start() {
  			var _this2 = this;

  			if (this.recording) {
  				ERROR("tried to start recording, but recording is already in progress.");
  				return;
  			}

  			this.recording = true;
  			INFO("started recording at " + ts());

  			navigator.getUserMedia({ audio: true }, function (stream) {
  				_this2.stream = stream;
  				_this2.mic = _this2.AC.createMediaStreamSource(stream);
  				_this2.mic.connect(_this2.recorder);
  			});
  		}
  	}, {
  		key: "stop",
  		value: function stop(cb) {
  			if (!this.recording) {
  				ERROR("tried to stop recording, but no recording is in progress.");
  				return;
  			}

  			this.mic.disconnect();
  			this.stream.getTracks()[0].stop();

  			this.recording = false;
  			INFO("stopped recording at " + ts());

  			TRACE(this.buffer);

  			cb && cb(this.buffer);

  			// clear buffer for future recording
  			this.buffer = [];
  		}
  	}]);
  	return Recorder;
  }();

  function makeAudioContext() {
  	var ac = window.AudioContext || window.webkitAudioContext;
  	return new ac();
  }

  function polyfillGetUserMedia() {
  	if (navigator.mediaDevices === undefined) navigator.mediaDevices = {};
  	if (navigator.mediaDevices.getUserMedia === undefined) {
  		WARN("mediaDevices.getUserMedia not found, patching during init");
  		navigator.mediaDevices.getUserMedia = function (constraints, onSuccess, onError) {
  			// Maybe we have a browser-prefixed one?
  			var getUserMedia = navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia || navigator.msGetUserMedia;

  			// Some browsers just don't implement it - return a rejected promise with an error
  			if (!getUserMedia) {
  				return Promise.reject(new Error('getUserMedia not implemented in this browser'));
  			}

  			// Otherwise, wrap getUserMedia in a promise to mimic modern API
  			return new Promise(function (onSuccess, onError) {
  				getUserMedia.call(navigator, constraints, onSuccess, onError);
  			});
  		};
  	}
  }

  var Workspace = function Workspace(hub) {
  	classCallCheck(this, Workspace);

  	this.$pane = $('.workspace .pane');

  	this.overlay = new Overlay(hub);
  	this.controls = new Controls(hub);
  	this.paneScroller = new PaneScroller(hub);
  	this.recorder = new RecordingModal(hub);
  	this.slider = new PlaybackPosition(hub);
  };

  // manages showing and hiding non-interactive overlays

  var Overlay = function () {
  	function Overlay(hub) {
  		var _this = this;

  		classCallCheck(this, Overlay);

  		this.emit = hub.emit;

  		hub.on('projectloadstarted', function () {
  			_this.show(_this.OVERLAYS.WAIT);
  		});
  		hub.on('loadedproject', this.hide);
  		hub.on('projectloadfailed', function (err) {
  			_this.show(_this.OVERLAYS.ERROR, err);
  		});

  		hub.on('recordpressed', function () {
  			_this.show(_this.OVERLAYS.GREYEDOUT);
  		});

  		this.$overlay = $('.workspace .overlay');
  		this.$overlay.find('.close').click(this.hide);
  	}

  	createClass(Overlay, [{
  		key: "show",
  		value: function show(overlay) {
  			for (var _len = arguments.length, args = Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
  				args[_key - 1] = arguments[_key];
  			}

  			var renderedOverlay = typeof overlay == "function" ? overlay.apply(undefined, args) : overlay;

  			this.$overlay.html(renderedOverlay);
  			this.$overlay.show();
  			this.emit('overlayshown');
  		}
  	}, {
  		key: "hide",
  		value: function hide() {
  			this.$overlay.hide();
  			this.emit('overlayhidden');
  		}
  	}]);
  	return Overlay;
  }();

  Overlay.prototype.OVERLAYS = {
  	GREYEDOUT: '<span class="greyed-out-overlay"></span>',
  	EDITING: '<span class="editing-overlay">EDITING...</span>',
  	RENDERING: '<span class="rendering-overlay">RENDERING...</span>',
  	ERROR: function ERROR(err) {
  		"<span class=\"error-overlay\">" + err + "</span>";
  	}
  };

  // manages interacting with the overall play/pause/stop/record controls

  var Controls = function () {
  	function Controls(hub) {
  		classCallCheck(this, Controls);

  		this.emit = hub.emit;

  		this.$allControls = $('.workspace .controls');

  		this.$allControls.find('.control.play').click(this.makeClickControl('playpressed'));
  		this.$allControls.find('.control.pause').click(this.makeClickControl('pausepressed'));
  		this.$allControls.find('.control.stop').click(this.makeClickControl('stoppressed'));
  		this.$allControls.find('.control.record').click(this.makeClickControl('recordpressed'));

  		hub.on('overlayshown', this.disable);
  		hub.on('overlayhidden', this.enable);
  	}

  	createClass(Controls, [{
  		key: "disable",
  		value: function disable() {
  			this.$allControls.disabled = true;
  		}
  	}, {
  		key: "enable",
  		value: function enable() {
  			this.$allControls.disabled = false;
  		}
  	}, {
  		key: "makeClickControl",
  		value: function makeClickControl(eventToEmit) {
  			var _this2 = this;

  			return function () {
  				button.addClass('active');
  				_this2.$allControls.not(button).removeClass('active');
  				_this2.emit(eventToEmit);
  			};
  		}
  	}]);
  	return Controls;
  }();

  // deals with the physical and temporal position of the current-position slider


  var PlaybackPosition = function () {
  	function PlaybackPosition(hub) {
  		var _this3 = this;

  		classCallCheck(this, PlaybackPosition);

  		this.emit = hub.emit;

  		this.$slider = $('.workspace .slider');

  		hub.on('playpressed', this.animate);
  		hub.on('pausepressed', this.pause);
  		hub.on('stoppressed', this.reset);

  		draggable(this.$slider, {
  			x: true,
  			cursor: 'pointer',
  			ondrag: function ondrag() {
  				_this3.emit('slidermoved', _this3.getTime());
  			}
  		});
  	}

  	createClass(PlaybackPosition, [{
  		key: "getTime",
  		value: function getTime() {}
  	}, {
  		key: "getPos",
  		value: function getPos() {}
  	}, {
  		key: "animate",
  		value: function animate() {}
  	}, {
  		key: "pause",
  		value: function pause() {}
  	}, {
  		key: "reset",
  		value: function reset() {}
  	}, {
  		key: "moveToTime",
  		value: function moveToTime(time) {}
  	}, {
  		key: "moveToPos",
  		value: function moveToPos(px) {}
  	}]);
  	return PlaybackPosition;
  }();

  // shows and sets up interactions with the recording modals


  var RecordingModal = function () {
  	function RecordingModal(hub) {
  		var _this4 = this;

  		classCallCheck(this, RecordingModal);

  		this.$modal = $('.recording-modal');
  		this.$countdown = $('.recording-modal .countdown');
  		this.count = 3;
  		this.recorder = new Recorder();

  		this.prefs = {
  			'backingTrack': false };

  		this.$modal.on('changed', 'input.backing-track', this.setPref('backingTrack'));

  		this.$modal.on('click', 'button.close', function () {
  			_this4.$modal.hide();
  			_this4.recorder.stop();
  		});

  		hub.on('recordpressed', this.show);
  	}

  	createClass(RecordingModal, [{
  		key: "setPref",
  		value: function setPref(pref) {
  			var _this5 = this;

  			return function (evt) {
  				_this5.prefs[pref] = evt.checked;
  			};
  		}
  	}, {
  		key: "show",
  		value: function show() {
  			var _this6 = this;

  			this.$modal.show();
  			// TODO switch 1-second interval to something off the metronome
  			this.interval = window.setInterval(function () {
  				_this6.countdown();
  			}, 1000);
  		}
  	}, {
  		key: "countdown",
  		value: function countdown(timeToGo) {
  			if (this.count == 0) {
  				window.clearInterval(this.interval);
  				this.count = 3;
  				this.recorder.start();
  			} else {
  				this.$countdown.html(this.count);
  			}
  		}
  	}, {
  		key: "hide",
  		value: function hide() {
  			var _this7 = this;

  			this.$modal.hide();
  			this.recorder.stop(function (buffer) {
  				_this7.emit('bufferrecorded', buffer);
  			});
  		}
  	}]);
  	return RecordingModal;
  }();

  var w = function Warble(config) {
  	classCallCheck(this, Warble);

  	var hub = new EventHub();
  	this.workspace = new Workspace(hub);
  	this.project = new Project(hub);
  	this.project.load();
  }();

}());
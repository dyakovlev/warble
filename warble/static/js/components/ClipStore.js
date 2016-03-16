// ClipStore holds all the clips in sorted order to enable efficient search and ordered playback

import {INFO, WARN, ERROR, DEBUG, TRACE} from "./utils/misc";

export class ClipStore {
	constructor(){
		this._clipsByStartTime = new Array();
		this._clipsByEndTime = new Array();
	}

	add(clip){
		INFO(`adding clip ${clip.id} to clip store`);

		let startInsertPos = bisect(this._clipsByStartTime, clip.startTime),
			endInsertPos   = bisect(this._clipsByEndTime, clip.endTime);

		this._clipsByStartTime.splice(startInsertPos, 0, [clip.startTime, clip.id]);
		DEBUG(`clip ${clip.id} added to position ${startInsertPos} in start-time list`);

		this._clipsByEndTime.splice(endInsertPos, 0, [clip.endTime, clip.id]);
		DEBUG(`clip ${clip.id} added to position ${endInsertPos} in end-time list`);
	}

	delete(id){
		INFO(`deleting clip ${id} from clip store`);

		let startClipDataIndex = this._clipsByStartTime.findIndex((item) => { item[1].id == id; });
		if (startClipDataIndex == -1){ WARN(`clip ${id} not found in start-time list`); }
		this._clipsByStartTime.splice(startClipDataIndex, 1);

		let endClipDataIndex = this._clipsByEndTime.findIndex((item) => { item[1].id == id; });
		if (endClipDataIndex == -1){ WARN(`clip ${id} not found in end-time list`); }
		this._clipsByEndTime.splice(endClipDataIndex, 1);
	}

	update(clip){
		INFO(`updating clip ${clip.id} in clip store`);

		this.delete(clip.id);
		this.add(clip);
	}

	get(time){
		// get all clip IDs that occur at specific time

		// TODO figure out how this is used
	}

	getRange(start, end){
		// get all clip IDs that start after start and end before end

		// TODO figure out how this is used
	}
}

// binary search-based alg find insertion index for pivot that keeps array sorted
let bisect = function(array, pivot){
	let low = 0, high = array.length, mid;

	while (low < high) {
		mid = (low + high) >> 1;
		(pivot < array[mid][0]) ? high = mid : low = mid + 1;
	}

	return low;
}

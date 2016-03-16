// component that deals with outputting audio through audiocontext

import {INFO, WARN, ERROR, DEBUG, TRACE} from "./utils/misc";
import makeAudioContext from "./utils/audio";

export class Player {
	constructor(hub){
		this.ac = makeAudioContext();

		this.emit = hub.emit;
	}

	play(clips){

	}

	pause(){

	}

	stop(){

	}
}


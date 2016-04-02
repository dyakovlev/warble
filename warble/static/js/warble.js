import {EventHub} from "components/EventHub";
import {Project} from "components/Project";
import {Workspace} from "components/Workspace";


let w = class Warble {
	constructor(config){
		let hub = new EventHub();
		this.workspace = new Workspace(hub);
		this.project = new Project(hub);
		this.project.load();
	}
}();

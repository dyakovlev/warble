(function (components_EventHub,components_Project,components_Workspace) {
  'use strict';

  var classCallCheck = function (instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  };

  var w = function Warble(config) {
  	classCallCheck(this, Warble);

  	var hub = new components_EventHub.EventHub();
  	this.workspace = new components_Workspace.Workspace(hub);
  	this.project = new components_Project.Project(hub);
  	this.project.load();
  }();

}(components_EventHub,components_Project,components_Workspace));
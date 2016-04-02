# Development Notes

### 3/14/16

evaluating a hosting service. heroku, AWS, GAE seem to be mature candidates

https://devcenter.heroku.com/articles/getting-started-with-go

need better MVC adherence.. division of responsibility is not great currently

added an event hub.. has some namespacing provision, unclear if useful

started converting injection scheme into events

moved ClipStore to inside Project

started moving fake modules into real modules

### 3/15/16

priorities today
- clean up JS modules, make sure all imports are right
- break up the injected shims into event handlers
- investigate MOVE vs. MVC vs. other ways to communicate between Model and View

started implementing ClipStore, using two arrays and some traversal logic
cleaned up EventHub, added once(), off()

need to formalize interface for getting things out of ClipStore - maybe some drafts would be good

started whipping Player, Clip, ClipEditor into shape.. just stubbing for now


arrived at a pretty general class setup:

class Wiper {
  constructor(hub){
	// big list of events which this component subscribes to
    hub.on('poo', this.reactToPooEvent);

	// methods of this component can only emit events
	this.emit = hub.emit;
  }

  reactToPooEvent(){
	this.emit('wipe', event, data, goes, here);
  }
}

also would be neat to be able to set log level on a per-module basis but not going to look into
that until the thing at least builds or something

started adding views to Workspace module

## 4/1/16

added draggable() to misc lib.. needs to actually be invented.

expanded a few more workspace widgets. file getting a bit ungainly.. might be worthwhile to make
components/workspace module folder

components/
- Main
- Project
- EventHub
- ClipStore
- workspace/
  - Main
  - Overlay
  - Controls
  - PlaybackPosition
  - RecordingModal
utils/
- misc
- audio
- time

moved Recorder into utils/audio

moved error stuff into utils/error

open question: how do i make PlaybackPosition and RecordingModal talk to each other neatly?

open question: can there be more than one active audio context

https://github.com/eatonphil/smalljs is a near almost-drop-in for jquery

https://github.com/bcherny/draggable good reference for draggable

https://github.com/bcherny/uxhr good reference for xhr

TODO icky, addClip and clipadded events. generally events in Project are fucky



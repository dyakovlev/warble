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

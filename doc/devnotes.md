# Development Notes

## major arch decisions

### object dependency injection

- maybe solved by picking a toolchain and module interface

### build toolchain

use JSPM! offers easiest forward transition to ES6

alternatives:
- vanilla js and make/cat
- RequireJS/almond/AMD packages
- webpack or browserify and CJS packages


## TODO

### misc
- scaffold out project framework and build process

### JS
~ make a config format and make Warble care about it
~ add 'use strict' to places
- convert ghetto modules to ES6 modules
- write or transpile WAV->mp3 JS codec
- write project definition file format, make Warble load it
- write playAll function to schedule clips
+ switch clip IDs to Symobls

### backend
- decide on a backend language
	- go <--
- write storage backend


## Log

### 3/13/16

evaluating a hosting service. heroku, AWS, GAE seem to be mature candidates

https://devcenter.heroku.com/articles/getting-started-with-go

need better MVC adherence.. division of responsibility is not great currently

added an event hub.. has some namespacing provision, unclear if useful

started converting injection scheme into events

moved ClipStore to inside Project

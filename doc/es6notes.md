# variables

var is hoisted

let is block-scope

const is block-scope but complains on reassignment

let and const shadow globals instead of overwriting

## spread operator

... makes a collection-type variable splat out its components

	...varname

# functions

=> fat arrow attaches this to current scope, removes some fiddly bits

python-like default params, pass undefined to use (or omit at tail), no named args still


# objects

## inheritance

mix b props into a

	Object.assign(a, b)

## object and array destructuring

break nested object properties out the weird, crappy way

	let {object structure: { bla bla } =  obj;

arrays work too

	let A = [1, 2, 3];
	let [ , second, third] = A;

also works in function arguments, e.g.

    function eat(food, {type, quantity, tastiness} = default){};


## sets and maps are now legit to use

for example,

	let s = new Set([1,2,3,4]);
	s.add(5)
	console.log(s.has(3)) // true
	console.log(s.size)   // 5
	s.clear()

maps are like objects but they have some nicer lookup properties and can take non-string keys

WeakSet and WeakMap are like Set and Map, but do not block the GC from cleaning up objects if
they are the last reference to those objects

## Generators

function ---> * <--- bla(){ yield val };

calling this function makes a generator that works similar to python


## collection comprehensions

entries(), keys(), values() are things now

## class declarations

similar to ES5 functionality, but slightly different syntax for no absolute fuckin reason yay

	class Bla {
		constructor(bla) {
		}

		memberFunction(){
		}
	};

~or~

	let Bla = class {
		constructor(bla) {
		}

		memberFunction(){
		}
	};

## class extension

	class B {
		member() { }
	}

	class A extends B {
		member() {
			super() // calls B.member, but with whatever params
		}
	}


## arrays

make an array out of something that you can traverse (instead of Array.prototype.slice.call)

    Array.from(iterable)
	Array.from(iterable, (arg) => {transformed output})
	Array.from(iterable, obj.transformer, obj)

make an array out of a list of args (sidesteps some constructor ambiguity)

    Array.of(arg1, arg2, arg3)

find values and indices in array with a condition

	numbers = [1,2,35,4]
	console.log(numbers.find(n => n > 33));         // 35
	console.log(numbers.findIndex(n => n > 33));    // 2

fill elements of an array

	let numbers = [1, 2, 3, 4];
	numbers.fill(1, 2);    // 1,2,1,1
	numbers.fill(0, 1, 3); // 1,0,0,1


## arraybuffers

Arrays hold random crap, resize as needed (chunky lists or smth)
buffers hold literal hard cold bits

views allow you to interact with them, might be useful for manually editing WAVs

typed arrays are a shorthand for declaring a buffer and a typed view over it

	let ints = new Int16Array(2),
		floats = new Float32Array(5);

	ints.byteLength == 4
	ints.length == 2

	floats.byteLength == 20
	floats.length == 5

## promises



// math
var clamp = (a, b, c) => Math.min(Math.max(a, b), c)


// let things be dragged
export function draggable(el, opts){
	let defaults = {
		x: false, // do not allow drag in the x direction
		y: false, // do not allow drag in the y direction
		cursor: 'pointer', // change cursor to pointer when hovering over el
		ondrag: () => {}, // call this function when a drag event is finished (mouseup)
	}

	// update opts with defaults

	// set up cursor

	// set up drag-ability

	// fire ondrag on mouse up
}

/* Audio utils */

function generateGUMPatch(){
	return function(constraints, successCallback, errorCallback) {
		var getUserMedia = (navigator.getUserMedia ||
				navigator.webkitGetUserMedia ||
				navigator.mozGetUserMedia ||
				navigator.msGetUserMedia);

		// Some browsers just don't implement it - return a rejected promise with an error
		if(!getUserMedia) {
			return Promise.reject(new Error('getUserMedia is not implemented in this browser'));
		}

		// Otherwise, wrap getUserMedia in a promise to mimic modern API
		return new Promise(function(successCallback, errorCallback) {
			getUserMedia.call(navigator, constraints, successCallback, errorCallback);
		});
	}
}

(function () {
	if (!window._js_hooked) {
		window._js_hooked = true;
		_raw_send_method = XMLHttpRequest.prototype.send
		XMLHttpRequest.prototype.send = function () {
		}
	};
})();
(function() {
    if (!window._js_hooked) {
        window._js_hooked = true;

        //disable all ajax request
        _raw_send_method = XMLHttpRequest.prototype.send
        XMLHttpRequest.prototype.send = function() {}

        
    };
})();

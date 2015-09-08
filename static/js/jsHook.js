/*
This js file was injected into the proxied html for some functional enhancements
*/
(function() {
    if (!window._js_hooked) {
        window._js_hooked = true;

        //disable all ajax request
        var _raw_send_method = XMLHttpRequest.prototype.send
        XMLHttpRequest.prototype.send = function() {}
        document.write('<script src="/js/statusBar.js" type="text/javascript"></script>');
    };
})();

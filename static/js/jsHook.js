/*
 * @Author gongshw
 *
 * This js file was injected into the proxied html for some functional enhancements
 */
(function() {
    if (!window._js_hooked) {
        window._js_hooked = true;

        function proxyUrl(url) {
            if (url.startsWith('http://') || url.startsWith('https://')) {
                var token = url.match(/^(http.?):\/\/(.*)$/);
                return '/proxy/' + token[1] + '/' + token[2];
            } else if (url.startsWith('//')) {
                return location.pathname.split('/', 3).join('/') + '/' + url.substring(2)
            } else if (url.startsWith('/')) {
                return location.pathname.split('/', 4).join('/') + '/' + url.substring(1)
            } else {
                return url;
            }
        }

        //disable all ajax request
        var _raw_send_method = XMLHttpRequest.prototype.send
        XMLHttpRequest.prototype.send = function() {}

        function hookProperty(type, attrName) {
            var _raw_prop = Object.getOwnPropertyDescriptor(type.prototype, attrName);
            Object.defineProperty(type.prototype, attrName, {
                set: function function_name(url) {
                    var proxiedUrl = proxyUrl(url);
                    _raw_prop.set.call(this, proxiedUrl);
                },
                get: _raw_prop.get
            });
        }

        //hook HTMLScriptElement.src property
        var elementToHook = {
            'HTMLScriptElement':['src'],
            'HTMLImageElement':['src'],
        }
        for (var eleName in elementToHook) {
            for (var i = elementToHook[eleName].length - 1; i >= 0; i--) {
                hookProperty(window[eleName], elementToHook[eleName][i]);
            };
        };


        document.write('<script src="/js/statusBar.js" type="text/javascript"></script>');
    };
})();

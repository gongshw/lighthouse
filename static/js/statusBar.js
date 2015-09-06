(function() {
    document.addEventListener("DOMContentLoaded", function(event) {
        var toggleMenu = document.createElement('span')
        toggleMenu.id = 'toggleMenu';
        toggleMenu.title = '点击展开Lighthouse菜单'
        var cssLink = document.createElement('link');
        cssLink.setAttribute('rel', 'stylesheet');
        cssLink.setAttribute("type", "text/css")
        cssLink.setAttribute('href', '/css/statusBar.css');
        document.body.appendChild(toggleMenu);
        document.getElementsByTagName("head")[0].appendChild(cssLink);
    });
})();

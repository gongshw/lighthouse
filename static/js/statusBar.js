(function() {
    document.addEventListener("DOMContentLoaded", function(event) {

        var html = '<div class="reset-this"><div id="lighthouseMenu" class="closed">' +
            '<div class="title"><a href="https://github.com/gongshw/lighthouse">Lighthouse</a></div>' +
            '<div class="links"><a href="/">Go to lighthouse home</a>' +
            '<a href="javascript:void(0)">Visit this page directly</a></div>' +
            '<span id="menuToggle" title="点击展开Lighthouse菜单"></span>' +
            '</div></div>';
        document.body.innerHTML += html;
        var menuToggle = document.getElementById('menuToggle');
        var lighthouseMenu = document.getElementById('lighthouseMenu');
        var cssLink = document.createElement('link');
        cssLink.setAttribute('rel', 'stylesheet');
        cssLink.setAttribute("type", "text/css")
        cssLink.setAttribute('href', '/css/statusBar.css');
        document.getElementsByTagName("head")[0].appendChild(cssLink);
        menuToggle.addEventListener('click', function() {
            if (lighthouseMenu.className == 'closed') {
                lighthouseMenu.className = 'opened';
            } else if (lighthouseMenu.className == 'opened') {
                lighthouseMenu.className = 'closed';
            };;
        })
    });
})();

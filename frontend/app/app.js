(function() {
  "use strict";

  var app = angular.module("app", [
    "ngRoute",
    "ngResource",
    "ngAnimate",
    "ui.bootstrap",
    "luegg.directives",
    "hljs",
    "chart.js",
    "app",
    "app.filters",
    "app.services"
  ]);

  app.config(function(hljsServiceProvider) {
    hljsServiceProvider.setOptions({ tabReplace: "  " });
  });

  app.run(["$route", function($route) {}]);
})();

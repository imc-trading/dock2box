(function() {
  "use strict";
  var app = angular.module("app");

  app.config([
    "$compileProvider",
    function($compileProvider) {
      $compileProvider.aHrefSanitizationWhitelist(/^\s*(http?|https?|ssh):/);
    }
  ]);
})();

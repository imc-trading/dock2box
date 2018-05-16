(function() {
  "use strict";
  angular.module("app.services").factory("state", ["$http", state]);

  function state($http) {
    var service = {
      setTab: setTab,
      getTab: getTab
    };

    return service;

    var activeTab = "hosts";

    function setTab(tab) {
      activeTab = tab;
    }

    function getTab(tab) {
      return activeTab;
    }
  }
})();

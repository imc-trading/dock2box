(function() {
  "use strict";
  angular.module("app").controller("layout", ["state", layout]);

  function layout(state) {
    var vm = this;
    vm.getTab = getTab;

    function getTab() {
      return state.getTab();
    }
  }
})();

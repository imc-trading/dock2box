(function() {
  "use strict";

  angular
    .module("app")
    .controller("taskDefSettings", [
      "$scope",
      "$uibModalInstance",
      "data",
      "taskDef",
      taskDefSettings
    ]);

  function taskDefSettings($scope, $uibModalInstance, data, taskDef) {
    var vm = this;
    vm.taskDef = taskDef;
    vm.alerts = [];

    if (vm.taskDef.env !== undefined) {
      var keys = Object.keys(vm.taskDef.env);
      if (keys.length > 0) {
        vm.env = vm.taskDef.env[keys[0]];
      }
    }

    if (vm.taskDef.pre !== undefined) {
      var keys = Object.keys(vm.taskDef.pre);
      if (keys.length > 0) {
        vm.pre = vm.taskDef.pre[keys[0]];
      }
    }

    if (vm.taskDef.post !== undefined) {
      var keys = Object.keys(vm.taskDef.post);
      if (keys.length > 0) {
        vm.post = vm.taskDef.post[keys[0]];
      }
    }

    vm.closeAlert = function(index) {
      vm.alerts.splice(index, 1);
    };

    vm.close = function() {
      $uibModalInstance.close("cancel");
    };
  }
})();

(function() {
  "use strict";

  angular
    .module("app")
    .controller("taskStop", [
      "$scope",
      "$uibModalInstance",
      "data",
      "task",
      taskStop
    ]);

  function taskStop($scope, $uibModalInstance, data, task) {
    var vm = this;
    vm.task = task;
    vm.alerts = [];

    vm.closeAlert = function(index) {
      vm.alerts.splice(index, 1);
    };

    vm.stop = function() {
      data.stopTask(vm.task).then(function(data) {
        console.log(data);
        if (data.status !== undefined && data.status == 200) {
          $uibModalInstance.close("stopTask");
        } else {
          vm.alerts.push({ msg: data });
        }
      });
    };

    vm.close = function() {
      $uibModalInstance.close("close");
    };
  }
})();

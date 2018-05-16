(function() {
  "use strict";

  angular
    .module("app")
    .controller("hostRun", [
      "$scope",
      "$uibModalInstance",
      "data",
      "host",
      hostRun
    ]);

  function hostRun($scope, $uibModalInstance, data, host) {
    var vm = this;
    vm.host = host;
    vm.alerts = [];

    data.getHostTaskDefs(vm.host.uuid).then(function(taskDefs) {
      vm.taskDefs = taskDefs;
      if (vm.taskDefs.length > 0) {
        vm.taskDef = taskDefs[0];
        vm.taskDef.value = vm.taskDef.name;
      }
    });

    vm.closeAlert = function(index) {
      vm.alerts.splice(index, 1);
    };

    vm.run = function() {
      var runTask = {
        hostUUID: vm.host.uuid,
        hostname: vm.host.hostname,
        taskDef: vm.taskDef.name
      };

      var json = angular.toJson(runTask);
      data.createTask(json).then(function(data) {
        console.log(data);
        if (data.status !== undefined && data.status == 200) {
          $uibModalInstance.close("run");
        } else {
          vm.alerts.push({ msg: data });
        }
      });
    };

    vm.close = function() {
      $uibModalInstance.close("cancel");
    };
  }
})();

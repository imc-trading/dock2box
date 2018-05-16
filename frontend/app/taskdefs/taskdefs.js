(function() {
  "use strict";

  angular
    .module("app")
    .controller("taskDefs", [
      "$uibModal",
      "$interval",
      "$filter",
      "$scope",
      "state",
      "data",
      taskDefs
    ]);

  function taskDefs($uibModal, $interval, $filter, $scope, state, data) {
    state.setTab("taskdefs");

    var vm = this;
    vm.module = "taskdefs";
    vm.tableControls = true;

    vm.fields = [
      { name: "Name", field: "name", kind: "string", show: true },
      { name: "Description", field: "descr", kind: "string", show: true },
      { name: "User", field: "user", kind: "string" },
      { name: "Group", field: "group", kind: "string" },
      { name: "Work Dir.", field: "dir", kind: "string" },
      { name: "Environment", field: "env", kind: "string", dict: true },
      { name: "Command", field: "cmd", kind: "string", show: true },
      { name: "Arguments", field: "args", kind: "string", array: true },
      { name: "File Exist", field: "fileExist", kind: "string" },
      { name: "File Not Exist", field: "fileNotExist", kind: "string" },
      { name: "Reboot", field: "reboot", kind: "boolean" },
      {
        name: "RebootExitCode",
        field: "rebootExitCode",
        kind: "integer",
        show: true
      },
      { name: "RebootFileExist", field: "rebootFileExist", kind: "string" },
      {
        name: "RebootFileNotExist",
        field: "rebootFileNotExist",
        kind: "string"
      },
      {
        name: "Concurrency",
        field: "concurrency",
        kind: "integer",
        show: true
      },
      { name: "Timeout", field: "timeout", kind: "integer", show: true }
    ];

    vm.unfilteredRows = [];
    vm.rows = [];
    vm.rowLimit = 10;
    vm.page = 1;
    vm.numPages = 1;
    vm.fromPage = 1;
    vm.toPage = 1;
    vm.orderBy = "name";
    vm.reverse = false;
    vm.refreshing = false;
    vm.refresh = refresh;
    vm.setPage = setPage;
    vm.setRowLimit = setRowLimit;
    vm.getSortClass = getSortClass;
    vm.setOrderBy = setOrderBy;
    vm.getRowClass = getRowClass;

    refresh();

    /*
    vm.intervalPromise = $interval(function() {
      console.log("auto refresh tasks");
      vm.refresh();
    }, 5000);

    vm.$onDestroy = function() {
      console.log("cancel auto refresh tasks");
      $interval.cancel(vm.intervalPromise);
    };
*/

    function refresh() {
      if (vm.refreshing) {
        return;
      }
      vm.refreshing = true;

      data.getTaskDefs().then(function(data) {
        vm.unfilteredRows = data;
        vm.refreshing = false;
        filter();
      });
    }

    function filter() {
      var showFields = [];
      for (var k in vm.fields) {
        if (vm.fields[k].show) {
          showFields.push(vm.fields[k].field);
        }
      }

      if (typeof vm.search !== "undefined" || vm.search != "") {
        vm.rows = $filter("booleanSearch")(
          vm.unfilteredRows,
          vm.search,
          showFields
        );
      } else {
        vm.rows = vm.unfilteredRows;
      }

      paging();
    }

    function paging() {
      vm.numPages = Math.ceil(vm.rows.length / vm.rowLimit);
      vm.fromPage = vm.page - 3;
      vm.toPage = vm.page + 3;

      if (vm.fromPage < 1) {
        var diff = 1 - vm.fromPage;
        vm.toPage += diff;
      } else if (vm.toPage > vm.numPages) {
        var diff = vm.toPage - vm.numPages;
        vm.fromPage -= diff;
      }

      if (vm.fromPage < 1) {
        vm.fromPage = 1;
      }
      if (vm.toPage > vm.numPages) {
        vm.toPage = vm.numPages;
      }
    }

    $scope.$watch("vm.search", function(nVal, oVal, scope) {
      vm.page = 1;
      filter();
    });

    function setPage(page) {
      vm.page = page;
      paging();
    }

    function setRowLimit(rowLimit) {
      vm.rowLimit = rowLimit;
      vm.page = 1;
      paging();
    }

    function getSortClass(field) {
      if (field == vm.orderBy) {
        if (vm.reverse) {
          return "sort-desc";
        }
        return "sort-asc";
      }
    }

    function setOrderBy(field) {
      if (field == vm.orderBy) {
        vm.reverse = !vm.reverse;
        return;
      }
      vm.orderBy = field;
      vm.reverse = false;
    }

    function getRowClass(row) {
      return "table-row-striped";
    }

    vm.pullTaskDefs = function() {
      console.log("DO");
      data.pullTaskDefs().then(function(data) {
        console.log(data);
        if (data.status !== undefined && data.status == 200) {
        } else {
          //          vm.alerts.push({ msg: data });
        }
      });
    };

    vm.settings = function(taskDef) {
      $uibModal
        .open({
          size: "lg",
          templateUrl: "taskdefs/settings/settings.html",
          controller: "taskDefSettings",
          controllerAs: "vm",
          resolve: {
            taskDef: function() {
              return taskDef;
            }
          }
        })
        .result.then(function() {}, function(res) {});
    };
  }
})();

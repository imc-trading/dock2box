(function() {
  "use strict";

  angular
    .module("app")
    .controller("tasks", [
      "$uibModal",
      "$interval",
      "$filter",
      "$scope",
      "state",
      "data",
      tasks
    ]);

  function tasks($uibModal, $interval, $filter, $scope, state, data) {
    state.setTab("tasks");

    var vm = this;
    vm.module = "tasks";
    vm.tableControls = true;
    vm.tableButtons = true;

    vm.fields = [
      { name: "Created", field: "created", kind: "datetime", show: true },
      { name: "Updated", field: "updated", kind: "datetime" },
      { name: "Hostname", field: "hostname", kind: "string", show: true },
      { name: "Task Def.", field: "taskDef", kind: "string" },
      { name: "Name", field: "name", kind: "string", show: true },
      { name: "Environment", field: "env", kind: "string", dict: true },
      { name: "Command", field: "cmd", kind: "string" },
      { name: "Args", field: "args", kind: "string", array: true },
      { name: "Running", field: "running", kind: "boolean", show: true },
      { name: "Terminated", field: "terminated", kind: "boolean", show: true },
      { name: "Orphaned", field: "orphaned", kind: "boolean", show: true },
      { name: "Started", field: "started", kind: "datetime" },
      { name: "Finished", field: "finished", kind: "datetime" },
      {
        name: "Progress / Duration",
        field: "progress",
        kind: "progress",
        dispField: "durationLeft",
        dispKind: "duration",
        show: true,
        progress: true
      },
      { name: "Exit Code", field: "exitCode", kind: "integer" },
      { name: "Tail", field: "tail", kind: "string", array: true }
    ];

    vm.unfilteredRows = [];
    vm.rows = [];
    vm.rowLimit = 10;
    vm.page = 1;
    vm.numPages = 1;
    vm.fromPage = 1;
    vm.toPage = 1;
    vm.orderBy = "created";
    vm.reverse = true;
    vm.refreshing = false;
    vm.refresh = refresh;
    vm.setPage = setPage;
    vm.setRowLimit = setRowLimit;
    vm.getSortClass = getSortClass;
    vm.setOrderBy = setOrderBy;
    vm.getRowClass = getRowClass;

    refresh();

    vm.intervalPromise = $interval(function() {
      console.log("auto refresh tasks");
      vm.refresh();
    }, 5000);

    vm.$onDestroy = function() {
      console.log("cancel auto refresh tasks");
      $interval.cancel(vm.intervalPromise);
    };

    function refresh() {
      if (vm.refreshing) {
        return;
      }
      vm.refreshing = true;

      data.getTasks().then(function(data) {
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
      if (row.running) {
        return "table-row-success";
      }
      if (row.exitCode != 0) {
        return "table-row-danger";
      }
      return "table-row-striped";
    }

    vm.taskStop = function(task) {
      $uibModal
        .open({
          size: "sm",
          templateUrl: "tasks/stop/stop.html",
          controller: "taskStop",
          controllerAs: "vm",
          resolve: {
            task: function() {
              return task;
            }
          }
        })
        .result.then(function() {}, function(res) {});
    };

    vm.logs = function(task) {
      $uibModal
        .open({
          size: "lg",
          templateUrl: "tasks/logs/logs.html",
          controller: "taskLogs",
          controllerAs: "vm",
          resolve: {
            task: function() {
              return task;
            }
          }
        })
        .result.then(function() {}, function(res) {});
    };
  }
})();

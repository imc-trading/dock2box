(function() {
  "use strict";

  angular
    .module("app")
    .controller("hosts", [
      "$uibModal",
      "$interval",
      "$filter",
      "$scope",
      "state",
      "data",
      hosts
    ]);

  function hosts($uibModal, $interval, $filter, $scope, state, data) {
    state.setTab("hosts");

    var vm = this;
    vm.module = "hosts";
    vm.tableControls = true;
    vm.tableButtons = true;

    vm.labels = ["Running", "Successful", "Failed"];
    vm.data = [1, 10, 1];
    vm.colors = ["#03BB4F", "#4c9fd1", "#ff3232"];
    vm.options = {
      cutoutPercentage: 40
    };

    vm.datasets = {
      backgroundColor: ["#03BB4F", "#4c9fd1", "#ff3232"],
      hoverBackgroundColor: ["#03e260", "#71b2da", "#ff6060"],
      hoverBorderColor: "rgba(28, 32, 34, .9)",
      borderColor: "rgba(28, 32, 34, .7)",
      borderWidth: 5,
      hoverBorderWidth: 10
    };

    vm.fields = [
      { name: "Created", field: "created", kind: "datetime", show: true },
      { name: "Updated", field: "updated", kind: "datetime" },
      { name: "Hostname", field: "hostname", kind: "string", show: true },
      {
        name: "Running",
        field: "running",
        kind: "integer",
        show: true,
        class: "text-success"
      },
      { name: "Address", field: "addr", kind: "string" },
      {
        name: "Distribution",
        field: "distroName",
        kind: "string",
        show: true
      },
      {
        name: "Release",
        field: "distroRelease",
        kind: "string",
        show: true
      },
      {
        name: "Code Name",
        field: "distroCodeName",
        kind: "string"
      },
      { name: "CPU Model", field: "cpuModel", kind: "string" },
      { name: "CPU Flags", field: "cpuFlags", kind: "string" },
      {
        name: "CPU Speed GHz",
        field: "cpuSpeedGHz",
        kind: "number"
      },
      {
        name: "CPU Logical Cores",
        field: "cpuLogicalCores",
        kind: "integer"
      },
      {
        name: "CPU Physical Cores",
        field: "cpuPhysicalCores",
        kind: "integer"
      },
      { name: "CPU Sockets", field: "cpuSockets", kind: "integer" },
      {
        name: "CPU Cores Per Socket",
        field: "cpuCoresPerSocket",
        kind: "integer"
      },
      {
        name: "CPU Threads Per Core",
        field: "cpuThreadsPerCore",
        kind: "integer"
      },
      { name: "Mem. Total GB", field: "memoryTotalGB", kind: "number" },
      {
        name: "Manufacturer",
        field: "manufacturer",
        kind: "string"
      },
      { name: "Product", field: "product", kind: "string", show: true },
      {
        name: "Product Version",
        field: "productVersion",
        kind: "string",
        show: true
      },
      {
        name: "Serial Number",
        field: "serialNumber",
        kind: "string",
        show: true
      }
    ];
    vm.unfilteredRows = [];
    vm.rows = [];
    vm.rowLimit = 10;
    vm.page = 1;
    vm.numPages = 1;
    vm.fromPage = 1;
    vm.toPage = 1;
    vm.orderBy = "alive";
    vm.reverse = true;
    vm.refreshing = false;
    vm.refresh = refresh;
    vm.setPage = setPage;
    vm.setRowLimit = setRowLimit;
    vm.getSortClass = getSortClass;
    vm.setOrderBy = setOrderBy;
    vm.getRowClass = getRowClass;

    vm.openHardware = openHardware;
    vm.openCreate = openCreate;
    vm.openEdit = openEdit;
    vm.openRun = openRun;

    refresh();

    vm.intervalPromise = $interval(function() {
      console.log("auto refresh hosts");
      vm.refresh();
    }, 5000);

    vm.$onDestroy = function() {
      console.log("cancel auto refresh hosts");
      $interval.cancel(vm.intervalPromise);
    };

    function refresh() {
      if (vm.refreshing) {
        return;
      }
      vm.refreshing = true;

      data.getHosts().then(function(data) {
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
      if (row.running > 0) {
        return "table-row-success";
      }
      return "table-row-striped";
    }

    function openHardware(data) {
      $uibModal
        .open({
          size: "lg",
          templateUrl: "hosts/hardware/hardware.html",
          controller: "hardware",
          controllerAs: "vm",
          resolve: {
            data: function() {
              return data;
            }
          }
        })
        .result.then(function() {}, function(res) {});
    }

    function openCreate(host) {
      $uibModal
        .open({
          size: "lg",
          templateUrl: "hosts/create/create.html",
          controller: "hostCreate",
          controllerAs: "vm",
          resolve: {
            host: function() {
              return host;
            }
          }
        })
        .result.then(function() {}, function(res) {});
    }

    function openEdit(host) {
      $uibModal
        .open({
          size: "lg",
          templateUrl: "hosts/edit/edit.html",
          controller: "hostEdit",
          controllerAs: "vm",
          resolve: {
            host: function() {
              return host;
            }
          }
        })
        .result.then(function() {}, function(res) {});
    }

    function openRun(host) {
      $uibModal
        .open({
          size: "md",
          templateUrl: "hosts/run/run.html",
          controller: "hostRun",
          controllerAs: "vm",
          resolve: {
            host: function() {
              return host;
            }
          }
        })
        .result.then(function() {}, function(res) {});
    }
  }
})();

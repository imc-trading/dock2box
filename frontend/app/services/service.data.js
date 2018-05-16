(function() {
  "use strict";
  angular.module("app.services").factory("data", ["$http", data]);

  function data($http) {
    var service = {
      getHost: getHost,
      getHosts: getHosts,
      getStatuses: getStatuses,
      getSubnets: getSubnets,
      getSites: getSites,
      getRoles: getRoles,
      getEnvironments: getEnvironments,
      getTenants: getTenants,
      getSecurityGroups: getSecurityGroups,
      getImages: getImages,
      getTaskDefs: getTaskDefs,
      getHostTaskDefs: getHostTaskDefs,
      createHost: createHost,
      updateHost: updateHost,
      getTasks: getTasks,
      createTask: createTask,
      stopTask: stopTask,
      pullTaskDefs: pullTaskDefs
    };

    return service;

    function getHost(host) {
      return $http({
        method: "GET",
        url: "/api/hosts/" + host
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getHosts() {
      return $http({
        method: "GET",
        url: "/api/hosts"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getStatuses() {
      return $http({
        method: "GET",
        url: "/api/statuses"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getSubnets() {
      return $http({
        method: "GET",
        url: "/api/subnets"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getSites() {
      return $http({
        method: "GET",
        url: "/api/sites"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getRoles() {
      return $http({
        method: "GET",
        url: "/api/roles"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getEnvironments() {
      return $http({
        method: "GET",
        url: "/api/environments"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getTenants() {
      return $http({
        method: "GET",
        url: "/api/tenants"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getSecurityGroups() {
      return $http({
        method: "GET",
        url: "/api/securitygroups"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getImages() {
      return $http({
        method: "GET",
        url: "/api/images"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getTaskDefs() {
      return $http({
        method: "GET",
        url: "/api/taskdefs"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getHostTaskDefs(hostUUID) {
      return $http({
        method: "GET",
        url: "/api/hosts/" + hostUUID + "/taskdefs"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function createHost(data) {
      return $http({
        method: "POST",
        url: "/api/hosts",
        data: data
      }).then(
        function(data, status, headers, config) {
          return data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function updateHost(data) {
      return $http({
        method: "PUT",
        url: "/api/hosts/" + data.uuid,
        data: data
      }).then(
        function(data, status, headers, config) {
          return data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function getTasks() {
      return $http({
        method: "GET",
        url: "/api/tasks"
      }).then(
        function(data, status, headers, config) {
          return data.data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function createTask(data) {
      return $http({
        method: "POST",
        url: "/api/tasks",
        data: data
      }).then(
        function(data, status, headers, config) {
          return data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function stopTask(task) {
      return $http({
        method: "PUT",
        url: "/api/hosts/" + task.hostUUID + "/tasks/" + task.name + "/stop"
      }).then(
        function(data, status, headers, config) {
          return data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }

    function pullTaskDefs() {
      return $http({
        method: "PUT",
        url: "http://dock2box-dev/api/taskdefs/pull"
      }).then(
        function(data, status, headers, config) {
          return data;
        },
        function(error) {
          console.log(error);
          return error;
        }
      );
    }
  }
})();

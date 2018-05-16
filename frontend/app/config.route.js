(function() {
  "use strict";
  var app = angular.module("app");

  app.constant("routes", getRoutes());

  app.config(["$routeProvider", "routes", routeConfigurator]);
  function routeConfigurator($routeProvider, routes) {
    routes.forEach(function(r) {
      $routeProvider.when(r.url, r.config);
    });
    $routeProvider.otherwise({ redirectTo: "/" });
  }

  function getRoutes() {
    return [
      {
        url: "/",
        config: {
          templateUrl: "hosts/hosts.html",
          title: "hosts"
        }
      },
      {
        url: "/hosts",
        config: {
          templateUrl: "hosts/hosts.html",
          title: "hosts"
        }
      },
      {
        url: "/tasks",
        config: {
          templateUrl: "tasks/tasks.html",
          title: "tasks"
        }
      },
      {
        url: "/taskdefs",
        config: {
          templateUrl: "taskdefs/taskdefs.html",
          title: "taskdefs"
        }
      }
    ];
  }
})();

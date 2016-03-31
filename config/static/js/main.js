var app = angular.module('dock2box', [
  'ngRoute',
  'ngResource',
  'ui.bootstrap',
  'smart-table'
])

/*
 * Routes
 */

app.config(['$routeProvider', function ($routeProvider) {
  $routeProvider
    .when("/", {templateUrl: "partials/dashboard.html", controller: "PageCtrl"})
    .when("/hosts", {templateUrl: "partials/hosts.html", controller: "PageCtrl"})
    .when("/hosts/edit", {templateUrl: "partials/edit_host.html", controller: "PageCtrl"})
    .when("/hosts/overview", {templateUrl: "partials/host/overview.html", controller: "PageCtrl"})
    .when("/hosts/network", {templateUrl: "partials/host/network.html", controller: "PageCtrl"})
    .when("/hosts/storage", {templateUrl: "partials/host/storage.html", controller: "PageCtrl"})
    .when("/hosts/system", {templateUrl: "partials/host/system.html", controller: "PageCtrl"})
    .when("/hosts/docker", {templateUrl: "partials/host/docker.html", controller: "PageCtrl"})
    .when("/hosts/new", {templateUrl: "partials/new_host.html", controller: "PageCtrl"})
    .when("/images", {templateUrl: "partials/images.html", controller: "PageCtrl"})
}]);

/**
 * Controls all other Pages
 */
app.controller('PageCtrl', function (/* $scope, $location, $http */) {
  console.log("Page Controller reporting for duty.");
});

/*
 * Filters
 */

// joinBy
app.filter('joinBy', function () {
  return function (input,delimiter) {
    return (input || []).join(delimiter || ',');
  };
});

// replace
app.filter('replace', function () {
  return function (input,restr,newstr) {
    var re = new RegExp(restr);
    return input.replace(re, newstr)
  };
});

/*
 * Controllers
 */

// Hosts
app.controller('hostsController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/hosts?table=true');

  resource.query().$promise.then(function(value) {
    $scope.hosts = value;
//    console.log (value); 
  });

} ]);

// Images
app.controller('imagesController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/images?table=true');

  resource.query().$promise.then(function(value) {
    $scope.images = value;
//    console.log (value);
  });
} ]);

/*
 * Controllers Peekaboo
 */

// CPU
app.controller('cpuController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/system/cpu');

  resource.get().$promise.then(function(value) {
    $scope.cpu = value;
//    console.log (value);
  });
} ]);

// Memory
app.controller('memoryController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/system/memory');

  resource.get().$promise.then(function(value) {
    $scope.memory = value;
//    console.log (value);
  });
} ]);

// System
app.controller('systemController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/system');

  resource.get().$promise.then(function(value) {
    $scope.system = value;
//    console.log (value);
  });
} ]);

// OS
app.controller('osController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/system/os');

  resource.get().$promise.then(function(value) {
    $scope.os = value;
//    console.log (value);
  });
} ]);

// Interfaces
app.controller('interfacesController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/network/interfaces');

  resource.query().$promise.then(function(value) {
    $scope.interfaces = value;
//    console.log (value);
  });
} ]);

// Routes
app.controller('routesController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/network/routes');

  resource.query().$promise.then(function(value) {
    $scope.routes = value;
//    console.log (value);
  });
} ]);

// Disks
app.controller('disksController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/storage/disks');

  resource.query().$promise.then(function(value) {
    $scope.disks = value;
//    console.log (value);
  });
} ]);

// Mounts
app.controller('mountsController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/storage/mounts');

  resource.query().$promise.then(function(value) {
    $scope.mounts = value;
//    console.log (value);
  });
} ]);

// LVM Physical Volumes
app.controller('physVolsController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/storage/lvm/physvols');

  resource.query().$promise.then(function(value) {
    $scope.physVols = value;
//    console.log (value);
  });
} ]);

// LVM Logical Volumes
app.controller('logVolsController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/storage/lvm/logvols');

  resource.query().$promise.then(function(value) {
    $scope.logVols = value;
//    console.log (value);
  });
} ]);

// Sysctls
app.controller('sysctlsController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/system/sysctls');

  resource.query().$promise.then(function(value) {
    $scope.sysctls = value;
//    console.log (value);
  });
} ]);

// Docker
app.controller('dockerController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/docker');

  resource.get().$promise.then(function(value) {
    $scope.docker = value;
//    console.log (value);
  });
} ]);

// Images
/*
app.controller('imagesController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/docker/images');

  resource.query().$promise.then(function(value) {
    $scope.images = value;
//    console.log (value);
  });
} ]);
*/

// Containers
app.controller('containersController', [ '$scope', '$resource', function($scope, $resource) {
  var resource = $resource('/api/v1/docker/containers');

  resource.query().$promise.then(function(value) {
    $scope.containers = value;
//    console.log (value);
  });
} ]);

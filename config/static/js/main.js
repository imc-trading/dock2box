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
  return function (input,oldstr,newstr) {
    var re = new RegExp(oldstr, "g");
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

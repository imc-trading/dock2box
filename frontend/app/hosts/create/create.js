(function() {
  "use strict";

  angular
    .module("app")
    .controller("hostCreate", [
      "$scope",
      "$uibModalInstance",
      "data",
      "host",
      hostCreate
    ]);

  function hostCreate($scope, $uibModalInstance, data, host) {
    var vm = this;
    vm.alerts = [];
    vm.host = host;
    vm.domain = "";
    vm.hostname = host.hostname;
    vm.allowBuild = true;

    data.getStatuses().then(function(statuses) {
      vm.statuses = statuses;
      vm.status = vm.statuses.find(function(status) {
        return status.status == "";
      });
      if (vm.status !== undefined) {
        vm.status.value = vm.status.status;
      }
    });

    data.getRoles().then(function(roles) {
      vm.roles = roles;
      vm.role = vm.roles.find(function(role) {
        return role.role == "";
      });
      if (vm.role !== undefined) {
        vm.role.value = vm.role.role;
      }
    });

    data.getImages().then(function(images) {
      vm.images = images;
      vm.image = vm.images.find(function(image) {
        return image.image == "";
      });
      if (vm.image !== undefined) {
        vm.image.value = vm.image.image;
      }
    });

    data.getTenants().then(function(tenants) {
      vm.tenant = "";
      vm.tenants = tenants;
      vm.tenant = vm.tenants.find(function(tenant) {
        return tenant.tenant == "";
      });
      if (vm.tenant !== undefined) {
        vm.tenant.value = vm.tenant.tenant;
      }
    });

    data.getSecurityGroups().then(function(securityGroups) {
      vm.securityGroups = securityGroups;
      vm.securityGroup = vm.securityGroups.find(function(group) {
        return group.group == "";
      });
      if (vm.securityGroup !== undefined) {
        vm.securityGroup.value = vm.securityGroup.group;
      }
    });

    data.getEnvironments().then(function(environments) {
      vm.environments = environments;
      vm.environment = vm.environments.find(function(environment) {
        return environment.environment == "";
      });
      if (vm.environment !== undefined) {
        vm.environment.value = vm.environment.environment;
      }
    });

    var promises = [
      data.getSubnets().then(function(subnets) {
        vm.subnets = subnets;
      }),
      data.getSites().then(function(sites) {
        vm.sites = sites;
      })
    ];

    Promise.all(promises).then(function() {
      var intf = vm.host.interfaces.find(function(intf) {
        return intf.primary;
      });

      if (intf !== undefined) {
        vm.setInterface(intf);
      }

      $scope.$apply();
    });

    vm.setInterface = function(intf) {
      vm.interface = intf;
      vm.interface.value = intf.hwAddr;
      if (intf.network !== undefined) {
        var subnet = vm.subnets.find(function(subnet) {
          return subnet.subnet == intf.network;
        });
        if (subnet !== undefined) {
          vm.setSubnet(subnet);
        }
      }
    };

    vm.setSubnet = function(subnet) {
      vm.subnet = subnet;
      vm.subnet.value = subnet.subnet;

      if (subnet.siteID !== undefined) {
        var site = vm.sites.find(function(site) {
          return site.id == subnet.siteID;
        });
        if (site !== undefined) {
          vm.site = site;
          vm.site.value = site.site;
        }
      }
    };

    vm.closeAlert = function(index) {
      vm.alerts.splice(index, 1);
    };

    vm.create = function() {
      var newHost = {
        hostUUID: vm.host.uuid,
        hostname: vm.hostname.split(".")[0] + "." + vm.domain,
        descr: vm.descr,
        matchHwAddr: vm.interface.hwAddr,
        matchSerialNumber: vm.host.serialNumber,
        statusID: vm.status.id,
        roleID: vm.role.id,
        environmentID: vm.environment.id,
        siteID: vm.site.id,
        tenantID: vm.tenant.id,
        imageID: vm.image.id,
        allowBuild: vm.allowBuild,
        securityGroupID: vm.securityGroup.id
      };

      var status = vm.statuses.find(function(status) {
        return status.id == newHost.statusID;
      });
      newHost.status = status.status;

      var role = vm.roles.find(function(role) {
        return role.id == newHost.roleID;
      });
      newHost.role = role.role;

      var environment = vm.environments.find(function(environment) {
        return environment.id == newHost.environmentID;
      });
      newHost.environment = environment.environment;

      var site = vm.sites.find(function(site) {
        return site.id == newHost.siteID;
      });
      newHost.site = site.site;

      var tenant = vm.tenants.find(function(tenant) {
        return tenant.id == newHost.tenantID;
      });
      newHost.tenant = tenant.tenant;

      var image = vm.images.find(function(image) {
        return image.id == newHost.imageID;
      });
      newHost.image = image.image;

      var securityGroup = vm.securityGroups.find(function(group) {
        return group.id == newHost.securityGroupID;
      });
      newHost.securityGroup = securityGroup.group;

      var json = angular.toJson(newHost);
      data.createHost(json).then(function(data) {
        console.log(data);
        if (data.status !== undefined && data.status == 200) {
          $uibModalInstance.close("create");
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

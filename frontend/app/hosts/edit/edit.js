(function() {
  "use strict";

  angular
    .module("app")
    .controller("hostEdit", [
      "$scope",
      "$uibModalInstance",
      "data",
      "host",
      hostEdit
    ]);

  function hostEdit($scope, $uibModalInstance, data, host) {
    var vm = this;
    vm.alerts = [];
    vm.host = host;
    vm.domain = "";
    vm.hostname = host.hostname.split(".")[0];
    vm.allowBuild = host.allowBuild;

    data.getStatuses().then(function(statuses) {
      vm.statuses = statuses;
      vm.status = vm.statuses.find(function(status) {
        return status.id == vm.host.statusID;
      });
      if (vm.status !== undefined) {
        vm.status.value = vm.status.status;
      }
    });

    data.getRoles().then(function(roles) {
      vm.roles = roles;
      vm.role = vm.roles.find(function(role) {
        return role.id == vm.host.roleID;
      });
      if (vm.role !== undefined) {
        vm.role.value = vm.role.role;
      }
    });

    data.getImages().then(function(images) {
      vm.images = images;
      vm.image = vm.images.find(function(image) {
        return image.id == vm.host.imageID;
      });
      if (vm.image !== undefined) {
        vm.image.value = vm.image.image;
      }
    });

    data.getTenants().then(function(tenants) {
      vm.tenant = "";
      vm.tenants = tenants;
      vm.tenant = vm.tenants.find(function(tenant) {
        return tenant.id == vm.host.tenantID;
      });
      if (vm.tenant !== undefined) {
        vm.tenant.value = vm.tenant.tenant;
      }
    });

    data.getSecurityGroups().then(function(securityGroups) {
      vm.securityGroups = securityGroups;
      vm.securityGroup = vm.securityGroups.find(function(group) {
        return group.id == vm.host.securityGroupID;
      });
      if (vm.securityGroup !== undefined) {
        vm.securityGroup.value = vm.securityGroup.group;
      }
    });

    data.getEnvironments().then(function(environments) {
      vm.environments = environments;
      vm.environment = vm.environments.find(function(environment) {
        return environment.id == vm.host.environmentID;
      });
      if (vm.environment !== undefined) {
        vm.environment.value = vm.environment.environment;
      }
    });

    data.getSubnets().then(function(subnets) {
      vm.subnets = subnets;
    });

    data.getSites().then(function(sites) {
      vm.sites = sites;
      vm.site = vm.sites.find(function(site) {
        return site.id == vm.host.siteID;
      });
      if (vm.site !== undefined) {
        vm.site.value = vm.site.site;
      }
    });

    if (vm.host.matchHwAddr !== undefined) {
      vm.interface = vm.host.interfaces.find(function(intf) {
        return intf.hwAddr.toLowerCase() == vm.host.matchHwAddr.toLowerCase();
      });
      if (vm.interface !== undefined) {
        vm.interface.value = vm.interface.hwAddr;
      }
    }

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

    vm.edit = function() {
      var editHost = {
        hostname: vm.hostname.split(".")[0] + "." + vm.domain,
        hostUUID: vm.host.uuid,
        hostID: vm.host.hostID,
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
        return status.id == editHost.statusID;
      });
      editHost.status = status.status;

      var role = vm.roles.find(function(role) {
        return role.id == editHost.roleID;
      });
      editHost.role = role.role;

      var environment = vm.environments.find(function(environment) {
        return environment.id == editHost.environmentID;
      });
      editHost.environment = environment.environment;

      var site = vm.sites.find(function(site) {
        return site.id == editHost.siteID;
      });
      editHost.site = site.site;

      var tenant = vm.tenants.find(function(tenant) {
        return tenant.id == editHost.tenantID;
      });
      editHost.tenant = tenant.tenant;

      var image = vm.images.find(function(image) {
        return image.id == editHost.imageID;
      });
      editHost.image = image.image;

      var securityGroup = vm.securityGroups.find(function(group) {
        return group.id == editHost.securityGroupID;
      });
      editHost.securityGroup = securityGroup.group;

      var json = angular.toJson(editHost);
      data.updateHost(json).then(function(data) {
        console.log(data);
        if (data.status !== undefined && data.status == 200) {
          $uibModalInstance.close("edit");
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

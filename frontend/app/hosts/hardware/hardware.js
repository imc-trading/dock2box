(function() {
  "use strict";
  var controllerId = "hardware";

  angular
    .module("app")
    .controller("hardware", ["$uibModalInstance", "data", hardware]);

  function hardware($uibModalInstance, data) {
    var vm = this;
    vm.rows = [];
    vm.row = {};
    vm.search = "";
    vm.showSearch = false;
    vm.close = close;
    vm.getSystem = getSystem;
    vm.getCPU = getCPU;
    vm.getDNS = getDNS;
    vm.getIPMI = getIPMI;
    vm.getPCICards = getPCICards;
    vm.getDisks = getDisks;
    vm.getInterfaces = getInterfaces;
    vm.getRoutes = getRoutes;
    vm.host = data;

    getSystem();

    function getSystem() {
      vm.fields = [
        { name: "Memory", field: "memoryTotalGB", kind: "number", unit: "GB" },
        { name: "Product", field: "product", kind: "string" },
        { name: "Product Version", field: "productVersion", kind: "string" },
        { name: "Serial Number", field: "serialNumber", kind: "string" },
        { name: "BIOS Vendor", field: "biosVendor", kind: "string" },
        { name: "BIOS Date", field: "biosDate", kind: "datetime" },
        { name: "BIOS Version", field: "biosVersion", kind: "string" }
      ];

      vm.row = vm.host;
    }

    function getCPU() {
      vm.fields = [
        { name: "Model", field: "cpuModel", kind: "string" },
        { name: "Flags", field: "cpuFlags", kind: "string" },
        { name: "Speed", field: "cpuSpeedGHz", kind: "number", unit: "GHz" },
        { name: "Logical Cores", field: "cpuLogicalCores", kind: "integer" },
        { name: "Physical Cores", field: "cpuPhysicalCores", kind: "integer" },
        { name: "Sockets", field: "cpuSockets", kind: "integer" },
        {
          name: "Cores Per Socket",
          field: "cpuCoresPerSocket",
          kind: "integer"
        },
        {
          name: "Threads Per Core",
          field: "cpuThreadsPerCore",
          kind: "integer"
        }
      ];

      vm.row = vm.host;
    }

    function getDNS() {
      vm.fields = [
        { name: "Servers", field: "dnsServers", kind: "string", array: true },
        { name: "Search", field: "dnsSearch", kind: "string", array: true },
        { name: "Num. Dots", field: "dnsNDots", kind: "integer" },
        { name: "Timeout", field: "dnsTimeout", kind: "duration" },
        { name: "Attempts", field: "dnsAttempts", kind: "integer" },
        { name: "Rotate", field: "dnsRotate", kind: "boolean" }
      ];

      vm.row = vm.host;
    }

    function getIPMI() {
      vm.fields = [
        { name: "Module Present", field: "ipmiModulePresent", kind: "boolean" },
        {
          name: "Software Present",
          field: "ipmiSoftwarePresent",
          kind: "boolean"
        },
        { name: "Uses DHCP", field: "ipmiUsesDHCP", kind: "boolean" },
        { name: "IPv4 Addr.", field: "ipmiIPv4", kind: "string" },
        { name: "Netmask", field: "ipmiNetmask", kind: "string" },
        { name: "Hardware Addr.", field: "ipmiHwAddr", kind: "string" },
        { name: "Gateway", field: "ipmiGw", kind: "string" }
      ];

      vm.row = vm.host;
    }

    function getPCICards() {
      vm.fields = [
        { name: "Slot", field: "slot", kind: "string", show: true },
        { name: "Class", field: "class", kind: "string", show: true },
        { name: "Class ID", field: "classID", kind: "string" },
        { name: "Vendor", field: "vendor", kind: "string", show: true },
        { name: "Vendor ID", field: "vendorID", kind: "string" },
        { name: "Sub Vendor", field: "subVendor", kind: "string" },
        { name: "Sub Vendor ID", field: "subVendorID", kind: "string" },
        { name: "Device", field: "device", kind: "string", show: true },
        { name: "Device ID", field: "deviceID", kind: "string" },
        { name: "Revision", field: "revision", kind: "string", show: true },
        { name: "Sub Device", field: "subDevice", kind: "string", show: true },
        { name: "Sub Device ID", field: "subDeviceID", kind: "string" }
      ];

      vm.rows = vm.host.pciCards;
      vm.search = "";
    }

    function getDisks() {
      vm.fields = [
        { name: "Name", field: "name", kind: "string", show: true },
        { name: "Device", field: "device", kind: "string", show: true },
        { name: "Host", field: "host", kind: "integer", show: true },
        { name: "Channel", field: "channel", kind: "integer", show: true },
        { name: "ID", field: "id", kind: "integer", show: true },
        { name: "Lun", field: "lun", kind: "integer", show: true },
        { name: "Type", field: "hwType", kind: "string", show: true },
        { name: "Sectors", field: "sectors", kind: "integer", show: true },
        {
          name: "Sector Size",
          field: "sectorSize",
          kind: "integer",
          unit: "B",
          show: true
        },
        {
          name: "Size",
          field: "sizeGB",
          kind: "number",
          unit: "GB",
          show: true
        }
      ];

      vm.rows = vm.host.disks;
      vm.search = "";
    }

    function getInterfaces() {
      vm.fields = [
        { name: "Primary", field: "primary", kind: "boolean", show: true },
        { name: "Slot", field: "slot", kind: "string" },
        { name: "Bus", field: "bus", kind: "string" },
        { name: "Func", field: "func", kind: "string" },
        { name: "Name", field: "name", kind: "string", show: true },
        { name: "HW Addr", field: "hwAddr", kind: "string", show: true },
        { name: "MTU", field: "mtu", kind: "integer", show: true },
        { name: "IPv4", field: "ipv4", kind: "string", show: true },
        { name: "IPv6", field: "ipv6", kind: "string" },
        { name: "Netmask", field: "netmask", kind: "string", show: true },
        { name: "Network", field: "network", kind: "string", show: true },
        {
          name: "Flags",
          field: "flags",
          kind: "string",
          array: true,
          show: true
        },
        { name: "Sw Chassis ID", field: "swChassisID", kind: "string" },
        { name: "Sw Port ID", field: "swPortID", kind: "string" },
        {
          name: "Sw Port Descr",
          field: "swPortDescr",
          kind: "string",
          show: true
        },
        { name: "Sw VLAN", field: "swVLAN", kind: "string" }
      ];

      vm.rows = vm.host.interfaces;
      vm.search = "";
    }

    function getRoutes() {
      vm.fields = [
        { name: "Default", field: "default", kind: "boolean", show: true },
        {
          name: "Destination",
          field: "destination",
          kind: "string",
          show: true
        },
        { name: "Gateway", field: "gateway", kind: "string", show: true },
        { name: "Genmask", field: "genmask", kind: "string", show: true },
        { name: "Flags", field: "flags", kind: "string", show: true },
        { name: "MSS", field: "mss", kind: "integer", show: true },
        { name: "Window", field: "window", kind: "integer", show: true },
        { name: "IRTT", field: "irtt", kind: "integer", show: true },
        { name: "Interface", field: "interface", kind: "string", show: true }
      ];

      vm.rows = vm.host.routes;
      vm.search = "";
    }

    function close() {
      $uibModalInstance.close("close");
    }
  }
})();

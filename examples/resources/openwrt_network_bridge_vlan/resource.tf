resource "openwrt_network_device" "br_testing" {
  id   = "br_testing"
  name = "br-testing"
  ports = [
    "eth0",
    "eth1",
  ]
  type = "bridge"
}

resource "openwrt_network_bridge_vlan" "br_testing" {
	id = "br_testing"
	device = "br-testing"
	ports = [
		"eth0:t",
		"eth1",
	]
	vlan = 4
}
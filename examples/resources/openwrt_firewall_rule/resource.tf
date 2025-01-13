resource "openwrt_firewall_zone" "wan" {
  id       = "wantest"
  name     = "wan"
  forward  = "REJECT"
  input    = "REJECT"
  output   = "ACCEPT"
  network = [
    "wan"
  ]
}

resource "openwrt_firewall_zone" "lan" {
  id       = "lantest"
  name     = "lan"
  forward  = "ACCEPT"
  input    = "ACCEPT"
  output   = "ACCEPT"
  network = [
    "lan"
  ]
}

resource "openwrt_firewall_rule" "this" {
  id       = "testing"
  name     = "example-rule"
  src      = "wan"
  src_port = 5050
  src_ip = [
    "127.0.0.1"
  ]
  dest      = "lan"
  dest_port = 8080
  dest_ip = [
    "192.168.0.0"
  ]
  target = "DROP"
  family = "ipv4"
  proto = [
    "udp",
    "tcp",
  ]
}

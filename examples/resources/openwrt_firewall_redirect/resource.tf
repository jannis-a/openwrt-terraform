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

resource "openwrt_firewall_redirect" "this" {
  id       = "testing"
  name     = "example-rule"
  src      = openwrt_firewall_zone.wan.name
  src_dport = 8080
  dest      = openwrt_firewall_zone.lan.name
  dest_port = 8080
  dest_ip = [
    "192.168.0.0"
  ]
  target = "DNAT"
  family = "ipv4"
  proto = [
    "tcp",
  ]
}

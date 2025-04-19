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

resource "openwrt_firewall_forwarding" "this" {
  id   = "testing"
  src  = "lan"
  dest = "wan"
}

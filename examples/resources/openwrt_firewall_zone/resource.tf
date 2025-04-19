resource "openwrt_firewall_zone" "this" {
  name       = "wan"
  forward    = "REJECT"
  input      = "REJECT"
  output     = "ACCEPT"
  masquerade = true
  mssclamp   = true
  network = [
    "wan"
  ]
}

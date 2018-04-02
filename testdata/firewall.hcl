"allowed" = {
  "IPProtocol" = "tcp"

  "ports" = ["22"]
}

"allowed" = {
  "IPProtocol" = "tcp"

  "ports" = ["80"]
}

"allowed" = {
  "IPProtocol" = "tcp"

  "ports" = ["443"]
}

"allowed" = {
  "IPProtocol" = "icmp"
}

"creationTimestamp" = "2017-11-02T13:22:51.626-07:00"

"description" = ""

"direction" = "INGRESS"

"disabled" = false

"id" = "xxx"

"kind" = "compute#firewall"

"name" = "test-firewall-rule"

"network" = "projects/test-project/global/networks/network"

"priority" = 1000

"selfLink" = "projects/test-project/global/firewalls/test-firewall-rule"

"sourceRanges" = ["127.0.0.1/32"]
"canIpForward" = false

"cpuPlatform" = "Intel Broadwell"

"creationTimestamp" = "2017-08-03T22:09:15.827-07:00"

"deletionProtection" = false

"description" = "build servers"

"disks" = {
  "autoDelete" = true

  "boot" = true

  "deviceName" = "boot"

  "diskSizeGb" = "100"

  "index" = 0

  "interface" = "SCSI"

  "kind" = "compute#attachedDisk"

  "licenses" = ["projects/debian-cloud/global/licenses/debian-8-jessie"]

  "mode" = "READ_WRITE"

  "source" = "projects/test-prod/zones/us-west1-a/disks/test-vm-a"

  "type" = "PERSISTENT"
}

"id" = "4400288426170630037"

"kind" = "compute#instance"

"labelFingerprint" = "42WmSpB8rSM="

"machineType" = "projects/test-prod/zones/us-west1-a/machineTypes/n1-standard-1"

"metadata" = {
  "fingerprint" = "DI8XmRGsqlU="

  "items" = {
    "key" = "block-project-ssh-keys"

    "value" = "true"
  }

  "items" = {
    "key" = "ssh-keys"

    "value" = "bootstrap:ssh-rsa AAAAxxx bootstrap"
  }

  "kind" = "compute#metadata"
}

"name" = "test-vm-a"

"networkInterfaces" = {
  "accessConfigs" = {
    "kind" = "compute#accessConfig"

    "name" = "External NAT"

    "natIP" = "104.198.2.191"

    "networkTier" = "PREMIUM"

    "type" = "ONE_TO_ONE_NAT"
  }

  "fingerprint" = "YpX-fdWr-Yw="

  "kind" = "compute#networkInterface"

  "name" = "nic0"

  "network" = "projects/test-prod/global/networks/network"

  "networkIP" = "192.168.13.6"

  "subnetwork" = "projects/test-prod/regions/us-west1/subnetworks/subnet-13"
}

"scheduling" = {
  "automaticRestart" = true

  "onHostMaintenance" = "MIGRATE"

  "preemptible" = false
}

"selfLink" = "projects/test-prod/zones/us-west1-a/instances/test-vm-a"

"serviceAccounts" = {
  "email" = "bogus-compute@developer.gserviceaccount.com"

  "scopes" = ["https://www.googleapis.com/auth/compute.readonly", "https://www.googleapis.com/auth/cloud.useraccounts.readonly", "https://www.googleapis.com/auth/devstorage.read_only", "https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/monitoring.write"]
}

"startRestricted" = false

"status" = "RUNNING"

"tags" = {
  "fingerprint" = "bogus="

  "items" = ["testvm"]
}

"zone" = "projects/test-prod/zones/us-west1-a"
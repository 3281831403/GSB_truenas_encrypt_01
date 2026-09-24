# Create a pool from scratch.
#
# Identify disks by their STABLE SERIAL, not the kernel device name (sdX).
# Kernel names are reassigned across reboots and controller changes; the
# provider resolves any accepted form (serial, TrueNAS identifier,
# /dev/disk/by-id path, or sdX) to the same physical disk at plan time, so a
# config that pins disks by serial never plans a pool replacement when the
# device names shuffle. Discover serials with
# `midclt call disk.query '[]' '{"select": ["name", "serial"]}'` on the box.
#
# `topology` is a nested attribute, so it takes the `= { ... }` assignment
# form (not a bare `topology { ... }` block). Use "MIRROR"/"RAIDZ2"/etc.;
# for a single-disk vdev use "STRIPE".
#
# name and topology are effectively immutable: the provider refuses a
# post-create change to them (it will not destroy and recreate a pool). Adding
# prevent_destroy is recommended as defense-in-depth so a `terraform destroy`
# or `-replace` cannot erase the pool's data by accident either.
resource "truenas_pool" "tank" {
  name     = "tank"
  autotrim = false
  topology = {
    data = [
      { type = "MIRROR", disks = ["WD-WCC7K5PACL0V", "WD-WCC7K6ABXYZ1"] },
      { type = "MIRROR", disks = ["WD-WCC7K7CDEFG2", "WD-WCC7K8HIJKL3"] },
    ]
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Datasets managed alongside the pool should build their name from the pool
# resource, so Terraform creates the pool before its datasets rather than
# racing them. See the truenas_dataset example for the full pattern.
resource "truenas_dataset" "data" {
  name = "${truenas_pool.tank.name}/data"
}

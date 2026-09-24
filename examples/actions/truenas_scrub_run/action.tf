# Kick off a pool scrub (Terraform 1.14+).
# pool_id is the numeric pool id. Look it up by pool name with the
# truenas_pool data source rather than hardcoding it.
data "truenas_pool" "tank" {
  name = "tank"
}

action "truenas_scrub_run" "example" {
  config {
    pool_id = data.truenas_pool.tank.id

    # By default the action blocks until the scrub job finishes (which can be
    # a long time). Set wait = false to start it and return immediately.
    # wait = false
  }
}

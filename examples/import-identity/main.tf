# Identity-based import (Terraform v1.14+): the `identity` attribute takes
# the values shown under "Identity Schema" on the resource's documentation
# page (or found via `terraform query`, see examples/list) instead of the
# single opaque `id` string used by classic `import { id = "..." }` blocks.

import {
  to       = truenas_user.example
  identity = { id = 1 }
}

resource "truenas_user" "example" {
  username  = "deploy"
  full_name = "Deploy User"
  password  = "changeme"
  shell     = "/bin/bash"
}

# Live-smoke testing for actions

The `.tf` files in each subdirectory here are the example blocks embedded in
the generated docs (`docs/actions/*.md`) — they show the shape of an
`action` block's `config {}`, but an `action` block alone does nothing.
Terraform only *invokes* an action when something triggers it: a
`lifecycle { action_trigger }` block on a resource, evaluated during
`terraform apply`.

This file is about that second half: how to point a real provider build at a
real TrueNAS box and confirm an action actually fires. It's manual, not CI —
see "Why this isn't automated" below.

## Requirements

- Terraform **1.14 or newer**. Actions did not exist as an invocable language
  construct before 1.14; older Terraform will reject the `action` block and
  `action_trigger` with a parse error.
- A local build of the provider and a `dev_overrides` entry (see the main
  [README](../../README.md#build-from-source-optional) for the full setup):

  ```hcl
  provider_installation {
    dev_overrides {
      "truenas/truenas" = "/path/to/terraform-provider-truenas"
    }
    direct {}
  }
  ```

  With `dev_overrides` active, `terraform init` is skipped — Terraform loads
  whatever `make build` last produced.
- A reachable TrueNAS box and credentials, same as any other manual test
  (see `TRUENAS_ENDPOINT` / `TRUENAS_API_KEY` in the main README).

## Wiring an action to fire

An action only runs when a resource's `lifecycle` block references it in an
`action_trigger`. The simplest harness resource is `terraform_data` (the
built-in null-resource replacement) with `triggers_replace` bumped by hand
each time you want to re-fire the action:

```hcl
terraform {
  required_providers {
    truenas = {
      source  = "truenas/truenas"
      version = "~> 1.0"
    }
  }
}

resource "terraform_data" "smoke" {
  triggers_replace = [timestamp()]

  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.truenas_ui_restart.smoke]
    }
  }
}

action "truenas_ui_restart" "smoke" {
  config {
    delay = 5
  }
}
```

`terraform apply` will show the action queued alongside the resource change;
confirm at the `apply` prompt to actually invoke it. Swap the `action` block
for `truenas_service_control`, `truenas_scrub_run`, etc. and update
`actions = [...]` to match.

## Low-risk actions to try first

Try these against a lab/test box (not production) in roughly this order:

1. **`truenas_ui_restart`** with `delay = 5` or more — restarts the web UI.
   Low risk: the UI comes back on its own, and the delay gives you time to
   watch it drop and recover.
2. **`truenas_service_control`** with `verb = "RESTART"` on a service that
   isn't serving anything important on the test box (e.g. `ssh` if you have
   console access, or `cron`) — confirms the version gate (26.0+) and the
   verb enum both behave.
3. **`truenas_scrub_run`** against a small test pool — confirms job
   invocation and that the action returns once the job starts rather than
   blocking on completion.

Avoid `truenas_replication_run`, `truenas_cloudsync_run`,
`truenas_snapshot_task_run`, and the `truenas_app_*` actions for smoke
testing unless you have a disposable task/app set up — they act on
whatever `id` or `app_name` you give them.

## Why this isn't automated

`terraform-plugin-testing` v1.16 has no action `TestStep` — the acceptance
test harness this provider otherwise relies on has no way to declare an
`action_trigger` or assert an action fired. Each action's request-building
and schema logic (including `service_control`'s version gate) is covered by
unit tests under `internal/actions/`; end-to-end invocation against a real
box is verified manually with the steps above until plugin-testing gains
action support.

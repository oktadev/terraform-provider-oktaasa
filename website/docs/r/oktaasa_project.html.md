---
layout: "oktaasa"
page_title: "Advanced Server Access: oktaasa_project"
sidebar_current: "docs-resource-oktaasa-project"
description: |-
  The oktaasa_project resource creates projects in Okta's ASA.
---

# oktaasa\_project

The oktaasa_project resource creates projects in Okta's ASA.

## Example Usage

```hcl
resource "oktaasa_project" "demo-stack" {
  project_name = "tf-test"
  next_unix_uid = 60120
  next_unix_gid = 63020
}
```


## Argument Reference

The following arguments are supported:

* `project_name` (Required) - name of the project.
* `next_unix_uid` (Optional - Default: 60101) - Okta's ASA will start assigning Unix user IDs from this value
* `next_unix_gid` (Optional - Default: 63001) - Okta's ASA will start assigning Unix group IDs from this value


## Attributes Reference

No further attributes are exported.



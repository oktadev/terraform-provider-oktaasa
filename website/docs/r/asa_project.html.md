---
layout: "asa"
page_title: "Advanced Server Access: asa_project"
sidebar_current: "docs-resource-asa-project"
description: |-
  The asa_project resource creates projects in ASA.
---

# asa\_project

The asa_project resource creates projects in ASA.

## Example Usage

```hcl
resource "asa_project" "demo-stack" {
  project_name = "tf-test"
  next_unix_uid = 60120
  next_unix_gid = 63020
}
```


## Argument Reference

The following arguments are supported:

* `project_name` (Required) - name of the project.
* `next_unix_uid` (Optional - Default: 60101) - ASA will start assigning Unix user IDs from this value
* `next_unix_gid` (Optional - Default: 63001) - ASA will start assigning Unix group IDs from this value


## Attributes Reference

No further attributes are exported.



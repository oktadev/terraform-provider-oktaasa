---
layout: "oktaasa"
page_title: "Advanced Server Access: oktaasa_assign_group"
sidebar_current: "docs-resource-oktaasa-assign-group"
description: |-
  The oktaasa_assign_group resource assigns control access levels in group definitions in Okta's ASA.
---

# oktaasa\_assign\_group

The oktaasa_assign_group resource assigns control access levels in group definitions in Okta's ASA.  In order to give access to project, you need to assign an Okta group to a project. Use "server_access" and "server_admin" parameters to control access level.

## Example Usage

```hcl
resource "oktaasa_assign_group" "sg-cloud-group-access" {
  project_name = "tf-test"
  group_name = "cloud-ro"
  server_access = true
  server_admin = false
  create_server_group = true
}
```


## Argument Reference

The following arguments are supported:

* `project_name` (Required) - name of the project.
* `group_name` (Required) - name of the group 
* `server_access` (bool) (Required) - Whether users in this group have access permissions on the servers in this project.
* `server_admin` (bool) (Required) - Whether users in this group have sudo permissions on the servers in this project.
* `create_server_group` (bool) (Optional - Default: true) - will make Okta's ASA synchronize group name to linux box. To avoid naming collision, group created by Okta's ASA will have prefix of "oktaasa_"


## Attributes Reference

No further attributes are exported.



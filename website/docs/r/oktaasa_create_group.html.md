---
layout: "oktaasa"
page_title: "Advanced Server Access: oktaasa_create_group"
sidebar_current: "docs-resource-oktaasa-create-group"
description: |-
  The oktaasa_create_group resource creates projects in Okta's ASA.
---

# oktaasa\_create\_group

The oktaasa_create_group resource creates groups in Okta's ASA.  If groups are not synced from Okta, you may need to create groups in Okta's ASA using this resource.

NOTE: group is created with basic access_user access. It does not give any privileges in Okta's ASA console. Creation of groups with access_admin and reporting_user is currently not supported in the provider.


## Example Usage

```hcl
resource "oktaasa_create_group" "test-tf-group" {
  name = "test-tf-group"
}
```


## Argument Reference

The following arguments are supported:

* `name` (Required) - name for Okta's ASA group.


## Attributes Reference

No further attributes are exported.



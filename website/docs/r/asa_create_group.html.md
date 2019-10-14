---
layout: "asa"
page_title: "Advanced Server Access: asa_create_group"
sidebar_current: "docs-resource-asa-create-group"
description: |-
  The asa_create_group resource creates projects in ASA.
---

# asa\_create\_group

The asa_create_group resource creates groups in ASA.  If groups are not synced from Okta, you may need to create groups in ASA using this resource.

NOTE: group is created with basic access_user access. It does not give any privileges in ASA console. Creation of groups with access_admin and reporting_user is currently not supported in the provider.


## Example Usage

```hcl
resource "asa_create_group" "test-tf-group" {
  name = "test-tf-group"
}
```


## Argument Reference

The following arguments are supported:

* `name` (Required) - name for the ASA group.


## Attributes Reference

No further attributes are exported.



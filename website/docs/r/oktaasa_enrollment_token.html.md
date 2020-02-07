---
layout: "oktaasa"
page_title: "Advanced Server Access: oktaasa_enrollment_token"
sidebar_current: "docs-resource-oktaasa-enrollment-token"
description: |-
  The oktaasa_token resource creates enrollment tokens which are base64-encoded objects with metadata that Okta's ASA Agent can configure itself from.  Enrollment is the process where Okta's ASA agent configures a server to be managed by a specific project.

---

# oktaasa\_token

The oktaasa_token resource creates enrollment tokens which are base64-encoded objects with metadata that Okta's ASA Agent can configure itself from.  Enrollment is the process where Okta's ASA agent configures a server to be managed by a specific project.

## Example Usage

```hcl
resource "oktaasa_token" "stack-x-token" {
  project_name = "tf-test"
  description = "Token for X stack"
}
```


## Argument Reference

The following arguments are supported:

* `project_name` (Required) - name of the project.
* `Description` (Required) - free form text field to provide description. You will NOT be able to change it later without recreating a token.


## Attributes Reference

No further attributes are exported.



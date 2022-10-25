---
layout: "oktaasa"
page_title: "Provider: Okta Advanced Server Access"
sidebar_current: "docs-oktaasa-index"
description: |-
  Okta's Advanced Server Access (ASA) provider configures projects, groups, server enrollment tokens and group permissions in Advanced Server Access.
---

# DEPRECATION NOTICE

This provider is now deprecated. Okta has released a replacement provider https://registry.terraform.io/providers/okta/oktapam/latest that contains all functionality provided by this provider as well as support for additional resource types. Please migrate to the OktaPAM Terraform provider.

NOTE: The OktaPAM Terraform Provider is not backwards compatible with this provider. However, the new provider supports data sources which will allow users to import current ASA configuration into Terraform for use with the new provider.

# Okta Advanced Server Access Provider

Okta's Advanced Server Access (ASA) provider configures key parameters such as projects, groups, server enrollment tokens and group permissions in [Advanced Server Access](https://www.okta.com/products/advanced-server-access/), which provides zero trust access management for infrastructure that extends Okta’s core platform to Linux and Windows servers via SSH and RDP.  It does so in a manner that replaces static keys with a more elegant approach based on an ephemeral client certificate architecture.

These parameters are defined as:
* Project - project is an authorization scope. It associates a collection of resources with a set of configurations, including RBAC and access policies.
* Enrollment token - using enrollment token, you can add servers to a project.
* Groups - in order to give access to a project, you need to create, configure permissions and assign a group to it.


The provider requires these environment variables to be set with your respective values:

* export OKTAASA_KEY="Your OKTAASA key"
* export OKTAASA_KEY_SECRET="Your OKTAASA key secret"
* export OKTAASA_TEAM="Your OKTAASA team name"

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
resource "oktaasa_project" "demo-project" {
  project_name = "tf-test"
  next_unix_uid = 60120
  next_unix_gid = 63020
}

resource "oktaasa_enrollment_token" "test-token" {
  project_name = "${oktaasa_project.demo-project.project_name}"
  description  = "Token for X"
}

resource "oktaasa_enrollment_token" "test-import-token" {
 project_name = "${oktaasa_project.demo-project.project_name}"
  description  = "Token for Y"
}

resource "oktaasa_create_group" "test-tf-group" {
  name = "test-tf-group"
}

resource "oktaasa_create_group" "cloud-sre" {
  name = "cloud-sre"
}

resource "oktaasa_create_group" "cloud-support" {
  name = "cloud-support"
}

resource "oktaasa_assign_group" "test-sft-group-assignment" {
  project_name = "${oktaasa_project.demo-project.project_name}"
  group_name   = "${oktaasa_create_group.test-tf-group.name}"
}

resource "oktaasa_assign_group" "group-assignment" {
  project_name        = "${oktaasa_project.demo-project.project_name}"
  group_name          = "cloud-sre"
  server_access       = true
  server_admin        = true
  create_server_group = false
}

resource "oktaasa_assign_group" "dev-group-assignment" {
  project_name  = "${oktaasa_project.demo-project.project_name}"
  group_name    = "cloud-support"
  server_access = true
  server_admin  = false
}

output "enrollment_token" {
  value = "${oktaasa_enrollment_token.test-token.token_value}"
}
```

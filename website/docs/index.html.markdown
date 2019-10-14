---
layout: "asa"
page_title: "Provider: Advanced Server Access"
sidebar_current: "docs-asa-index"
description: |-
  The Advanced Server Access (ASA) provider configures projects, groups, server enrollment tokens and group permissions in Advanced Server Access.
---

# Advanced Server Access Provider

The Advanced Server Access (ASA) provider configures key parameters such as projects, groups, server enrollment tokens and group permissions in [Advanced Server Access](https://www.okta.com/products/advanced-server-access/), which provides zero trust access management for infrastructure that extends Okta’s core platform to Linux and Windows servers via SSH and RDP.  It does so in a manner that replaces static keys with a more elegant approach based on an ephemeral client certificate architecture.

These parameters are defined as:
* Project - project is an authorization scope. It associates a collection of resources with a set of configurations, including RBAC and access policies.
* Enrollment token - using enrollment token, you can add servers to a project.
* Groups - in order to give access to a project, you need to create, configure permissions and assign a group to it.


The provider requires these environment variables to be set with your respective values:

* export ASA_KEY="Your ASA key"
* export ASA_KEY_SECRET="Your ASA key secret"
* export ASA_TEAM="Your ASA team name"

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
resource "asa_project" "demo-project" {
  project_name = "tf-test"
  next_unix_uid = 60120
  next_unix_gid = 63020
}

resource "asa_enrollment_token" "test-token" {
  project_name = "${asa_project.demo-project.project_name}"
  description  = "Token for X"
}

resource "asa_enrollment_token" "test-import-token" {
 project_name = "${asa_project.demo-project.project_name}"
  description  = "Token for Y"
}

resource "asa_create_group" "test-tf-group" {
  name = "test-tf-group"
}

resource "asa_create_group" "cloud-sre" {
  name = "cloud-sre"
}

resource "asa_create_group" "cloud-support" {
  name = "cloud-support"
}

resource "asa_assign_group" "test-sft-group-assignment" {
  project_name = "${asa_project.demo-project.project_name}"
  group_name   = "${asa_create_group.test-tf-group.name}"
}

resource "asa_assign_group" "group-assignment" {
  project_name        = "${asa_project.demo-project.project_name}"
  group_name          = "cloud-sre"
  server_access       = true
  server_admin        = true
  create_server_group = false
}

resource "asa_assign_group" "dev-group-assignment" {
  project_name  = "${asa_project.demo-project.project_name}"
  group_name    = "cloud-support"
  server_access = true
  server_admin  = false
}

output "enrollment_token" {
  value = "${asa_enrollment_token.test-token.token_value}"
}
```

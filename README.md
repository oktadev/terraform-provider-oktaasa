Terraform Provider for Advanced Server Access (ASA) 
==================

Maintainers
-----------

This provider plugin is maintained by the Hashicorp Terraform and Okta teams. 
Credit to Aleksei Denisov and Splunk Team for providing original source.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html)
-	[Go](https://golang.org/doc/install)

Build
---------------------

```
go build -o terraform-provider-asa
```

Usage
---------------------

This plugin requires a couple of inputs to run: ASA API key, secret and team name. Set them as environment variables.
```
export ASA_KEY_SECRET=<secret here>
export ASA_KEY=<key here>
export ASA_TEAM=<team name>
```

Using the provider
----------------------
There are three main resources that you need to provision:
1. Project - project is an authorization scope. It associates a collection of resources with a set of configurations, including RBAC and access policies.
1. Enrollment token - using enrollment token, you can add servers to a project.
1. Groups - in order to give access to a project, you need to assign a group to it.


Sample terraform plan:

```
resource "asa_project" "demo-stack" {
  project_name = "tf-test"
}

resource "asa_enrollment_token" "test-token" {
  project_name = "${asa_project.demo-stack.project_name}"
  description = "Token for X"
}

resource "asa_enrollment_token" "test-import-token" {
  project_name = "${asa_project.demo-stack.project_name}"
  description = "Token for Y"
}

resource "asa_assign_group" "group-assignment" {
  project_name = "${asa_project.demo-stack.project_name}"
  group_name = "cloud-sre"
  server_access = true
  server_admin = true
  create_server_group = false
}

resource "asa_assign_group" "dev-group-assignment" {
  project_name = "${asa_project.demo-stack.project_name}"
  group_name = "cloud-support"
  server_access = true
  server_admin = true
}


```

## Resources
### Project

Example usage:
```
resource "asa_project" "demo-stack" {
  project_name = "tf-test"
}
```
Parameters:
* project_name (Required) - name of the project.
* next_unix_uid (Optional - Default: 60101) - ASA will start assigning Unix user IDs from this value
* next_unix_gid (Optional - Default: 63001) - ASA will start assigning Unix group IDs from this value

### Enrollment token
Enrollment is the process where the ASA agent configures a server to be managed by a specific project. An enrollment token is a base64 encoded object with metadata that the ASA Agent can configure itself from.  

Example usage:
```
resource "asa_token" "stack-x-token" {
  project_name = "tf-test"
  description = "Token for X stack"
}
```
Parameters:
* project_name (Required) - name of the project.
* Description (Required) - free form text field to provide description. You will NOT be able to change it later without recreating a token.

### Create Group
If groups is not synced from Okta, you may need to create in ASA.
```
resource "asa_create_group" "test-tf-group" {
  name = "test-tf-group"
}
```
Parameters:
* name (Required) - name for the ASA group.

NOTE: group is created with basic access_user access. It does not give any privileges in ASA console.
Creation of groups with access_admin and reporting_user is currently not supported in the provider.

### Assign group to project
In order to give access to project, you need to assign Okta group to a project. Use "server_access" and "server_admin" parameters to control access level.

Example usage:
```
resource "asa_assign_group" "sg-cloud-group-access" {
  project_name = "tf-test"
  group_name = "cloud-ro"
  server_access = true
  server_admin = false
  create_server_group = true
}
```
Parameters:
* project_name (Required) - name of the project.
* server_access (bool) (Required) - Whether users in this group have access permissions on the servers
in this project.
* server_admin (bool) (Required) - Whether users in this group have sudo permissions on the servers in this project.
* create_server_group (bool) (Optional - Default: true) - will make ASA synchronize group name to linux box. To avoid naming collision, group created by ASA will have prefix of "asa_"

Remove/destroy configured parameters
----------------------
To remove configurations that were created with Terraform, do the following.  Refer to the main.tf for the defined resources and review the "resource_type" and "resource_name" values.  

```
resource "resource_type" "resource_name" {
  ...
}
```

Or, at a commandline prompt where you have sourced the terraform binaries, enter "terraform state list" (without the quotes).  You should get a list like the following.  

```
[root@ip-172-31-20-120 terraform-provider-asa]# terraform state list
asa_assign_group.dev-group-assignment
asa_assign_group.group-assignment
asa_assign_group.test-sft-group-assignment
asa_create_group.cloud-sre
asa_create_group.cloud-support
asa_create_group.test-tf-group
asa_enrollment_token.test-import-token
asa_enrollment_token.test-token
asa_project.demo-project
```

Compile the terraform destroy command with the -target arguments in the following format:

terraform destroy -target RESOURCE_TYPE.NAME -target RESOURCE_TYPE2.NAME -target RESOURCE_TYPE3.NAME

An example would be like the following.  Note that destroying the demo-project also destroys certain parameters like the tokens.  However, the groups still must be destroyed.

```
terraform destroy -target asa_project.demo-project -target asa_create_group.cloud-sre -target asa_create_group.cloud-support -target asa_create_group.test-tf-group
```

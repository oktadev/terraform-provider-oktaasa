Terraform Provider for Okta's Advanced Server Access (Okta's ASA)
=========================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.13+ (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-oktaasa`

```sh
$ git clone git@github.com:terraform-providers/terraform-provider-oktaasa $GOPATH/src/github.com/terraform-providers/terraform-provider-oktaasa
```

Go to the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-oktaasa
$ make build
```

Using the provider
----------------------
Set the following environment variables prior to running: Okta's ASA API key, secret and team name. 
```
export OKTAASA_KEY_SECRET=<secret here>
export OKTAASA_KEY=<key here>
export OKTAASA_TEAM=<team name>
```

There are three main resources that you need to provision:
1. Project - project is an authorization scope. It associates a collection of resources with a set of configurations, including RBAC and access policies.
1. Enrollment token - using enrollment token, you can add servers to a project.
1. Groups - in order to give access to a project, you need to assign a group to it.


Sample terraform plan:

```
resource "oktaasa_project" "demo-stack" {
  project_name = "tf-test"
}

resource "oktaasa_enrollment_token" "test-token" {
  project_name = "${oktaasa_project.demo-stack.project_name}"
  description = "Token for X"
}

resource "oktaasa_enrollment_token" "test-import-token" {
  project_name = "${oktaasa_project.demo-stack.project_name}"
  description = "Token for Y"
}

resource "oktaasa_assign_group" "group-assignment" {
  project_name = "${oktaasa_project.demo-stack.project_name}"
  group_name = "cloud-sre"
  server_access = true
  server_admin = true
  create_server_group = false
}

resource "oktaasa_assign_group" "dev-group-assignment" {
  project_name = "${oktaasa_project.demo-stack.project_name}"
  group_name = "cloud-support"
  server_access = true
  server_admin = true
}


```

## Resources
### Project

Example usage:
```
resource "oktaasa_project" "demo-stack" {
  project_name = "tf-test"
}
```
Parameters:
* project_name (Required) - name of the project.
* next_unix_uid (Optional - Default: 60101) - Okta's ASA will start assigning Unix user IDs from this value
* next_unix_gid (Optional - Default: 63001) - Okta's ASA will start assigning Unix group IDs from this value

### Enrollment token
Enrollment is the process where Okta's ASA agent configures a server to be managed by a specific project. An enrollment token is a base64 encoded object with metadata that Okta's ASA Agent can configure itself from.  

Example usage:
```
resource "oktaasa_enrollment_token" "stack-x-token" {
  project_name = "tf-test"
  description = "Token for X stack"
}
```
Parameters:
* project_name (Required) - name of the project.
* Description (Required) - free form text field to provide description. You will NOT be able to change it later without recreating a token.

### Create Group
If groups is not synced from Okta, you may need to create in Okta's ASA.
```
resource "oktaasa_create_group" "test-tf-group" {
  name = "test-tf-group"
}
```
Parameters:
* name (Required) - name for Okta's ASA group.

NOTE: group is created with basic access_user access. It does not give any privileges in Okta's ASA console.
Creation of groups with access_admin and reporting_user is currently not supported in the provider.

### Assign group to project
In order to give access to project, you need to assign Okta group to a project. Use "server_access" and "server_admin" parameters to control access level.

Example usage:
```
resource "oktaasa_assign_group" "sg-cloud-group-access" {
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
* create_server_group (bool) (Optional - Default: true) - will make Okta's ASA synchronize group name to linux box. To avoid naming collision, group created by Okta's ASA will have prefix of "oktaasa_"

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
[root@ip-172-31-20-120 terraform-provider-oktaasa]# terraform state list
oktaasa_assign_group.dev-group-assignment
oktaasa_assign_group.group-assignment
oktaasa_assign_group.test-sft-group-assignment
oktaasa_create_group.cloud-sre
oktaasa_create_group.cloud-support
oktaasa_create_group.test-tf-group
oktaasa_enrollment_token.test-import-token
oktaasa_enrollment_token.test-token
oktaasa_project.demo-project
```

Compile the terraform destroy command with the -target arguments in the following format:

terraform destroy -target RESOURCE_TYPE.NAME -target RESOURCE_TYPE2.NAME -target RESOURCE_TYPE3.NAME

An example would be like the following.  Note that destroying the demo-project also destroys certain parameters like the tokens.  However, the groups still must be destroyed.

```
terraform destroy -target oktaasa_project.demo-project -target oktaasa_create_group.cloud-sre -target oktaasa_create_group.cloud-support -target oktaasa_create_group.test-tf-group
```


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-oktaasa
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

**Special thanks to Aleksei Denisov and the Splunk Cloud team members who were the authors of the original provider which has since been re-purposed for this certified version**

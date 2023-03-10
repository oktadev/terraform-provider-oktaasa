resource "oktaasa_project" "demo-project" {
  project_name = "tf-test"
  next_unix_uid = 60120
  next_unix_gid = 63020
  require_preathorization = false
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

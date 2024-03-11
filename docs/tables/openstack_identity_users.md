# Table: openstack_identity_users

This table shows data for Openstack Identity Users.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_identity_users:
  - [openstack_identity_user_keypairs](openstack_identity_user_keypairs.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|ignore_change_password_upon_first_use|`bool`|
|ignore_lockout_failure_attempts|`bool`|
|ignore_password_expiry|`bool`|
|default_project_id|`utf8`|
|description|`utf8`|
|domain_id|`utf8`|
|enabled|`bool`|
|id|`utf8`|
|name|`utf8`|
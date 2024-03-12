# Table: openstack_identity_projects

This table shows data for Openstack Identity Projects.

The primary key for this table is **id**.

## Relations

The following tables depend on openstack_identity_projects:
  - [openstack_compute_project_limits](openstack_compute_project_limits.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id|`uuid`|
|_cq_parent_id|`uuid`|
|is_domain|`bool`|
|description|`utf8`|
|domain_id|`utf8`|
|enabled|`bool`|
|id (PK)|`utf8`|
|name|`utf8`|
|parent_id|`utf8`|
|tags|`list<item: utf8, nullable>`|
|options|`json`|
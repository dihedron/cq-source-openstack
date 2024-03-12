# Table: openstack_compute_project_limits

This table shows data for Openstack Compute Project Limits.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_identity_projects](openstack_identity_projects.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|max_total_cores|`int64`|
|max_image_meta|`int64`|
|max_server_meta|`int64`|
|max_personality|`int64`|
|max_personality_size|`int64`|
|max_total_keypairs|`int64`|
|max_security_groups|`int64`|
|max_security_group_rules|`int64`|
|max_server_groups|`int64`|
|max_server_group_members|`int64`|
|max_total_floating_ips|`int64`|
|max_total_instances|`int64`|
|max_total_ram_size|`int64`|
|total_cores_used|`int64`|
|total_instances_used|`int64`|
|total_floating_ips_used|`int64`|
|total_ram_used|`int64`|
|total_security_groups_used|`int64`|
|total_server_groups_used|`int64`|
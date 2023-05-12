# Table: openstack_security_groups

This table shows data for Openstack Security Groups.

The primary key for this table is **id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|security_group_rule_ids|StringArray|
|id (PK)|String|
|name|String|
|description|String|
|security_group_rules|JSON|
|tenant_id|String|
|project_id|String|
|tags|StringArray|
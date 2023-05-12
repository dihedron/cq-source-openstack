# Table: openstack_security_group_rules

This table shows data for Openstack Security Group Rules.

The primary key for this table is **id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|id (PK)|String|
|direction|String|
|description|String|
|ethertype|String|
|security_group_id|String|
|port_range_min|Int|
|port_range_max|Int|
|protocol|String|
|remote_group_id|String|
|remote_ip_prefix|String|
|tenant_id|String|
|project_id|String|
# Table: openstack_networks

This table shows data for Openstack Networks.

The primary key for this table is **id**.

## Relations

The following tables depend on openstack_networks:
  - [openstack_network_subnets](openstack_network_subnets.md)
  - [openstack_network_tags](openstack_network_tags.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|id (PK)|String|
|name|String|
|description|String|
|admin_state_up|Bool|
|status|String|
|subnets|StringArray|
|tenant_id|String|
|project_id|String|
|shared|Bool|
|availability_zone_hints|StringArray|
|tags|StringArray|
|revision_number|Int|
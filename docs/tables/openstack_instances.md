# Table: openstack_instances

This table shows data for Openstack Instances.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_instances:
  - [openstack_instance_addresses](openstack_instance_addresses.md)
  - [openstack_instance_attached_volumes](openstack_instance_attached_volumes.md)
  - [openstack_instance_flavor_extra_specs](openstack_instance_flavor_extra_specs.md)
  - [openstack_instance_flavors](openstack_instance_flavors.md)
  - [openstack_instance_metadata](openstack_instance_metadata.md)
  - [openstack_instance_security_groups](openstack_instance_security_groups.md)
  - [openstack_instance_tags](openstack_instance_tags.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|image_id|String|
|power_state_name|String|
|id|String|
|tenant_id|String|
|user_id|String|
|name|String|
|created_at|Timestamp|
|launched_at|Timestamp|
|updated_at|Timestamp|
|terminated_at|Timestamp|
|hostid|String|
|status|String|
|progress|Int|
|access_ipv4|String|
|access_ipv6|String|
|flavor|JSON|
|addresses|JSON|
|metadata|JSON|
|key_name|String|
|admin_pass|String|
|security_groups|JSON|
|attached_volumes|JSON|
|tags|StringArray|
|server_groups|StringArray|
|disk_config|String|
|availability_zone|String|
|host|String|
|hostname|String|
|hypervisor_hostname|String|
|instance_name|String|
|kernel_id|String|
|launch_index|Int|
|ramdisk_id|String|
|reservation_id|String|
|root_device_name|String|
|user_data|String|
|power_state_id|Int|
|vm_state|String|
|config_drive|String|
|description|String|
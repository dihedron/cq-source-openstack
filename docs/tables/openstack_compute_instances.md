# Table: openstack_compute_instances

This table shows data for Openstack Compute Instances.

The primary key for this table is **_cq_id**.

## Relations

The following tables depend on openstack_compute_instances:
  - [openstack_compute_instance_addresses](openstack_compute_instance_addresses.md)
  - [openstack_compute_instance_attached_volumes](openstack_compute_instance_attached_volumes.md)
  - [openstack_compute_instance_flavor_extra_specs](openstack_compute_instance_flavor_extra_specs.md)
  - [openstack_compute_instance_flavors](openstack_compute_instance_flavors.md)
  - [openstack_compute_instance_metadata](openstack_compute_instance_metadata.md)
  - [openstack_compute_instance_security_groups](openstack_compute_instance_security_groups.md)
  - [openstack_compute_instance_tags](openstack_compute_instance_tags.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|image_id|`utf8`|
|power_state_name|`utf8`|
|id|`utf8`|
|tenant_id|`utf8`|
|user_id|`utf8`|
|name|`utf8`|
|created_at|`timestamp[us, tz=UTC]`|
|launched_at|`timestamp[us, tz=UTC]`|
|updated_at|`timestamp[us, tz=UTC]`|
|terminated_at|`timestamp[us, tz=UTC]`|
|hostid|`utf8`|
|status|`utf8`|
|progress|`int64`|
|access_ipv4|`utf8`|
|access_ipv6|`utf8`|
|flavor|`json`|
|addresses|`json`|
|metadata|`json`|
|key_name|`utf8`|
|admin_pass|`utf8`|
|security_groups|`json`|
|attached_volumes|`json`|
|tags|`list<item: utf8, nullable>`|
|server_groups|`list<item: utf8, nullable>`|
|disk_config|`utf8`|
|availability_zone|`utf8`|
|host|`utf8`|
|hostname|`utf8`|
|hypervisor_hostname|`utf8`|
|instance_name|`utf8`|
|kernel_id|`utf8`|
|launch_index|`int64`|
|ramdisk_id|`utf8`|
|reservation_id|`utf8`|
|root_device_name|`utf8`|
|user_data|`utf8`|
|power_state_id|`int64`|
|vm_state|`utf8`|
|config_drive|`utf8`|
|description|`utf8`|
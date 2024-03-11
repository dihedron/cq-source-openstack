# Table: openstack_baremetal_nodes

This table shows data for Openstack Baremetal Nodes.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|automated_clean|`bool`|
|uuid|`utf8`|
|name|`utf8`|
|power_state|`utf8`|
|target_power_state|`utf8`|
|provision_state|`utf8`|
|target_provision_state|`utf8`|
|maintenance|`bool`|
|maintenance_reason|`utf8`|
|fault|`utf8`|
|last_error|`utf8`|
|reservation|`utf8`|
|driver|`utf8`|
|driver_info|`json`|
|driver_internal_info|`json`|
|instance_info|`json`|
|instance_uuid|`utf8`|
|chassis_uuid|`utf8`|
|extra|`json`|
|console_enabled|`bool`|
|raid_config|`json`|
|target_raid_config|`json`|
|clean_step|`json`|
|deploy_step|`json`|
|resource_class|`utf8`|
|bios_interface|`utf8`|
|boot_interface|`utf8`|
|console_interface|`utf8`|
|deploy_interface|`utf8`|
|inspect_interface|`utf8`|
|management_interface|`utf8`|
|network_interface|`utf8`|
|power_interface|`utf8`|
|raid_interface|`utf8`|
|rescue_interface|`utf8`|
|storage_interface|`utf8`|
|traits|`list<item: utf8, nullable>`|
|vendor_interface|`utf8`|
|conductor_group|`utf8`|
|protected|`bool`|
|protected_reason|`utf8`|
|owner|`utf8`|
|network_data|`json`|
|created_at|`timestamp[us, tz=UTC]`|
|updated_at|`timestamp[us, tz=UTC]`|
|provision_updated_at|`timestamp[us, tz=UTC]`|
|inspection_started_at|`timestamp[us, tz=UTC]`|
|inspection_finished_at|`timestamp[us, tz=UTC]`|
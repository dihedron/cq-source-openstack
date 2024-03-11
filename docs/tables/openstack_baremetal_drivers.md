# Table: openstack_baremetal_drivers

This table shows data for Openstack Baremetal Drivers.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|name|`utf8`|
|hosts|`list<item: utf8, nullable>`|
|type|`utf8`|
|default_bios_interface|`utf8`|
|default_boot_interface|`utf8`|
|default_console_interface|`utf8`|
|default_deploy_interface|`utf8`|
|default_inspect_interface|`utf8`|
|default_management_interface|`utf8`|
|default_network_interface|`utf8`|
|default_power_interface|`utf8`|
|default_raid_interface|`utf8`|
|default_rescue_interface|`utf8`|
|default_storage_interface|`utf8`|
|default_vendor_interface|`utf8`|
|enabled_bios_interfaces|`list<item: utf8, nullable>`|
|enabled_boot_interfaces|`list<item: utf8, nullable>`|
|enabled_console_interfaces|`list<item: utf8, nullable>`|
|enabled_deploy_interfaces|`list<item: utf8, nullable>`|
|enabled_inspect_interfaces|`list<item: utf8, nullable>`|
|enabled_management_interfaces|`list<item: utf8, nullable>`|
|enabled_network_interfaces|`list<item: utf8, nullable>`|
|enabled_power_interfaces|`list<item: utf8, nullable>`|
|enabled_rescue_interfaces|`list<item: utf8, nullable>`|
|enabled_raid_interfaces|`list<item: utf8, nullable>`|
|enabled_storage_interfaces|`list<item: utf8, nullable>`|
|enabled_vendor_interfaces|`list<item: utf8, nullable>`|
|links|`json`|
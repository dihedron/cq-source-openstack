# Table: openstack_compute_hypervisors

This table shows data for Openstack Compute Hypervisors.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|current_workload|`int64`|
|status|`utf8`|
|state|`utf8`|
|disk_available_least|`int64`|
|host_ip|`utf8`|
|free_ram_mb|`int64`|
|hypervisor_hostname|`utf8`|
|hypervisor_type|`utf8`|
|local_gb_used|`int64`|
|memory_mb|`int64`|
|memory_mb_used|`int64`|
|running_vms|`int64`|
|service|`json`|
|servers|`json`|
|vcpus|`int64`|
|vcpus_used|`int64`|
# Table: openstack_hypervisors

This table shows data for Openstack Hypervisors.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|current_workload|Int|
|status|String|
|state|String|
|disk_available_least|Int|
|host_ip|String|
|free_ram_mb|Int|
|hypervisor_hostname|String|
|hypervisor_type|String|
|local_gb_used|Int|
|memory_mb|Int|
|memory_mb_used|Int|
|running_vms|Int|
|service|JSON|
|servers|JSON|
|vcpus|Int|
|vcpus_used|Int|
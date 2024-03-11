# Table: openstack_compute_serverusage

This table shows data for Openstack Compute Serverusage.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|flavor|`utf8`|
|hours|`float64`|
|instance_id|`utf8`|
|local_gb|`int64`|
|memory_mb|`int64`|
|name|`utf8`|
|state|`utf8`|
|tenant_id|`utf8`|
|uptime|`int64`|
|vcpus|`int64`|
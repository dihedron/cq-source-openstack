# Table: openstack_baremetal_allocations

This table shows data for Openstack Baremetal Allocations.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|uuid|`utf8`|
|candidate_nodes|`list<item: utf8, nullable>`|
|last_error|`utf8`|
|name|`utf8`|
|node_uuid|`utf8`|
|state|`utf8`|
|resource_class|`utf8`|
|traits|`list<item: utf8, nullable>`|
|extra|`json`|
|created_at|`timestamp[us, tz=UTC]`|
|updated_at|`timestamp[us, tz=UTC]`|
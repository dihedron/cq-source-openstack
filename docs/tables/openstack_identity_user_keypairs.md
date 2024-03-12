# Table: openstack_identity_user_keypairs

This table shows data for Openstack Identity User Keypairs.

The primary key for this table is **_cq_id**.

## Relations

This table depends on [openstack_identity_users](openstack_identity_users.md).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|name|`utf8`|
|fingerprint|`utf8`|
|public_key|`utf8`|
|private_key|`utf8`|
|user_id|`utf8`|
|type|`utf8`|
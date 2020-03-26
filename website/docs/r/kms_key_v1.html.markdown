---
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: huaweicloudstack_kms-key-v1"
sidebar_current: "docs-huaweicloudstack-resource-kms-key-v1"
description: |-
  Manages a V1 key resource within KMS.
---

# huaweicloudstack\_kms\_key_v1

Manages a V1 key resource within KMS.

## Example Usage

```hcl
resource "huaweicloudstack_kms_key_v1" "key_1" {
  key_alias       = "key_1"
  pending_days    = "7"
  key_description = "first test key"
  realm           = "cn-north-1"
  is_enabled      = true
}
```

## Argument Reference

The following arguments are supported:

* `key_alias` - (Required) The alias in which to create the key. It is required when
    we create a new key. Changing this updates the alias of key.

* `key_description` - (Optional) The description of the key as viewed in HCS console.
    Changing this updates the description of key.

* `realm` - (Optional) Region where a key resides. Changing this creates a new key.

* `pending_days` - (Optional) Duration in days after which the key is deleted
    after destruction of the resource, must be between 7 and 1096 days. It doesn't
    have default value. It only be used when delete a key.

* `is_enabled` - (Optional) Specifies whether the key is enabled. Defaults to true.
    Changing this updates the state of existing key.


## Attributes Reference

The following attributes are exported:

* `key_alias` - See Argument Reference above.
* `key_description` - See Argument Reference above.
* `realm` - See Argument Reference above.
* `key_id` - The globally unique identifier for the key.
* `default_key_flag` - Identification of a Master Key. The value 1 indicates a Default
    Master Key, and the value 0 indicates a key.
* `scheduled_deletion_date` - Scheduled deletion time (time stamp) of a key.
* `domain_id` - ID of a user domain for the key.
* `expiration_time` - Expiration time.
* `creation_date` - Creation time (time stamp) of a key.
* `is_enabled` - See Argument Reference above.


## Import

KMS Keys can be imported using the `id`, e.g.

```
$ terraform import huaweicloudstack_kms_key_v1.key_1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```

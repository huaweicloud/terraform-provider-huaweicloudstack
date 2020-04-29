## 1.2.0 (Unreleased)

FEATURES:

* **New Resource:** `huaweicloudstack_kms_key_v1` ([#18](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/18))
* **New Data Source:** `huaweicloudstack_kms_key_v1` ([#18](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/18))
* **New Data Source:** `huaweicloudstack_kms_data_key_v1` ([#18](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/18))

ENHANCEMENTS:

* **provider:** Add custom endpoints parameter in the provider ([#19](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/19))
* resource/huaweicloudstack_as_group_v1: support tags in autoscaling group resource ([#21](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/21))
* resource/huaweicloudstack_compute_instance_v2: support tags in compute instance resource ([#22](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/22))
* resource/huaweicloudstack_compute_instance_v2: support import compute instance ([#24](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/24))

BUG FIXES:

* data source/huaweicloudstack_networking*: raising an exception when creating network client failed ([#23](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/23))


## 1.1.0 (December 27, 2019)

FEATURES:

* **New Data Source:** `huaweicloudstack_images_image_v2` ([#5](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/5))
* **New Data Source:** `huaweicloudstack_networking_port_v2` ([#3](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/3))
* **New Data Source:** `huaweicloudstack_networking_subnet_v2` ([#3](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/3))
* **New Resource:** `huaweicloudstack_lb_certificate_v2` ([#11](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/11))
* **New Resource:** `huaweicloudstack_lb_l7policy_v2` ([#11](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/11))
* **New Resource:** `huaweicloudstack_lb_l7rule_v2` ([#11](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/11))
* **New Resource:** `huaweicloudstack_lb_whitelist_v2` ([#11](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/11))
* **New Resource:** `huaweicloudstack_lb_listener_v2` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/10))
* **New Resource:** `huaweicloudstack_lb_loadbalancer_v2` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/10))
* **New Resource:** `huaweicloudstack_lb_member_v2` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/10))
* **New Resource:** `huaweicloudstack_lb_monitor_v2` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/10))
* **New Resource:** `huaweicloudstack_lb_pool_v2` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/10))

ENHANCEMENTS:
* resource/huaweicloudstack_as_group_v1: Add `lbaas_listeners`, `scaling_group_status` and `current_instance_number` attributes ([#12](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/12))
* resource/huaweicloudstack_as_group_v1: Mark `lb_listener_id` as deprecated ([#12](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/12))
* resource/huaweicloudstack_networking_port_v2: Add `no_security_groups` and remove `extra_dhcp_option` attribute ([#3](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/3))
* resource/huaweicloudstack_networking_router_v2: Add `external_network_id` and remove `external_gateway` attribute ([#2](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/2))

BUG FIXES:
* resource/huaweicloudstack_as_configuration_v1: Update validated values of `volume_type` ([#7](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/7))
* Clean up unsupported `availability_zone_hints` parameter ([#2](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/2))
* Clean up unused `value_specs` parameter ([#1](https://github.com/terraform-providers/terraform-provider-huaweicloudstack/pull/1))

## 1.0.0 (November 13, 2019)

FEATURES:

* **New Data Source:** `huaweicloudstack_networking_network_v2`
* **New Data Source:** `huaweicloudstack_networking_secgroup_v2`
* **New Resource:** `huaweicloudstack_as_group_v1`
* **New Resource:** `huaweicloudstack_as_configuration_v1`
* **New Resource:** `huaweicloudstack_as_policy_v1`
* **New Resource:** `huaweicloudstack_blockstorage_volume_v2`
* **New Resource:** `huaweicloudstack_compute_instance_v2`
* **New Resource:** `huaweicloudstack_compute_interface_attach_v2`
* **New Resource:** `huaweicloudstack_compute_keypair_v2`
* **New Resource:** `huaweicloudstack_compute_servergroup_v2`
* **New Resource:** `huaweicloudstack_compute_floatingip_associate_v2`
* **New Resource:** `huaweicloudstack_compute_volume_attach_v2`
* **New Resource:** `huaweicloudstack_networking_network_v2`
* **New Resource:** `huaweicloudstack_networking_subnet_v2`
* **New Resource:** `huaweicloudstack_networking_floatingip_v2`
* **New Resource:** `huaweicloudstack_networking_floatingip_associate_v2`
* **New Resource:** `huaweicloudstack_networking_port_v2`
* **New Resource:** `huaweicloudstack_networking_router_v2`
* **New Resource:** `huaweicloudstack_networking_router_interface_v2`
* **New Resource:** `huaweicloudstack_networking_router_route_v2`
* **New Resource:** `huaweicloudstack_networking_secgroup_v2`
* **New Resource:** `huaweicloudstack_networking_secgroup_rule_v2`
* **New Resource:** `huaweicloudstack_networking_vip_v2`
* **New Resource:** `huaweicloudstack_networking_vip_associate_v2`

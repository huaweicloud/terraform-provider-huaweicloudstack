---
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: huaweicloudstack_lb_pool_v2"
sidebar_current: "docs-huaweicloudstack-resource-lb-pool-v2"
description: |-
  Manages a V2 pool resource within HuaweiCloudStack.
---

# huaweicloudstack\_lb\_pool\_v2

Manages a V2 pool resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "huaweicloudstack_lb_pool_v2" "pool_1" {
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Human-readable name for the pool.

* `description` - (Optional) Human-readable description for the pool.

* `protocol` - (Required) The IP protocol, can either be TCP, HTTP or UDP.
    Changing this creates a new pool.

* `loadbalancer_id` - (Optional) The load balancer on which to provision this
    pool. Changing this creates a new pool.
    Note:  One of LoadbalancerID or ListenerID must be provided.

* `listener_id` - (Optional) The Listener on which the members of the pool
    will be associated with. Changing this creates a new pool.
    Note:  One of LoadbalancerID or ListenerID must be provided.

* `lb_method` - (Required) The load balancing algorithm to
    distribute traffic to the pool's members. Must be one of
    ROUND_ROBIN, LEAST_CONNECTIONS, or SOURCE_IP.

* `persistence` - Omit this field to prevent session persistence.  Indicates
    whether connections in the same session will be processed by the same Pool
    member or not. Changing this creates a new pool.

* `admin_state_up` - (Optional) The administrative state of the pool.
    A valid value is true (UP) or false (DOWN).

The `persistence` argument supports:

* `type` - (Required) The type of persistence mode. The current specification
    supports SOURCE_IP, HTTP_COOKIE, and APP_COOKIE.

* `cookie_name` - (Optional) The name of the cookie if persistence mode is set
    appropriately. It's only supported in the `APP_COOKIE` type.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the pool.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `lb_method` - See Argument Reference above.
* `persistence` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.

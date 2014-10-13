---
layout: "consul"
page_title: "Consul: consul_node"
sidebar_current: "docs-consul-resource-node"
---

# consul\_node

Provides a Consul node resource.

## Example Usage

```
resource "consul_node" "app" {
    name = "app-node"
    address = "10.1.2.3"
    datacenter = "nyc1"
    service {
        name = "app"
        port = 80
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the node.

* `address` - (Required) The address of the node. Can be any format (IP address,
  domain name, etc).

* `datacenter` - (Optional) The datacenter to use. This overrides the
  datacenter in the provider setup and the agent's default datacenter.

* `service` - (Required) Can be specified once for each service. The
  fields supported are documented below.

The `service` block supports the following:

* `name` - (Required) The name of the service.

* `id` - (Optional) The ID of the service. If unspecified, defaults to the
  service name.

* `port` - (Required) The TCP/UDP port of the service.

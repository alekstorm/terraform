package consul

import (
	"github.com/hashicorp/terraform/helper/resource"
)

// resourceMap is the mapping of resources we support to their basic
// operations. This makes it easy to implement new resource types.
var resourceMap *resource.Map

func init() {
	resourceMap = &resource.Map{
		Mapping: map[string]resource.Resource{
			"consul_keys": resource.Resource{
				ConfigValidator: resource_consul_keys_validation(),
				Create:          resource_consul_keys_create,
				Destroy:         resource_consul_keys_destroy,
				Update:          resource_consul_keys_update,
				Diff:            resource_consul_keys_diff,
				Refresh:         resource_consul_keys_refresh,
			},
			"consul_node": resource.Resource{
				ConfigValidator: resource_consul_node_validation(),
				Create:          resource_consul_node_create,
				Destroy:         resource_consul_node_destroy,
				Diff:            resource_consul_node_diff,
				Refresh:         resource_consul_node_refresh,
			},
		},
	}
}

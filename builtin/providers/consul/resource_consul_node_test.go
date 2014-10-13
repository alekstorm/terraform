package consul

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/armon/consul-api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccConsulNode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulNodeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulNodeConfig,
				Check: testAccCheckConsulNodeExists(),
			},
		},
	})
}

func testAccCheckConsulNodeDestroy(s *terraform.State) error {
	c := testAccProvider.client.Catalog()
	qOpts := consulapi.QueryOptions{
		Datacenter: "nyc1",
		AllowStale: false,
		RequireConsistent: true,
	}
	node, _, err := c.Node("app-node", &qOpts)
	if err != nil {
		return err
	}
	if node != nil {
		return fmt.Errorf("Node still exists: %#v", node)
	}
	return nil
}

func testAccCheckConsulNodeExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.client.Catalog()
		qOpts := consulapi.QueryOptions{
			Datacenter: "nyc1",
			AllowStale: false,
			RequireConsistent: true,
		}
		node, _, err := c.Node("app-node", &qOpts)
		if err != nil {
			return err
		}
		if node == nil {
			return fmt.Errorf("Node 'app-node' does not exist")
		}
		if len(node.Services) != 1 {
			return fmt.Errorf("Node does not contain exactly one service")
		}
		expected := consulapi.AgentService{
			ID: "app",
			Service: "app",
			Tags: nil,
			Port: 80,
		}
		service := node.Services["app"]
		if !reflect.DeepEqual(*node.Services["app"], expected) {
			return fmt.Errorf("Service does not match expected: %#v", service)
		}
		return nil
	}
}

const testAccConsulNodeConfig = `
resource "consul_node" "app" {
    name = "app-node"
    address = "10.1.2.3"
    datacenter = "nyc1"
    service {
        name = "app"
        port = 80
    }
}
`

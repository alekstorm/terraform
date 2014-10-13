package consul

import (
	"fmt"
	"log"
	"strconv"

	"github.com/armon/consul-api"
	"github.com/hashicorp/terraform/helper/config"
	"github.com/hashicorp/terraform/helper/diff"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/terraform"
)

func resource_consul_node_validation() *config.Validator {
	return &config.Validator{
		Required: []string{
			"name",
			"address",
			"service.*.name",
			"service.*.port",
		},
		Optional: []string{
			"datacenter",
			"service.*.id",
			//"service.*.tags", // TODO
		},
	}
}

func resource_consul_node_create(
	s *terraform.InstanceState,
	d *terraform.InstanceDiff,
	meta interface{}) (*terraform.InstanceState, error) {
	p := meta.(*ResourceProvider)

	/*f, _ := os.OpenFile("/Users/astorm/foo.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	defer f.Close()

	log.SetOutput(f)*/

	// Merge the diff into the state so that we have all the attributes
	// properly.
	rs := s.MergeDiff(d)
	rs.ID = "consul" // TODO should be nil or ""?

	// Check if the datacenter should be computed
	dc := rs.Attributes["datacenter"]
	if aDiff, ok := d.Attributes["datacenter"]; ok && aDiff.NewComputed {
		var err error
		dc, err = get_dc(p.client)
		if err != nil {
			return rs, fmt.Errorf("Failed to get agent datacenter: %v", err)
		}
		rs.Attributes["datacenter"] = dc
	}

	c := p.client.Catalog()
	if _, ok := rs.Attributes["service.#"]; ok {
		for _, foo := range flatmap.Expand(rs.Attributes, "service").([]interface{}) {
			sv := foo.(map[string]interface{})
			port, err := strconv.ParseInt(sv["port"].(string), 0, 0)
			if err != nil {
				return rs, err
			}

			id := ""
			if attr, ok := sv["id"].(string); ok {
				id = attr
			}

			asv := consulapi.AgentService{
				ID: id,
				Service: sv["name"].(string),
				Port: int(port),
			}
			opts := consulapi.CatalogRegistration{
				Node: rs.Attributes["name"],
				Address: rs.Attributes["address"],
				Service: &asv,
			}
			wOpts := consulapi.WriteOptions{
				Datacenter: rs.Attributes["datacenter"],
			}
			_, err = c.Register(&opts, &wOpts)
			if err != nil {
				return rs, err
			}
		}
	}

	return rs, nil
}

func resource_consul_node_destroy(
	s *terraform.InstanceState,
	meta interface{}) error {
	p := meta.(*ResourceProvider)
	c := p.client.Catalog()

	opts := consulapi.CatalogDeregistration{
		Node: s.Attributes["name"],
		Address: s.Attributes["address"],
	}

	wOpts := consulapi.WriteOptions{
		Datacenter: s.Attributes["datacenter"],
	}

	_, err := c.Deregister(&opts, &wOpts)
	return err
}

func resource_consul_node_diff(
	s *terraform.InstanceState,
	c *terraform.ResourceConfig,
	meta interface{}) (*terraform.InstanceDiff, error) {

	b := &diff.ResourceBuilder{
		Attrs: map[string]diff.AttrType{
			"name":       diff.AttrTypeCreate,
			"address":    diff.AttrTypeCreate,
			"datacenter": diff.AttrTypeCreate,
			"service":    diff.AttrTypeCreate, // TODO AttrTypeUpdate
		},

		ComputedAttrs: []string{
			"address",
		},
	}

	return b.Diff(s, c)
}

func resource_consul_node_refresh(
	s *terraform.InstanceState,
	meta interface{}) (*terraform.InstanceState, error) {
	p := meta.(*ResourceProvider)
	c := p.client.Catalog()

	qOpts := consulapi.QueryOptions{
		Datacenter: s.Attributes["Datacenter"],
		AllowStale: false,
		RequireConsistent: true,
	}
	n, _, err := c.Node(s.Attributes["name"], &qOpts)
	if err != nil {
		return nil, fmt.Errorf("Error refreshing Node: %s", err)
	}
	if n == nil {
		return nil, nil
	}
	s.Attributes["address"] = n.Node.Address

	toFlatten := make(map[string]interface{})
	numServices := len(n.Services)
	services := make([]map[string]interface{}, numServices)
	i := 0
	for _, service := range n.Services {
		log.Println("%i", i)
		n := make(map[string]interface{})
		n["name"] = service.Service
		n["port"] = service.Port
		services[i] = n
		i++
	}

	toFlatten["service"] = services
	for k, v := range flatmap.Flatten(toFlatten) {
		s.Attributes[k] = v
	}

	return s, nil
}

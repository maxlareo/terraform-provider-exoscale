package exoscale

import (
	"fmt"
	"testing"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNetwork(t *testing.T) {
	net := new(egoscale.Network)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("exoscale_network.net", net),
					testAccCheckNetworkAttributes(net),
					testAccCheckNetworkCreateAttributes("terraform-test-network"),
				),
			},
		},
	})
}

func testAccCheckNetworkExists(n string, net *egoscale.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Network ID is set")
		}

		client := GetComputeClient(testAccProvider.Meta())
		net.ID = rs.Primary.ID
		if err := client.Get(net); err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckNetworkAttributes(net *egoscale.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if net.ID == "" {
			return fmt.Errorf("Network is nil")
		}

		return nil
	}
}

func testAccCheckNetworkCreateAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "exoscale_network" {
				continue
			}

			if rs.Primary.Attributes["name"] != name {
				continue
			}

			if rs.Primary.Attributes["display_text"] == "" {
				return fmt.Errorf("Network: expected display text to be set")
			}

			return nil
		}

		return fmt.Errorf("Could not find Network name: %s", name)
	}
}

func testAccCheckNetworkDestroy(s *terraform.State) error {
	client := GetComputeClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "exoscale_network" {
			continue
		}

		key := &egoscale.Network{ID: rs.Primary.ID}
		if err := client.Get(key); err != nil {
			if r, ok := err.(*egoscale.ErrorResponse); ok {
				if r.ErrorCode == egoscale.ParamError {
					return nil
				}
			}
			return err
		}
	}
	return fmt.Errorf("Network: still exists")
}

var testAccNetworkCreate = `
resource "exoscale_network" "net" {
  name = "terraform-test-network"
  display_text = "Terraform Acceptance Test"
  zone = "ch-dk-2"
  network_offering = "PrivNet"

  tags {
    managedby = "terraform"
  }
}
`

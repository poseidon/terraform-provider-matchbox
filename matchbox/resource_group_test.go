package matchbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/poseidon/matchbox/matchbox/storage/testfakes"
)

func TestResourceGroup(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go srv.Start()
	defer srv.Stop()

	hcl := `
		resource "matchbox_group" "default" {
 			name    = "default"
  			profile = "foo"
  			selector = {
				  qux = "baz"
			}

			metadata = {
				foo = "bar"
 			}
		}
	`

	check := func(s *terraform.State) error {
		grp, err := srv.Store.GroupGet("default")
		if err != nil {
			return err
		}

		if grp.GetId() != "default" {
			return fmt.Errorf("id, found %q", grp.GetId())
		}

		if grp.GetProfile() != "foo" {
			return fmt.Errorf("profile, found %q", grp.GetProfile())
		}

		selector := grp.GetSelector()
		if len(selector) != 1 || selector["qux"] != "baz" {
			return fmt.Errorf("selector.qux, found %q", selector["qux"])
		}

		if string(grp.GetMetadata()) != "{\"foo\":\"bar\"}" {
			return fmt.Errorf("metadata, found %q", grp.GetProfile())
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  providers,
		Steps: []resource.TestStep{{
			Config: srv.AddProviderConfig(hcl),
			Check:  check,
		}},
	})

}

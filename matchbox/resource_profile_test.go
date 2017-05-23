package matchbox

import (
	"fmt"
	"testing"

	"github.com/coreos/matchbox/matchbox/storage/testfakes"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestResourceProfile(t *testing.T) {
	srv := NewFixtureServer(testfakes.NewFixedStore())
	srv.Start()
	defer srv.Stop()

	hcl := `
		resource "matchbox_profile" "default" {
			name   = "default"
			kernel = "foo"

			initrd = [
				"bar",
			]

			args = [
				"qux",
			]

			container_linux_config = "baz"
		}
	`

	check := func(s *terraform.State) error {
		profile, err := srv.Store.ProfileGet("default")
		if err != nil {
			return err
		}

		if profile.GetId() != "default" {
			return fmt.Errorf("id, found %q", profile.GetId())
		}

		if profile.GetIgnitionId() != "default.yaml.tmpl" {
			return fmt.Errorf("profile, found %q", profile.GetIgnitionId())
		}

		boot := profile.GetBoot()
		if boot.GetKernel() != "foo" {
			return fmt.Errorf("kernel, found %d", boot.GetKernel())
		}

		initrd := boot.GetInitrd()
		if len(initrd) != 1 || initrd[0] != "bar" {
			return fmt.Errorf("kernel, found %v", initrd)
		}

		args := boot.GetArgs()
		if len(args) != 1 || args[0] != "qux" {
			return fmt.Errorf("args, found %v", args)
		}

		ignition, err := srv.Store.IgnitionGet("default.yaml.tmpl")
		if err != nil {
			return err
		}

		if ignition != "baz" {
			return fmt.Errorf("container_linux_config, found %q", ignition)
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

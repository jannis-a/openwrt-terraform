//go:build acceptance.test

package bridgevlan_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	client := openWrtServer.LuCIRPCClient(
		ctx,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()
	options := lucirpc.Options{
		"device": lucirpc.String("br-testing"),
		"ports":  lucirpc.ListString([]string{"eth0:t", "eth1"}),
		"vlan":   lucirpc.Integer(4),
	}
	ok, err := client.CreateSection(ctx, "network", "bridge-vlan", "br_testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_network_bridge_vlan" "this" {
	id = "br_testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "id", "br_testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "device", "br-testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "ports.0", "eth0:t"),
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "ports.1", "eth1"),
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "vlan", "4"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestResourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_bridge_vlan" "br_testing" {
	id = "br_testing"
	device = "br-testing"
	ports = [
		"eth0:t",
		"eth1",
	]
	vlan = 4
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "id", "br_testing"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.0", "eth0:t"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.1", "eth1"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "vlan", "4"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_bridgevlan.br_testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_bridge_vlan" "br_testing" {
	id = "br_testing"
	device = "br-testing"
	ports = [
		"eth0",
		"eth1:t",
	]
	vlan = 6
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "id", "br_testing"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.0", "eth0"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.1", "eth1:t"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "vlan", "6"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}

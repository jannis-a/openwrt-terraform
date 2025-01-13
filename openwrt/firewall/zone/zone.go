package zone

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	nameAttribute            = "name"
	nameAttributeDescription = "The name of the zone."
	nameUCIOption            = "name"

	forwardAttribute            = "forward"
	forwardAttributeDescription = "Zone forwarding policy."
	forwardUCIOption            = "forward"

	inputAttribute            = "input"
	inputAttributeDescription = "Zone input policy."
	inputUCIOption            = "input"

	outputAttribute            = "output"
	outputAttributeDescription = "Zone output policy."
	outputUCIOption            = "output"

	networkAttribute            = "network"
	networkAttributeDescription = "List of network interfaces this zone applies to."
	networkUCIOption            = "network"

	masqAttribute            = "masquerade"
	masqAttributeDescription = "Enable masquerading on this zone. Needed for NAT."
	masqUCIOption            = "masq"

	mtuFixAttribute            = "mssclamp"
	mtuFixAttributeDescription = "Enable MSS clamping for zones that have none default MTU."
	mtuFixUCIOption            = "mtu_fix"

	schemaDescription = "Firewall zone configurations to associate with network interfaces."

	uciConfig = "firewall"
	uciType   = "zone"

	typeAccept = "ACCEPT"
	typeReject = "REJECT"
	typeDrop   = "DROP"
)

var (
	TypeValidators = []validator.String{
		stringvalidator.OneOf(
			typeAccept,
			typeReject,
			typeDrop,
		),
	}

	nameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       nameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetName, nameAttribute, nameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetName, nameAttribute, nameUCIOption),
	}

	forwardSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       forwardAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetForward, forwardAttribute, forwardUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetForward, forwardAttribute, forwardUCIOption),
		Validators:        TypeValidators,
	}

	inputSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       inputAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetInput, inputAttribute, inputUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetInput, inputAttribute, inputUCIOption),
		Validators:        TypeValidators,
	}

	outputSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       outputAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetOutput, outputAttribute, outputUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetOutput, outputAttribute, outputUCIOption),
		Validators:        TypeValidators,
	}

	networkSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       networkAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetNetwork, networkAttribute, networkUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetNetwork, networkAttribute, networkUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	}

	masqSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       masqAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetMasq, masqAttribute, masqUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetMasq, masqAttribute, masqUCIOption),
	}

	mtuFixSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mtuFixAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetMtuFix, mtuFixAttribute, mtuFixUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetMtuFix, mtuFixAttribute, mtuFixUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		forwardAttribute:        forwardSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		inputAttribute:          inputSchemaAttribute,
		outputAttribute:         outputSchemaAttribute,
		nameAttribute:           nameSchemaAttribute,
		networkAttribute:        networkSchemaAttribute,
		masqAttribute:           masqSchemaAttribute,
		mtuFixAttribute:         mtuFixSchemaAttribute,
	}
)

func NewDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

func NewResource() resource.Resource {
	return lucirpcglue.NewResource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

type model struct {
	Id         types.String `tfsdk:"id"`
	Forward    types.String `tfsdk:"forward"`
	Output     types.String `tfsdk:"output"`
	Input      types.String `tfsdk:"input"`
	Name       types.String `tfsdk:"name"`
	Network    types.List   `tfsdk:"network"`
	Masquerade types.Bool   `tfsdk:"masquerade"`
	MssClamp   types.Bool   `tfsdk:"mssclamp"`
}

func modelGetOutput(m model) types.String  { return m.Output }
func modelGetId(m model) types.String      { return m.Id }
func modelGetInput(m model) types.String   { return m.Input }
func modelGetName(m model) types.String    { return m.Name }
func modelGetForward(m model) types.String { return m.Forward }
func modelGetNetwork(m model) types.List   { return m.Network }
func modelGetMasq(m model) types.Bool      { return m.Masquerade }
func modelGetMtuFix(m model) types.Bool    { return m.MssClamp }

func modelSetOutput(m *model, value types.String)  { m.Output = value }
func modelSetForward(m *model, value types.String) { m.Forward = value }
func modelSetInput(m *model, value types.String)   { m.Input = value }
func modelSetName(m *model, value types.String)    { m.Name = value }
func modelSetId(m *model, value types.String)      { m.Id = value }
func modelSetNetwork(m *model, value types.List)   { m.Network = value }
func modelSetMasq(m *model, value types.Bool)      { m.Masquerade = value }
func modelSetMtuFix(m *model, value types.Bool)    { m.MssClamp = value }

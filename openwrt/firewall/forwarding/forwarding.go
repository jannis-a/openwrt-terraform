package forwarding

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	destAttribute            = "dest"
	destAttributeDescription = "zone dest"
	destUCIOption            = "dest"

	srcAttribute            = "src"
	srcAttributeDescription = "zone src"
	srcUCIOption            = "src"

	familyAttribute            = "family"
	familyAttributeDescription = `Applies the rule to specific protocol families. Defaults to automatically determining which family if unset.`
	familyUCIOption            = "family"
	familyIpv4                 = "ipv4"
	familyIpv6                 = "ipv6"
	familyAny                  = "any"

	schemaDescription = "Firewall zone rules for controlling flow between zones."

	uciConfig = "firewall"
	uciType   = "forwarding"
)

var (
	destSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDest, destAttribute, destUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDest, destAttribute, destUCIOption),
	}

	srcSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrc, srcAttribute, srcUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrc, srcAttribute, srcUCIOption),
	}

	familySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       familyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetFamily, familyAttribute, familyUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetFamily, familyAttribute, familyUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				familyIpv6,
				familyIpv4,
				familyAny,
			),
		},
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		srcAttribute:            srcSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		destAttribute:           destSchemaAttribute,
		familyAttribute:         familySchemaAttribute,
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
	Id     types.String `tfsdk:"id"`
	Src    types.String `tfsdk:"src"`
	Dest   types.String `tfsdk:"dest"`
	Family types.String `tfsdk:"family"`
}

func modelGetSrc(m model) types.String    { return m.Src }
func modelGetId(m model) types.String     { return m.Id }
func modelGetDest(m model) types.String   { return m.Dest }
func modelGetFamily(m model) types.String { return m.Family }

func modelSetSrc(m *model, value types.String)    { m.Src = value }
func modelSetDest(m *model, value types.String)   { m.Dest = value }
func modelSetId(m *model, value types.String)     { m.Id = value }
func modelSetFamily(m *model, value types.String) { m.Family = value }

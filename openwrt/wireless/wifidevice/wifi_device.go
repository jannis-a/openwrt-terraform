package wifidevice

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	bandAttribute            = "band"
	bandAttributeDescription = `Channel width. Must be one of: "2g", "5g", "6g".`
	band2G                   = "2g"
	band5G                   = "5g"
	band6G                   = "6g"
	bandUCIOption            = "band"

	cellDensityAttribute            = "cell_density"
	cellDensityAttributeDescription = "Configures data rates based on the coverage cell density. Must be one of 0, 1, 2, 3."
	cellDensityDisabled             = 0
	cellDensityHigh                 = 2
	cellDensityNormal               = 1
	cellDensityUCIOption            = "cell_density"
	cellDensityVeryHigh             = 3

	channelAttribute            = "channel"
	channelAttributeDescription = `The wireless channel. Currently, only "auto" is supported.`
	channelAuto                 = "auto"
	channel1                    = "1"
	channel2                    = "2"
	channel3                    = "3"
	channel4                    = "4"
	channel5                    = "5"
	channel6                    = "6"
	channel7                    = "7"
	channel8                    = "8"
	channel9                    = "9"
	channel10                   = "10"
	channel11                   = "11"
	channel12                   = "12"
	channel13                   = "13"
	channel14                   = "14"
	channel32                   = "32"
	channel36                   = "36"
	channel38                   = "38"
	channel40                   = "40"
	channel42                   = "42"
	channel44                   = "44"
	channel46                   = "46"
	channel48                   = "48"
	channel50                   = "50"
	channel52                   = "52"
	channel54                   = "54"
	channel56                   = "56"
	channel58                   = "58"
	channel60                   = "60"
	channel62                   = "62"
	channel64                   = "64"
	channel68                   = "68"
	channel96                   = "96"
	channel100                  = "100"
	channel102                  = "102"
	channel104                  = "104"
	channel106                  = "106"
	channel108                  = "108"
	channel110                  = "110"
	channel112                  = "112"
	channel114                  = "114"
	channel116                  = "116"
	channel118                  = "118"
	channel120                  = "120"
	channel122                  = "122"
	channel124                  = "124"
	channel126                  = "126"
	channel128                  = "128"
	channel132                  = "132"
	channel134                  = "134"
	channel136                  = "136"
	channel138                  = "138"
	channel140                  = "140"
	channel142                  = "142"
	channel144                  = "144"
	channel149                  = "149"
	channel151                  = "151"
	channel153                  = "153"
	channel155                  = "155"
	channel157                  = "157"
	channel159                  = "159"
	channel161                  = "161"
	channel163                  = "163"
	channel165                  = "165"
	channel167                  = "167"
	channel169                  = "169"
	channel171                  = "171"
	channel173                  = "173"
	channel175                  = "175"
	channel177                  = "177"
	channelUCIOption            = "channel"

	countryCodeAttribute            = "country"
	countryCodeAttributeDescription = `Two-digit country code. E.g. "US".`
	countryCodeUCIOption            = "country"

	htModeAttribute            = "htmode"
	htModeAttributeDescription = `Channel width. Must be one of: "HE20", "HE40", "HE80", "HE160", "HT20", "HT40", "HT40-", "HT40+", "NONE", "VHT20", "VHT40", "VHT80", "VHT160".`
	htModeHE160                = "HE160"
	htModeHE20                 = "HE20"
	htModeHE40                 = "HE40"
	htModeHE80                 = "HE80"
	htModeHT20                 = "HT20"
	htModeHT40                 = "HT40"
	htModeHT40Minus            = "HT40-"
	htModeHT40Plus             = "HT40+"
	htModeNone                 = "NONE"
	htModeUCIOption            = "htmode"
	htModeVHT160               = "VHT160"
	htModeVHT20                = "VHT20"
	htModeVHT40                = "VHT40"
	htModeVHT80                = "VHT80"

	pathAttribute            = "path"
	pathAttributeDescription = "Path of the device in `/sys/devices`."
	pathUCIOption            = "path"

	schemaDescription = "The physical radio device."

	typeAttribute            = "type"
	typeAttributeDescription = `The type of device. Currently only "mac80211" is supported.`
	typeMac80211             = "mac80211"
	typeUCIOption            = "type"

	uciConfig = "wireless"
	uciType   = "wifi-device"
)

var (
	bandSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       bandAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetBand, bandAttribute, bandUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetBand, bandAttribute, bandUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				band2G,
				band5G,
				band6G,
			),
		},
	}

	cellDensitySchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       cellDensityAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetCellDensity, cellDensityAttribute, cellDensityUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetCellDensity, cellDensityAttribute, cellDensityUCIOption),
		Validators: []validator.Int64{
			int64validator.OneOf(
				cellDensityDisabled,
				cellDensityNormal,
				cellDensityHigh,
				cellDensityVeryHigh,
			),
		},
	}

	channelSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       channelAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetChannel, channelAttribute, channelUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetChannel, channelAttribute, channelUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				channelAuto,
				channel1,
				channel2,
				channel3,
				channel4,
				channel5,
				channel6,
				channel7,
				channel8,
				channel9,
				channel10,
				channel11,
				channel12,
				channel13,
				channel14,
				channel32,
				channel36,
				channel38,
				channel40,
				channel42,
				channel44,
				channel46,
				channel48,
				channel50,
				channel52,
				channel54,
				channel56,
				channel58,
				channel60,
				channel62,
				channel64,
				channel68,
				channel96,
				channel100,
				channel102,
				channel104,
				channel106,
				channel108,
				channel110,
				channel112,
				channel114,
				channel116,
				channel118,
				channel120,
				channel122,
				channel124,
				channel126,
				channel128,
				channel132,
				channel134,
				channel136,
				channel138,
				channel140,
				channel142,
				channel144,
				channel149,
				channel151,
				channel153,
				channel155,
				channel157,
				channel159,
				channel161,
				channel163,
				channel165,
				channel167,
				channel169,
				channel171,
				channel173,
				channel175,
				channel177,
			),
		},
	}

	countryCodeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       countryCodeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetCountryCode, countryCodeAttribute, countryCodeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetCountryCode, countryCodeAttribute, countryCodeUCIOption),
		Validators: []validator.String{
			stringvalidator.LengthBetween(2, 2),
		},
	}

	htModeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       htModeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetHTMode, htModeAttribute, htModeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetHTMode, htModeAttribute, htModeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				htModeHE160,
				htModeHE20,
				htModeHE40,
				htModeHE80,
				htModeHT20,
				htModeHT40,
				htModeHT40Minus,
				htModeHT40Plus,
				htModeNone,
				htModeVHT160,
				htModeVHT20,
				htModeVHT40,
				htModeVHT80,
			),
		},
	}

	pathSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       pathAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetPath, pathAttribute, pathUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetPath, pathAttribute, pathUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		bandAttribute:           bandSchemaAttribute,
		cellDensityAttribute:    cellDensitySchemaAttribute,
		channelAttribute:        channelSchemaAttribute,
		countryCodeAttribute:    countryCodeSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		htModeAttribute:         htModeSchemaAttribute,
		pathAttribute:           pathSchemaAttribute,
		typeAttribute:           typeSchemaAttribute,
	}

	typeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       typeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetType, typeAttribute, typeUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetType, typeAttribute, typeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				typeMac80211,
			),
		},
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
	Band        types.String `tfsdk:"band"`
	CellDensity types.Int64  `tfsdk:"cell_density"`
	Channel     types.String `tfsdk:"channel"`
	CountryCode types.String `tfsdk:"country"`
	HTMode      types.String `tfsdk:"htmode"`
	Id          types.String `tfsdk:"id"`
	Path        types.String `tfsdk:"path"`
	Type        types.String `tfsdk:"type"`
}

func modelGetBand(m model) types.String        { return m.Band }
func modelGetCellDensity(m model) types.Int64  { return m.CellDensity }
func modelGetChannel(m model) types.String     { return m.Channel }
func modelGetCountryCode(m model) types.String { return m.CountryCode }
func modelGetHTMode(m model) types.String      { return m.HTMode }
func modelGetId(m model) types.String          { return m.Id }
func modelGetPath(m model) types.String        { return m.Path }
func modelGetType(m model) types.String        { return m.Type }

func modelSetBand(m *model, value types.String)        { m.Band = value }
func modelSetCellDensity(m *model, value types.Int64)  { m.CellDensity = value }
func modelSetChannel(m *model, value types.String)     { m.Channel = value }
func modelSetCountryCode(m *model, value types.String) { m.CountryCode = value }
func modelSetHTMode(m *model, value types.String)      { m.HTMode = value }
func modelSetId(m *model, value types.String)          { m.Id = value }
func modelSetPath(m *model, value types.String)        { m.Path = value }
func modelSetType(m *model, value types.String)        { m.Type = value }

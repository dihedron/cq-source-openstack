package compute

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
)

func InstanceFlavors() *schema.Table {
	return &schema.Table{
		Name:     "openstack_instance_flavors",
		Resolver: fetchInstanceFlavors,
		Transform: transformers.TransformWithStruct(
			&InstanceFlavor{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Name", "ExtraSpecsObj", "ExtraSpecsMap", "ExtraSpecsRaw"),
		),
		Columns: []schema.Column{
			{
				Name:        "name",
				Type:        arrow.BinaryTypes.String,
				Description: "The original name of the flavor used to start the instance.",
				Resolver:    schema.PathResolver("Name"),
			},
			{
				Name:        "vcpus",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of virtual CPUs in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("VCPUs"),
			},
			{
				Name:        "vgpus",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of virtual GPUs in the flavor used to start the instance.",
				Resolver: transform.Apply(
					transform.OnObjectField("ExtraSpecsObj.VGPUs"),
					transform.ToInt(),
					transform.OrDefault(0),
				),
			},
			{
				Name:        "cores",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of virtual CPU cores in the flavor used to start the instance.",
				Resolver: transform.Apply(
					transform.OnObjectField("ExtraSpecsObj.CPUCores"),
					transform.ToInt(),
					transform.OrDefault(0),
				),
			},
			{
				Name:        "sockets",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of CPU sockets in the flavor used to start the instance.",
				Resolver: transform.Apply(
					transform.OnObjectField("ExtraSpecsObj.CPUSockets"),
					transform.ToInt(),
					transform.OrDefault(0),
				),
			},
			{
				Name:        "ram",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The amount of RAM in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("RAM"),
			},
			{
				Name:        "disk",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The size of the disk in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("Disk"),
			},
			{
				Name:        "swap",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The size of the swap disk in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("Swap"),
			},
			{
				Name:        "ephemeral",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The size of the ephemeral disk in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("Ephemeral"),
			},
			{
				Name:        "rng_allowed",
				Type:        arrow.FixedWidthTypes.Boolean,
				Description: "Whether the RNG is allowed on the flavor used to start the instance.",
				Resolver:    schema.PathResolver("ExtraSpecsObj.RNGAllowed"),
			},
			{
				Name:        "watchdog_action",
				Type:        arrow.BinaryTypes.String,
				Description: "The action to take when the Nova watchdog detects the instance is not responding.",
				Resolver:    schema.PathResolver("ExtraSpecsObj.WatchdogAction"),
			},
		},
	}
}

func fetchInstanceFlavors(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	instance := parent.Item.(*Instance)

	instance.Flavor.lock.Lock()
	defer instance.Flavor.lock.Unlock()

	instance.Flavor.ExtraSpecsMap = &map[string]string{}
	if err := json.Unmarshal(instance.Flavor.ExtraSpecsRaw, &instance.Flavor.ExtraSpecsMap); err != nil {
		api.Logger().Error().Err(err).Str("instance id", instance.ID).Msg("error parsing extra specs as map")
	}

	instance.Flavor.ExtraSpecsObj = &FlavorExtraSpecsData{}
	if err := json.Unmarshal(instance.Flavor.ExtraSpecsRaw, &instance.Flavor.ExtraSpecsObj); err != nil {
		api.Logger().Error().Err(err).Str("instance id", instance.ID).Msg("error parsing extra specs as object")
	}

	api.Logger().Debug().Str("instance id", instance.ID).Str("json", format.ToJSON(&instance.Flavor)).Msg("streaming instance flavor")
	res <- instance.Flavor

	return nil
}

type InstanceFlavor struct {
	lock          sync.RWMutex
	Name          string                `json:"original_name"`
	Disk          int                   `json:"disk"`
	RAM           int                   `json:"ram"`
	Swap          int                   `json:"-"`
	VCPUs         int                   `json:"vcpus"`
	Ephemeral     int                   `json:"OS-FLV-EXT-DATA:ephemeral" cq-name:"ephemeral"`
	ExtraSpecsRaw json.RawMessage       `json:"extra_specs"`
	ExtraSpecsObj *FlavorExtraSpecsData `json:"-" cq-name:"extra_specs"`
	ExtraSpecsMap *map[string]string    `json:"-"`
}

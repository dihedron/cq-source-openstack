package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
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
			transformers.WithSkipFields("Name", "ExtraSpecs"),
		),
		Columns: []schema.Column{
			{
				Name:        "name",
				Type:        schema.TypeString,
				Description: "The original name of the flavor used to start the instance.",
				Resolver:    schema.PathResolver("Name"),
			},
			{
				Name:        "vcpus",
				Type:        schema.TypeInt,
				Description: "The number of virtual CPUs in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("VCPUs"),
			},
			{
				Name:        "vgpus",
				Type:        schema.TypeInt,
				Description: "The number of virtual GPUs in the flavor used to start the instance.",
				Resolver: transform.Apply(
					transform.OnObjectField("ExtraSpecs.VGPUs"),
					transform.ToInt(),
					transform.OrDefault(0),
				),
			},
			{
				Name:        "cores",
				Type:        schema.TypeInt,
				Description: "The number of virtual CPU cores in the flavor used to start the instance.",
				Resolver: transform.Apply(
					transform.OnObjectField("ExtraSpecs.CPUCores"),
					transform.ToInt(),
					transform.OrDefault(0),
				),
			},
			{
				Name:        "sockets",
				Type:        schema.TypeInt,
				Description: "The number of CPU sockets in the flavor used to start the instance.",
				Resolver: transform.Apply(
					transform.OnObjectField("ExtraSpecs.CPUSockets"),
					transform.ToInt(),
					transform.OrDefault(0),
				),
			},
			{
				Name:        "ram",
				Type:        schema.TypeInt,
				Description: "The amount of RAM in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("RAM"),
			},
			{
				Name:        "disk",
				Type:        schema.TypeInt,
				Description: "The size of the disk in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("Disk"),
			},
			{
				Name:        "swap",
				Type:        schema.TypeInt,
				Description: "The size of the swap disk in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("Swap"),
			},
			{
				Name:        "ephemeral",
				Type:        schema.TypeInt,
				Description: "The size of the ephemeral disk in the flavor used to start the instance.",
				Resolver:    schema.PathResolver("Ephemeral"),
			},
			{
				Name:        "rng_allowed",
				Type:        schema.TypeBool,
				Description: "Whether the RNG is allowed on the flavor used to start the instance.",
				Resolver:    schema.PathResolver("ExtraSpecs.RNGAllowed"),
			},
			{
				Name:        "watchdog_action",
				Type:        schema.TypeString,
				Description: "The action to take when the Nova watchdog detects the instance is not responding.",
				Resolver:    schema.PathResolver("ExtraSpecs.WatchdogAction"),
			},
		},
	}
}

func fetchInstanceFlavors(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	instance := parent.Item.(*Instance)
	api.Logger.Debug().Str("instance id", instance.ID).Str("json", format.ToJSON(instance.Flavor)).Msg("streaming instance flavor")
	res <- instance.Flavor

	return nil
}

type InstanceFlavor struct {
	Name       string            `json:"original_name"`
	Disk       int               `json:"disk"`
	RAM        int               `json:"ram"`
	Swap       int               `json:"-"`
	VCPUs      int               `json:"vcpus"`
	Ephemeral  int               `json:"OS-FLV-EXT-DATA:ephemeral" cq-name:"ephemeral"`
	ExtraSpecs *FlavorExtraSpecs `json:"extra_specs"`
}

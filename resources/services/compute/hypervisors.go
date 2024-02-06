package compute

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
)

func Hypervisors() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_hypervisors",
		Resolver: fetchHypervisors,
		Transform: transformers.TransformWithStruct(
			&hypervisors.Hypervisor{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
		Relations: []*schema.Table{
			// ImageMetadata(),
			// ImageProperties(),
			// ImageTags(),
		},
	}
}

func fetchHypervisors(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	compute, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := hypervisors.ListOpts{}

	allPages, err := hypervisors.List(compute, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing hypervisors with options")
		return err
	}
	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting hypervisors")
		return err
	}
	api.Logger().Debug().Int("count", len(allHypervisors)).Msg("hypervisors retrieved")

	for _, hypervisor := range allHypervisors {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		hypervisor := hypervisor
		// api.Logger().Debug().Str("data", format.ToPrettyJSON(hypervisor)).Msg("streaming hypervisor")
		api.Logger().Debug().Str("id", hypervisor.ID).Msg("streaming hypervisor")
		res <- hypervisor
	}
	return nil
}

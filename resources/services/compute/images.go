package compute

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
)

func Images() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_images",
		Resolver: fetchImages,
		Transform: transformers.TransformWithStruct(
			&images.Image{},
		),
	}
}

func fetchImages(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	nova, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := images.ListOpts{}

	allPages, err := images.ListDetail(nova, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing images with options")
		return err
	}
	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting images")
		return err
	}
	api.Logger().Debug().Int("count", len(allImages)).Msg("images retrieved")

	for _, image := range allImages {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		image := image
		api.Logger().Debug().Str("id", image.ID).Msg("streaming image")
		res <- image
	}
	return nil
}

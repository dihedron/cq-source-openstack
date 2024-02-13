package glance

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func Images() *schema.Table {
	return &schema.Table{
		Name:     "openstack_glance_images",
		Resolver: fetchImages,
		Transform: transformers.TransformWithStruct(
			&images.Image{},
			//transformers.WithPrimaryKeys("ID"),
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			//transformers.WithSkipFields("Metadata", "Properties", "Tags"),
			transformers.WithSkipFields("Links"),
		),
		Relations: []*schema.Table{
			// ImageMembers(),
			ImageMetadata(),
			ImageProperties(),
			ImageTags(),
		},
	}
}

func fetchImages(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	glance, err := api.GetServiceClient(client.GlanceV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := images.ListOpts{}

	allPages, err := images.List(glance, opts).AllPages()
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
		// api.Logger().Debug().Str("data", format.ToPrettyJSON(image)).Msg("streaming image")
		api.Logger().Debug().Str("id", image.ID).Msg("streaming image")
		res <- image
	}
	return nil
}

package glance

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/dihedron/cq-source-openstack/resources/internal/utils"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func ImageMetadata() *schema.Table {
	return &schema.Table{
		Name:     "openstack_glance_image_metadata",
		Resolver: fetchImageMetadata,
		Transform: transformers.TransformWithStruct(
			&utils.Pair[string, string]{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			//transformers.WithSkipFields("OriginalName", "ExtraSpecs"),
		),
	}
}

func fetchImageMetadata(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	image := parent.Item.(images.Image)

	for k, v := range image.Metadata {
		pair := &utils.Pair[string, string]{
			Key:   k,
			Value: v,
		}
		api.Logger().Debug().Str("image id", image.ID).Msg("streaming image metadata")
		res <- pair
	}

	return nil
}

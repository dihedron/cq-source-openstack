package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func ImageTags() *schema.Table {
	return &schema.Table{
		Name:     "openstack_image_tags",
		Resolver: fetchImageTags,
		Transform: transformers.TransformWithStruct(
			&Tag{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			//transformers.WithSkipFields("OriginalName", "ExtraSpecs"),
		),
	}
}

func fetchImageTags(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	image := parent.Item.(images.Image)

	if image.Tags != nil {
		for _, v := range image.Tags {
			tag := &Tag{Value: v}
			api.Logger.Debug().Str("image id", image.ID).Msg("streaming image tag")
			res <- tag
		}
	}
	return nil
}

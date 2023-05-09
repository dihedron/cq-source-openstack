package resources

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func ImageProperties() *schema.Table {
	return &schema.Table{
		Name:     "openstack_image_properties",
		Resolver: fetchImageProperties,
		Transform: transformers.TransformWithStruct(
			&Pair[string, string]{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			//transformers.WithSkipFields("OriginalName", "ExtraSpecs"),
		),
	}
}

func fetchImageProperties(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	image := parent.Item.(images.Image)

	for k, v := range image.Properties {
		pair := &Pair[string, string]{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		}
		api.Logger.Debug().Str("image id", image.ID).Msg("streaming image property")
		res <- pair
	}

	return nil
}

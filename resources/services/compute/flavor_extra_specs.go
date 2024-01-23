package compute

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/dihedron/cq-source-openstack/resources/internal/utils"
)

func FlavorExtraSpecs() *schema.Table {
	return &schema.Table{
		Name:     "openstack_flavor_extra_specs",
		Resolver: fetchFlavorExtraSpecs,
		Transform: transformers.TransformWithStruct(
			&utils.Pair[string, string]{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
		),
	}
}

func fetchFlavorExtraSpecs(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	flavor := parent.Item.(*Flavor)

	if flavor.ExtraSpecsMap != nil {
		for k, v := range *flavor.ExtraSpecsMap {
			pair := &utils.Pair[string, string]{
				Key:   k,
				Value: v,
			}
			api.Logger().Debug().Str("flavor id", *flavor.ID).Msg("streaming flavor extra spec")
			res <- pair
		}
	}

	return nil
}

package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
)

func InstanceFlavorExtraSpecs() *schema.Table {
	return &schema.Table{
		Name:     "openstack_instance_flavor_extra_specs",
		Resolver: fetchInstanceFlavorExtraSpecs,
		Transform: transformers.TransformWithStruct(
			&Pair[string, string]{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
		),
	}
}

func fetchInstanceFlavorExtraSpecs(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	instance := parent.Item.(*Instance)

	instance.Flavor.lock.RLock()
	defer instance.Flavor.lock.RUnlock()

	if instance.Flavor.ExtraSpecsMap != nil {
		for k, v := range *instance.Flavor.ExtraSpecsMap {
			pair := &Pair[string, string]{
				Key:   k,
				Value: v,
			}
			api.Logger().Debug().Str("instance id", instance.ID).Msg("streaming instance flavor extra spec")
			res <- pair
		}
	}

	return nil
}

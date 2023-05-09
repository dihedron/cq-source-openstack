package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
)

func InstanceMetadata() *schema.Table {
	return &schema.Table{
		Name:     "openstack_instance_metadata",
		Resolver: fetchInstanceMetadata,
		Transform: transformers.TransformWithStruct(
			&Pair{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			//transformers.WithSkipFields("OriginalName", "ExtraSpecs"),
		),
	}
}

func fetchInstanceMetadata(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	instance := parent.Item.(*Instance)

	for k, v := range instance.Metadata {
		pair := &Pair{
			Key:   k,
			Value: v,
		}
		api.Logger.Debug().Str("instance id", instance.ID).Msg("streaming instance metadata")
		res <- pair
	}

	return nil
}

type Pair struct {
	Key   string `cq-name:"key"`
	Value string `cq-name:"value"`
}

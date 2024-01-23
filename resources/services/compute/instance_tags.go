package compute

import (
	"context"

	"github.com/dihedron/cq-source-openstack/resources/internal/utils"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
)

func InstanceTags() *schema.Table {
	return &schema.Table{
		Name:     "openstack_instance_tags",
		Resolver: fetchInstanceTags,
		Transform: transformers.TransformWithStruct(
			&utils.Tag{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
		),
	}
}

func fetchInstanceTags(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	instance := parent.Item.(*Instance)

	if instance.Tags != nil {
		for _, v := range *instance.Tags {
			tag := &utils.Tag{Value: v}
			api.Logger().Debug().Str("instance id", instance.ID).Msg("streaming instance tag")
			res <- tag
		}
	}
	return nil
}

package nova

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func InstanceAttachedVolumes() *schema.Table {
	return &schema.Table{
		Name:     "openstack_nova_instance_attached_volumes",
		Resolver: fetchInstanceAttachedVolumes,
		Transform: transformers.TransformWithStruct(
			&servers.AttachedVolume{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			//transformers.WithSkipFields("OriginalName", "ExtraSpecs"),
		),
	}
}

func fetchInstanceAttachedVolumes(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	instance := parent.Item.(*Instance)

	for _, av := range instance.AttachedVolumes {
		api.Logger().Debug().Str("instance id", instance.ID).Msg("streaming instance attached volume")
		res <- av
	}
	return nil
}

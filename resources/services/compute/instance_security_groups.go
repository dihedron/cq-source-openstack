package compute

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
)

func InstanceSecurityGroups() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_instance_security_groups",
		Resolver: fetchInstanceSecurityGroups,
		Transform: transformers.TransformWithStruct(
			&SecurityGroup{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			//transformers.WithSkipFields("OriginalName", "ExtraSpecs"),
		),
	}
}

func fetchInstanceSecurityGroups(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	instance := parent.Item.(*Instance)

	//if instance.SecurityGroups != nil {
	for _, group := range instance.SecurityGroups {
		api.Logger().Debug().Str("instance id", instance.ID).Msg("streaming instance security group")
		res <- SecurityGroup{
			Name: group.Name,
		}
	}
	//}
	return nil
}

type SecurityGroup struct {
	Name string `json:"name"`
}

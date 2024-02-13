package neutron

import (
	"context"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
)

func SecurityGroups() *schema.Table {
	return &schema.Table{
		Name:     "openstack_neutron_security_groups",
		Resolver: fetchSecurityGroups,
		Transform: transformers.TransformWithStruct(
			&groups.SecGroup{},
			transformers.WithPrimaryKeys("ID"),
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
		Columns: []schema.Column{
			{
				Name:        "security_group_rule_ids",
				Type:        arrow.ListOf(arrow.BinaryTypes.String),
				Description: "The collection of IP addresses associated with the port.",
				Resolver: transform.Apply(
					transform.OnObjectField("Rules.ID"),
					transform.NilIfZero(),
				),
			},
		},
	}
}

func fetchSecurityGroups(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	neutron, err := api.GetServiceClient(client.NeutronV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := groups.ListOpts{}

	allPages, err := groups.List(neutron, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing security groups with options")
		return err
	}
	allSecurityGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting security groups")
		return err
	}
	api.Logger().Debug().Int("count", len(allSecurityGroups)).Msg("security groups retrieved")

	for _, group := range allSecurityGroups {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		group := group
		// api.Logger().Debug().Str("data", format.ToPrettyJSON(group)).Msg("streaming security group")
		api.Logger().Debug().Str("data", group.ID).Msg("streaming security group")
		res <- group
	}
	return nil
}

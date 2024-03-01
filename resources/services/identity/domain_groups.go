package identity

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/groups"
)

func DomainGroups() *schema.Table {
	return &schema.Table{
		Name:     "openstack_identity_domain_groups",
		Resolver: fetchDomainGroups,
		Transform: transformers.TransformWithStruct(
			&groups.Group{},
			transformers.WithSkipFields("Links"),
		),
	}
}

func fetchDomainGroups(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	domain := parent.Item.(domains.Domain)

	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := groups.ListOpts{
		DomainID: domain.ID,
	}

	allPages, err := groups.List(identity, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing groups with options")
		return err
	}
	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting groups")
		return err
	}
	api.Logger().Debug().Int("count", len(allGroups)).Msg("groups retrieved")

	for _, group := range allGroups {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		api.Logger().Debug().Str("user id", group.ID).Str("data", format.ToPrettyJSON(group)).Msg("streaming group")
		res <- group
	}
	return nil
}

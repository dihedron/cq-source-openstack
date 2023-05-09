package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
)

func Users() *schema.Table {
	return &schema.Table{
		Name:     "openstack_users",
		Resolver: fetchUsers,
		Transform: transformers.TransformWithStruct(
			&users.User{},
			transformers.WithPrimaryKeys("ID"),
			transformers.WithSkipFields("Links"),
		),
		// Columns: []schema.Column{
		// 	{
		// 		Name:        "tags",
		// 		Type:        schema.TypeStringArray,
		// 		Description: "The set of tags on the project.",
		// 		Resolver: transform.Apply(
		// 			transform.OnObjectField("Tags"),
		// 		),
		// 	},
		// },
	}
}

func fetchUsers(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	keystone, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := users.ListOpts{}

	allPages, err := users.List(keystone, opts).AllPages()
	if err != nil {
		api.Logger.Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing users with options")
		return err
	}
	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error extracting users")
		return err
	}
	api.Logger.Debug().Int("count", len(allUsers)).Msg("users retrieved")

	for _, user := range allUsers {
		if ctx.Err() != nil {
			api.Logger.Debug().Msg("context done, exit")
			break
		}
		user := user
		// api.Logger.Debug().Str("data", format.ToPrettyJSON(user)).Msg("streaming user")
		api.Logger.Debug().Str("data", user.ID).Msg("streaming user")
		res <- user
	}
	return nil
}

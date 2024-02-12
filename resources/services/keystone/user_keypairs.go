package keystone

import (
	"context"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
)

func UserKeyPairs() *schema.Table {
	return &schema.Table{
		Name:     "openstack_keystone_user_keypairs",
		Resolver: fetchUserKeyPairs,
		Transform: transformers.TransformWithStruct(
			&keypairs.KeyPair{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
	}
}

func fetchUserKeyPairs(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	user := parent.Item.(*User)

	compute, err := api.GetServiceClient(client.NovaV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := keypairs.ListOpts{
		UserID: user.ID,
	}

	allPages, err := keypairs.List(compute, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing keypairs with options")
		return err
	}
	allKeyPairs, err := keypairs.ExtractKeyPairs(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting keypairs")
		return err
	}
	api.Logger().Debug().Int("count", len(allKeyPairs)).Msg("keypairs retrieved")

	for _, keypair := range allKeyPairs {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		keypair.UserID = user.ID
		api.Logger().Debug().Str("user id", keypair.UserID).Str("data", format.ToPrettyJSON(keypair)).Msg("streaming keypair")
		res <- keypair
	}
	return nil
}

type KeyPair struct {
	ID          *int       `json:"ID"`
	Name        string     `json:"name"`
	Fingerprint string     `json:"fingerprint"`
	PublicKey   string     `json:"public_key"`
	PrivateKey  string     `json:"private_key"`
	UserID      *string    `json:"user_id"`
	Type        string     `json:"type"`
	Deleted     bool       `json:"deleted"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

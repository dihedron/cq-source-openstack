package blockstorage

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/backups"
)

func Backups() *schema.Table {
	return &schema.Table{
		Name:     "openstack_blockstorage_backups",
		Resolver: fetchBackups,
		Transform: transformers.TransformWithStruct(
			&backups.Backup{},
			transformers.WithPrimaryKeys("ID"),
		),
	}
}

func fetchBackups(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	cinder, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	listOpts := backups.ListOpts{
		VolumeID: "uuid",
	}

	allPages, err := backups.List(cinder, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allBackups, err := backups.ExtractBackups(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error getting backups")
		panic(err)
	}

	for _, backup := range allBackups {
		res <- backup
	}

	return nil
}

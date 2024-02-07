package blockstorage

import (
	"context"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/limits"
)

func Limits() *schema.Table {
	return &schema.Table{
		Name:     "openstack_blockstorage_limits",
		Resolver: fetchLimits,
		Transform: transformers.TransformWithStruct(
			&limits.Limit{},
			transformers.WithSkipFields("Absolute", "Rate"),
		),
		Columns: []schema.Column{
			{
				Name:        "max_total_volumes",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The maximum number of volumes that can be created.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalVolumes"),
			},
			{
				Name:        "max_total_snapshots",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The maximum number of snapshots that can be created.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalSnapshots"),
			},
			{
				Name:        "max_total_volume_gigabytes",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The maximum number of gigabytes that can be used for volumes.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalVolumeGigabytes"),
			},
			{
				Name:        "max_total_backups",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The maximum number of backups that can be created.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalBackups"),
			},
			{
				Name:        "max_total_backup_gigabytes",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The maximum number of gigabytes that can be used for backups.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalBackupGigabytes"),
			},
			{
				Name:        "total_volumes_used",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of volumes that have been created.",
				Resolver:    schema.PathResolver("Absolute.TotalVolumesUsed"),
			},
			{
				Name:        "total_gigabytes_used",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes that have been used for volumes.",
				Resolver:    schema.PathResolver("Absolute.TotalGigabytesUsed"),
			},
			{
				Name:        "total_snapshots_used",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of snapshots that have been created.",
				Resolver:    schema.PathResolver("Absolute.TotalSnapshotsUsed"),
			},
			{
				Name:        "total_backups_used",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backups that have been created.",
				Resolver:    schema.PathResolver("Absolute.TotalBackupsUsed"),
			},
			{
				Name:        "total_backup_gigabytes_used",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes that have been used for backups.",
				Resolver:    schema.PathResolver("Absolute.TotalBackupGigabytesUsed"),
			},
			{
				Name:        "regex",
				Type:        arrow.ListOf(arrow.BinaryTypes.String),
				Description: "The regular expression used to match the URI.",
				Resolver:    schema.PathResolver("Rate.Regex"),
			},
			{
				Name:        "uri",
				Type:        arrow.ListOf(arrow.BinaryTypes.String),
				Description: "The URI that the regular expression matches.",
				Resolver:    schema.PathResolver("Rate.URI"),
			},
			{
				Name:        "verb",
				Type:        arrow.ListOf(arrow.ListOf(arrow.BinaryTypes.String)),
				Description: "The HTTP verb used to match the URI.",
				Resolver:    schema.PathResolver("Rate.Limit.Verb"),
			},
			{
				Name:        "next_available",
				Type:        arrow.ListOf(arrow.ListOf(arrow.BinaryTypes.String)),
				Description: "The next available time for the rate limit.",
				Resolver:    schema.PathResolver("Rate.Limit.NextAvailable"),
			},
			{
				Name:        "unit",
				Type:        arrow.ListOf(arrow.ListOf(arrow.BinaryTypes.String)),
				Description: "The unit of the rate limit.",
				Resolver:    schema.PathResolver("Rate.Limit.Unit"),
			},
			{
				Name:        "value",
				Type:        arrow.ListOf(arrow.ListOf(arrow.PrimitiveTypes.Int64)),
				Description: "The value of the rate limit.",
				Resolver:    schema.PathResolver("Rate.Limit.Value"),
			},
			{
				Name:        "remaining",
				Type:        arrow.ListOf(arrow.ListOf(arrow.PrimitiveTypes.Int64)),
				Description: "The number of requests remaining in the current rate limit window.",
				Resolver:    schema.PathResolver("Rate.Limit.Remaining"),
			},
		},
	}
}

func fetchLimits(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	cinder, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	allLimits, err := limits.Get(cinder).Extract()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error getting limits")
		panic(err)
	}

	res <- allLimits
	return nil
}

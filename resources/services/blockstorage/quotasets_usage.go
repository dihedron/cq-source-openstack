package blockstorage

import (
	"context"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func QuotaSetsUsage() *schema.Table {
	return &schema.Table{
		Name:     "openstack_blockstorage_quotasets_usage",
		Resolver: fetchQuotaSetsUsage,
		Transform: transformers.TransformWithStruct(
			&quotasets.QuotaUsageSet{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithSkipFields("Volumes", "Snapshots", "Gigabytes", "PerVolumeGigabytes", "Backups", "BackupGigabytes", "Groups"),
		),
		Columns: []schema.Column{
			{
				Name:        "volumes_in_use",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of volumes currently in use.",
				Resolver:    schema.PathResolver("Volumes.InUse"),
			},
			{
				Name:        "volumes_allocated",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of volumes currently allocated.",
				Resolver:    schema.PathResolver("Volumes.Allocated"),
			},
			{
				Name:        "volumes_reserved",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of volumes currently reserved.",
				Resolver:    schema.PathResolver("Volumes.Reserved"),
			},
			{
				Name:        "volumes_limit",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of volumes currently allowed.",
				Resolver:    schema.PathResolver("Volumes.Limit"),
			},
			{
				Name:        "snapshots_in_use",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of snapshots currently in use.",
				Resolver:    schema.PathResolver("Snapshots.InUse"),
			},
			{
				Name:        "snapshots_allocated",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of snapshots currently allocated.",
				Resolver:    schema.PathResolver("Snapshots.Allocated"),
			},
			{
				Name:        "snapshots_reserved",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of snapshots currently reserved.",
				Resolver:    schema.PathResolver("Snapshots.Reserved"),
			},
			{
				Name:        "snapshots_limit",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of snapshots currently allowed.",
				Resolver:    schema.PathResolver("Snapshots.Limit"),
			},
			{
				Name:        "gigabytes_in_use",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes currently in use.",
				Resolver:    schema.PathResolver("Gigabytes.InUse"),
			},
			{
				Name:        "gigabytes_allocated",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes currently allocated.",
				Resolver:    schema.PathResolver("Gigabytes.Allocated"),
			},
			{
				Name:        "gigabytes_reserved",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes currently reserved.",
				Resolver:    schema.PathResolver("Gigabytes.Reserved"),
			},
			{
				Name:        "gigabytes_limit",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes currently allowed.",
				Resolver:    schema.PathResolver("Gigabytes.Limit"),
			},
			{
				Name:        "per_volume_gigabytes_in_use",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes currently in use per volume.",
				Resolver:    schema.PathResolver("PerVolumeGigabytes.InUse"),
			},
			{
				Name:        "per_volume_gigabytes_allocated",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes currently allocated per volume.",
				Resolver:    schema.PathResolver("PerVolumeGigabytes.Allocated"),
			},
			{
				Name:        "per_volume_gigabytes_reserved",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes currently reserved per volume.",
				Resolver:    schema.PathResolver("PerVolumeGigabytes.Reserved"),
			},
			{
				Name:        "per_volume_gigabytes_limit",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of gigabytes currently allowed per volume.",
				Resolver:    schema.PathResolver("PerVolumeGigabytes.Limit"),
			},
			{
				Name:        "backups_in_use",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backups currently in use.",
				Resolver:    schema.PathResolver("Backups.InUse"),
			},
			{
				Name:        "backups_allocated",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backups currently allocated.",
				Resolver:    schema.PathResolver("Backups.Allocated"),
			},
			{
				Name:        "backups_reserved",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backups currently reserved.",
				Resolver:    schema.PathResolver("Backups.Reserved"),
			},
			{
				Name:        "backups_limit",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backups currently allowed.",
				Resolver:    schema.PathResolver("Backups.Limit"),
			},
			{
				Name:        "backup_gigabytes_in_use",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backup gigabytes currently in use.",
				Resolver:    schema.PathResolver("BackupGigabytes.InUse"),
			},
			{
				Name:        "backup_gigabytes_allocated",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backup gigabytes currently allocated.",
				Resolver:    schema.PathResolver("BackupGigabytes.Allocated"),
			},
			{
				Name:        "backup_gigabytes_reserved",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backup gigabytes currently reserved.",
				Resolver:    schema.PathResolver("BackupGigabytes.Reserved"),
			},
			{
				Name:        "backup_gigabytes_limit",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of backup gigabytes currently allowed.",
				Resolver:    schema.PathResolver("BackupGigabytes.Limit"),
			},
			{
				Name:        "groups_in_use",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of groups currently in use.",
				Resolver:    schema.PathResolver("Groups.InUse"),
			},
			{
				Name:        "groups_allocated",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of groups currently allocated.",
				Resolver:    schema.PathResolver("Groups.Allocated"),
			},
			{
				Name:        "groups_reserved",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of groups currently reserved.",
				Resolver:    schema.PathResolver("Groups.Reserved"),
			},
			{
				Name:        "groups_limit",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of groups currently allowed.",
				Resolver:    schema.PathResolver("Groups.Limit"),
			},
		},
	}
}

func fetchQuotaSetsUsage(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	// get list of all projects
	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving identity client")
		return err
	}

	allPages, err := projects.List(identity, &projects.ListOpts{}).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error listing projects")
		return err
	}
	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting projects")
		return err
	}
	api.Logger().Debug().Int("count", len(allProjects)).Msg("projects retrieved")

	projectIDs := []string{}
	for _, project := range allProjects {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		projectIDs = append(projectIDs, project.ID)
	}

	blockstorage, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving blockstorage client")
		return err
	}

	// for each project, get the associated QuotaUsageSet
	for _, projectID := range projectIDs {
		quotausageset, err := quotasets.GetUsage(blockstorage, projectID).Extract()
		if err != nil {
			api.Logger().Error().Err(err).Msg("error extracting quota sets for project " + projectID)
			return err
		}
		res <- quotausageset
	}
	return nil
}

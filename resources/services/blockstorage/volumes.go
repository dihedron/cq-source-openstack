package blockstorage

import (
	"context"

	"github.com/dihedron/cq-source-openstack/resources/internal/utils"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func Volumes() *schema.Table {
	return &schema.Table{
		Name:     "openstack_blockstorage_volumes",
		Resolver: fetchVolumes,
		Transform: transformers.TransformWithStruct(
			&Volume{},
			transformers.WithPrimaryKeys("ID"),
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type

			transformers.WithSkipFields("Links"),
		),
		Relations: []*schema.Table{
			VolumesBackups(),
		},
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

func fetchVolumes(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	blockstorage, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := volumes.ListOpts{
		AllTenants: true,
	}

	allPages, err := volumes.List(blockstorage, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing volumes with options")
		return err
	}
	allVolumes := []*Volume{}
	if err := volumes.ExtractVolumesInto(allPages, &allVolumes); err != nil {
		api.Logger().Error().Err(err).Msg("error extracting volumes")
		return err
	}
	api.Logger().Debug().Int("count", len(allVolumes)).Msg("volumes retrieved")

	for _, volume := range allVolumes {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		volume := volume
		// api.Logger().Debug().Str("data", format.ToPrettyJSON(volume)).Msg("streaming volume")
		api.Logger().Debug().Str("data", volume.ID).Msg("streaming volume")
		res <- volume
	}
	return nil
}

type Volume struct {
	// Unique identifier for the volume.
	ID string `json:"id"`
	// Current status of the volume.
	Status string `json:"status"`
	// Size of the volume in GB.
	Size int `json:"size"`
	// AvailabilityZone is which availability zone the volume is in.
	AvailabilityZone string `json:"availability_zone"`
	// The date when this volume was created.
	CreatedAt utils.Time `json:"created_at"`
	// The date when this volume was last updated
	UpdatedAt utils.Time `json:"updated_at"`
	// Instances onto which the volume is attached.
	Attachments []volumes.Attachment `json:"attachments"`
	// Human-readable display name for the volume.
	Name string `json:"name"`
	// Human-readable description for the volume.
	Description string `json:"description"`
	// The type of volume to create, either SATA or SSD.
	VolumeType string `json:"volume_type"`
	// The ID of the snapshot from which the volume was created
	SnapshotID string `json:"snapshot_id"`
	// The ID of another block storage volume from which the current volume was created
	SourceVolID string `json:"source_volid"`
	// The backup ID, from which the volume was restored
	// This field is supported since 3.47 microversion
	BackupID *string `json:"backup_id"`
	// The group ID; this field is supported since 3.47 microversion
	GroupID *string `json:"group_id"`
	// Arbitrary key-value pairs defined by the user.
	Metadata map[string]string `json:"metadata"`
	// UserID is the id of the user who created the volume.
	UserID string `json:"user_id"`
	// Indicates whether this is a bootable volume.
	Bootable string `json:"bootable"`
	// Encrypted denotes if the volume is encrypted.
	Encrypted bool `json:"encrypted"`
	// ReplicationStatus is the status of replication.
	ReplicationStatus string `json:"replication_status"`
	// ConsistencyGroupID is the consistency group ID.
	ConsistencyGroupID string `json:"consistencygroup_id"`
	// Multiattach denotes if the volume is multi-attach capable.
	Multiattach bool `json:"multiattach"`
	// Image metadata entries, only included for volumes that were created from an image, or from a snapshot of a volume originally created from an image.
	VolumeImageMetadata map[string]string `json:"volume_image_metadata"`
	// The volume migration status
	MigrationStatus string `json:"migration_status"`

	OsVolHostAttrHost         string `json:"os-vol-host-attr:host" cq-name:"host"`
	OsVolMigStatusAttrMigstat string `json:"os-vol-mig-status-attr:migstat" cq-name:"migration_status"`
	OsVolMigStatusAttrNameID  string `json:"os-vol-mig-status-attr:name_id" cq-name:"migration_status_name"`
	OsVolTenantAttrTenantID   string `json:"os-vol-tenant-attr:tenant_id" cq-name:"migration_status_tenant"`
	ProviderID                string `json:"provider_id"`
	ServiceUUID               string `json:"service_uuid"`
	SharedTargets             bool   `json:"shared_targets"`
}

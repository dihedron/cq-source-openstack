package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/attachments"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func Attachments() *schema.Table {
	return &schema.Table{
		Name:     "openstack_attachments",
		Resolver: fetchAttachments,
		Transform: transformers.TransformWithStruct(
			&Attachment{},
			transformers.WithPrimaryKeys("ID"),
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
		Columns: []schema.Column{
			{
				Name:        "attached_at",
				Type:        schema.TypeTimestamp,
				Description: "The time at which the attachment was created.",
				Resolver: transform.Apply(
					transform.OnObjectField("AttachedAt"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "detached_at",
				Type:        schema.TypeTimestamp,
				Description: "The time at which the attachment was removed.",
				Resolver: transform.Apply(
					transform.OnObjectField("DetachedAt"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "access_mode",
				Type:        schema.TypeString,
				Description: "The access mode of the volume attachment.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.AccessMode"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "attach_mode",
				Type:        schema.TypeString,
				Description: "The mode in which the volume is attached.",
				Resolver: transform.Apply(
					transform.OnObjectField("AttachMode"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "attachment_id",
				Type:        schema.TypeString,
				Description: "The ID of the volume attachment.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.AttachmentID"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "auth_enabled",
				Type:        schema.TypeBool,
				Description: "Whether the attachment authorisation is enabled.",
				Resolver:    schema.PathResolver("ConnectionInfo.AuthEnabled"),
			},
			{
				Name:        "auth_username",
				Type:        schema.TypeString,
				Description: "The name of the user that attached the volume.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.AuthUsername"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "cluster_name",
				Type:        schema.TypeString,
				Description: "The name of the cluster.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.ClusterName"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "discard",
				Type:        schema.TypeBool,
				Description: "The name of the cluster.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.Discard"),
				),
			},
			{
				Name:        "driver_volume_type",
				Type:        schema.TypeString,
				Description: "The type of the driver.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.DriverVolumeType"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "encrypted",
				Type:        schema.TypeBool,
				Description: "Whether the volume is encrypted.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.Encrypted"),
				),
			},
			{
				Name:        "hosts",
				Type:        schema.TypeStringArray,
				Description: "The storage hosts that hold the data in the volume.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.Hosts"),
				),
			},
			{
				Name:        "keyring",
				Type:        schema.TypeBool,
				Description: "The keyring associated with the attachment.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.Keyring"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "name",
				Type:        schema.TypeString,
				Description: "The name of the attachment.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.Name"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "ports",
				Type:        schema.TypeStringArray,
				Description: "The ports of the attachment.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.Ports"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "secret_type",
				Type:        schema.TypeString,
				Description: "The ports of the attachment.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.SecretType"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "secret_uuid",
				Type:        schema.TypeString,
				Description: "The ports of the attachment.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.SecretUUID"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "volume_id",
				Type:        schema.TypeString,
				Description: "The ports of the attachment.",
				Resolver: transform.Apply(
					transform.OnObjectField("ConnectionInfo.VolumeID"),
					transform.NilIfZero(),
				),
			},
		},
	}
}

func fetchAttachments(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	// get list of all projects
	keystone, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error retrieving keystone client")
		return err
	}

	allPages, err := projects.List(keystone, &projects.ListOpts{}).AllPages()
	if err != nil {
		api.Logger.Error().Err(err).Msg("error listing projects")
		return err
	}
	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error extracting projects")
		return err
	}
	api.Logger.Debug().Int("count", len(allProjects)).Msg("projects retrieved")

	projectIDs := []string{}
	for _, project := range allProjects {
		if ctx.Err() != nil {
			api.Logger.Debug().Msg("context done, exit")
			break
		}
		project := project
		projectIDs = append(projectIDs, project.ID)
	}

	cinder, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error retrieving cinder client")
		return err
	}

	// for each project, get the associated attachments
	opts := attachments.ListOpts{
		AllTenants: true,
	}
	for _, projectID := range projectIDs {
		opts.ProjectID = projectID

		allPages, err := attachments.List(cinder, opts).AllPages()
		if err != nil {
			api.Logger.Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing attachments with options")
			return err
		}
		allAttachments := []*Attachment{}
		err = attachments.ExtractAttachmentsInto(allPages, &allAttachments)
		if err != nil {
			api.Logger.Error().Err(err).Msg("error extracting attachments")
			return err
		}
		api.Logger.Debug().Int("count", len(allAttachments)).Msg("attachments retrieved")

		for _, attachment := range allAttachments {
			if ctx.Err() != nil {
				api.Logger.Debug().Msg("context done, exit")
				break
			}
			attachment := attachment
			attachment.ProjectID = projectID
			api.Logger.Debug().Str("data", format.ToPrettyJSON(attachment)).Msg("streaming attachment")
			res <- attachment
		}
	}
	return nil
}

type Attachment struct {
	ID             string `json:"id"`
	AttachedAt     Time   `json:"attached_at"`
	DetachedAt     Time   `json:"detached_at"`
	AttachmentID   string `json:"attachment_id"`
	VolumeID       string `json:"volume_id"`
	Instance       string `json:"instance" cq-name:"instance_id"`
	Status         string `json:"status"`
	AttachMode     string `json:"attach_mode"`
	ProjectID      string `json:"-" cq-name:"project_id"`
	ConnectionInfo struct {
		AccessMode       string   `json:"access_mode"`
		AttachmentID     string   `json:"attachment_id"`
		AuthEnabled      bool     `json:"auth_enabled"`
		AuthUsername     string   `json:"auth_username"`
		ClusterName      string   `json:"cluster_name"`
		Discard          bool     `json:"discard"`
		DriverVolumeType string   `json:"driver_volume_type"`
		Encrypted        bool     `json:"encrypted"`
		Hosts            []string `json:"hosts"`
		Keyring          string   `json:"keyring"`
		Name             string   `json:"name"`
		Ports            []string `json:"ports"`
		SecretType       string   `json:"secret_type"`
		SecretUUID       string   `json:"secret_uuid"`
		VolumeID         string   `json:"volume_id"`
	} `json:"connection_info"`
}

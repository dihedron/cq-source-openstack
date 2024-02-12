package keystone

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/dihedron/cq-source-openstack/resources/services/nova"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func Projects() *schema.Table {
	return &schema.Table{
		Name:     "openstack_keystone_projects",
		Resolver: fetchProjects,
		Transform: transformers.TransformWithStruct(
			&projects.Project{},
			transformers.WithPrimaryKeys("ID"),
			transformers.WithSkipFields("Links"),
		),
		Relations: []*schema.Table{
			nova.ProjectLimits(),
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

func fetchProjects(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	keystone, err := api.GetServiceClient(client.KeystoneV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := projects.ListOpts{}

	allPages, err := projects.List(keystone, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing projects with options")
		return err
	}
	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting projects")
		return err
	}
	api.Logger().Debug().Int("count", len(allProjects)).Msg("projects retrieved")

	for _, project := range allProjects {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		project := project
		//api.Logger().Debug().Str("data", format.ToPrettyJSON(project)).Msg("streaming project")
		api.Logger().Debug().Str("id", project.ID).Msg("streaming project")
		res <- project
	}
	return nil
}

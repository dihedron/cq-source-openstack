package blockstorage

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func QuotaSets() *schema.Table {
	return &schema.Table{
		Name:     "openstack_blockstorage_quotasets",
		Resolver: fetchQuotaSets,
		Transform: transformers.TransformWithStruct(
			&quotasets.QuotaSet{},
			// transformers.WithSkipFields("Extra"),
		),
	}
}

func fetchQuotaSets(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

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
		project := project
		projectIDs = append(projectIDs, project.ID)
	}

	blockstorage, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving blockstorage client")
		return err
	}

	// for each project, get the associated QuotaSets
	for _, projectID := range projectIDs {
		quotaset, err := quotasets.Get(blockstorage, projectID).Extract()
		if err != nil {
			api.Logger().Error().Err(err).Msg("error extracting quota sets for project " + projectID)
			return err
		}
		res <- quotaset
	}
	return nil
}

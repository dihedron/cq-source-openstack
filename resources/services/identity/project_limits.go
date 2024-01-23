package identity

import (
	"context"
	"errors"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/limits"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func ProjectLimits() *schema.Table {
	return &schema.Table{
		Name:     "openstack_project_limits",
		Resolver: fetchProjectLimits,
		Transform: transformers.TransformWithStruct(
			&limits.Absolute{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
		Relations: []*schema.Table{},
	}
}

func fetchProjectLimits(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	project := parent.Item.(projects.Project)

	compute, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	if ctx.Err() != nil {
		api.Logger().Debug().Msg("context done, exit")
		return errors.New("interrupted due to context done")
	}

	opts := limits.GetOpts{
		TenantID: project.ID,
	}
	allLimits, err := limits.Get(compute, opts).Extract()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing limits with options")
		return err
	}
	api.Logger().Debug().Str("project id", project.ID).Msg("streaming project limits")
	res <- allLimits.Absolute
	return nil
}

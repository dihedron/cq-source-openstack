package compute

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

func FlavorAccesses() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_flavor_accesses",
		Resolver: fetchFlavorAccesses,
		Transform: transformers.TransformWithStruct(
			&FlavorAccess{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
		),
	}
}

func fetchFlavorAccesses(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	nova, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	flavor := parent.Item.(*Flavor)
	api.Logger().Debug().Str("flavor id", *flavor.ID).Msg("retrieving accesses for flavor")

	allPages, err := flavors.ListAccesses(nova, *flavor.ID).AllPages()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			api.Logger().Warn().Err(err).Str("flavor id", *flavor.ID).Msg("no flavor accesses for flavor")
			return nil
		} else {
			api.Logger().Error().Err(err).Str("flavor id", *flavor.ID).Str("err type", fmt.Sprintf("%T", err)).Msg("error listing flavor accesses for flavor")
			return err
		}
	}
	api.Logger().Debug().Msg("flavor accesses retrieved")

	allAccesses := []*FlavorAccess{}
	if err = ExtractFlavorAccessInto(allPages, &allAccesses); err != nil {
		api.Logger().Error().Err(err).Msg("error extracting flavor accesses")
		return err
	}
	api.Logger().Debug().Int("count", len(allAccesses)).Msg("flavors accesses retrieved")
	for _, access := range allAccesses {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		access := access
		api.Logger().Debug().Str("flavor id", access.FlavorID).Str("project id", access.TenantID).Msg("streaming flavor access")
		res <- access
	}

	return nil

}

// type AutoGenerated struct {
// 	FlavorAccess []FlavorAccess `json:"flavor_access"`
// }

func ExtractFlavorAccessInto(r pagination.Page, v interface{}) error {
	return r.(flavors.AccessPage).Result.ExtractIntoSlicePtr(v, "flavor_access")
}

type FlavorAccess struct {
	FlavorID string `json:"flavor_id"`
	TenantID string `json:"tenant_id" cq-name:"project_id"`
}

package compute

import (
	"context"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/dihedron/cq-source-openstack/resources/internal/utils"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/usage"
	"github.com/gophercloud/gophercloud/pagination"
)

func ServerUsage() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_serverusage",
		Resolver: fetchServerUsage,
		Transform: transformers.TransformWithStruct(
			&Usage{},
		),
		Columns: []schema.Column{
			{
				Name:        "started_at",
				Type: 	 	 arrow.FixedWidthTypes.Timestamp_us,
				Description: "The time when the server usage started.",
				Resolver:	 schema.PathResolver("StartedAt"),
			},
			{
				Name:        "ended_at",
				Type: 	 	 arrow.FixedWidthTypes.Timestamp_us,
				Description: "The time when the server usage ended.",
				Resolver:	 schema.PathResolver("EndedAt"),
			},
		},
	}
}

func fetchServerUsage(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	nova, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	allTenantsOpts := usage.AllTenantsOpts{
		Detailed: true,
	}

	err = usage.AllTenants(nova, allTenantsOpts).EachPage(func(page pagination.Page) (bool, error) {
		allTenantsUsage, err := usage.ExtractAllTenants(page)
		if err != nil {
			return false, err
		}
	
		res <- allTenantsUsage
	
		return true, nil
	})
	
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting all tenants usage")
		return err
	}

	return nil
}

type Usage struct {
	EndedAt utils.Time `json:"-"`
	Flavor string `json:"flavor"`
	Hours float64 `json:"hours"`
	InstanceID string `json:"instance_id"`
	LocalGB int `json:"local_gb"`
	MemoryMB int `json:"memory_mb"`
	Name string `json:"name"`
	StartedAt utils.Time `json:"-"`
	State string `json:"state"`
	TenantID string `json:"tenant_id"`
	Uptime int `json:"uptime"`
	VCPUs int `json:"vcpus"`
}

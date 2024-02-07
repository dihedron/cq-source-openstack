package compute

import (
	"context"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/limits"
)

func Limits() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_limits",
		Resolver: fetchLimits,
		Transform: transformers.TransformWithStruct(
			&Limit{},
			transformers.WithSkipFields("Absolute"),
		),
		Columns: []schema.Column{
			{
				Name:        "max_total_cores",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of cores available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalCores"),
			},
			{
				Name:        "max_image_meta",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The amount of image metadata available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxImageMeta"),
			},
			{
				Name:        "max_server_meta",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The amount of server metadata available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxServerMeta"),
			},
			{
				Name:        "max_personality",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The amount of personality/files available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxPersonality"),
			},
			{
				Name:        "max_personality_size",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The personality file size available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxPersonalitySize"),
			},
			{
				Name:        "max_security_group_rules",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of security group rules available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxSecurityGroupRules"),
			},
			{
				Name:        "max_security_groups",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of security groups available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxSecurityGroups"),
			},
			{
				Name:        "max_server_groups",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of server groups available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxServerGroups"),
			},
			{
				Name:        "max_server_group_members",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of server group members available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxServerGroupMembers"),
			},
			{
				Name:        "max_total_floating_ips",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of floating IPs available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalFloatingIps"),
			},
			{
				Name:        "max_total_instances",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of instances/servers available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalInstances"),
			},
			{
				Name:        "max_total_keypairs",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The total keypairs available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalKeypairs"),
			},
			{
				Name:        "max_total_ram_size",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The amount of RAM available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalRAMSize"),
			},
			{
				Name:        "max_total_security_groups",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of security groups available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalSecurityGroups"),
			},
			{
				Name:        "max_total_snapshots",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of snapshots available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalSnapshots"),
			},
			{
				Name:        "max_total_volumes",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The number of volumes available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalVolumes"),
			},
			{
				Name:        "max_total_volume_gigabytes",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The amount of volume gigabytes available to a tenant.",
				Resolver:    schema.PathResolver("Absolute.MaxTotalVolumeGigabytes"),
			},
		},
	}
}

func fetchLimits(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	nova, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := limits.GetOpts{
		TenantID: "tenant-id",
	}

	limits, err := limits.Get(nova, opts).Extract()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving limits")
		return err
	}

	res <- limits
	return nil
}

type Limit struct {
	Absolute struct {
		// MaxTotalCores is the number of cores available to a tenant.
		MaxTotalCores int `json:"maxTotalCores"`
		// MaxImageMeta is the amount of image metadata available to a tenant.
		MaxImageMeta int `json:"maxImageMeta"`
		// MaxServerMeta is the amount of server metadata available to a tenant.
		MaxServerMeta int `json:"maxServerMeta"`
		// MaxPersonality is the amount of personality/files available to a tenant.
		MaxPersonality int `json:"maxPersonality"`
		// MaxPersonalitySize is the personality file size available to a tenant.
		MaxPersonalitySize int `json:"maxPersonalitySize"`
		// MaxTotalKeypairs is the total keypairs available to a tenant.
		MaxTotalKeypairs int `json:"maxTotalKeypairs"`
		// MaxSecurityGroups is the number of security groups available to a tenant.
		MaxSecurityGroups int `json:"maxSecurityGroups"`
		// MaxSecurityGroupRules is the number of security group rules available to
		// a tenant.
		MaxSecurityGroupRules int `json:"maxSecurityGroupRules"`
		// MaxServerGroups is the number of server groups available to a tenant.
		MaxServerGroups int `json:"maxServerGroups"`
		// MaxServerGroupMembers is the number of server group members available
		// to a tenant.
		MaxServerGroupMembers int `json:"maxServerGroupMembers"`
		// MaxTotalFloatingIps is the number of floating IPs available to a tenant.
		MaxTotalFloatingIps int `json:"maxTotalFloatingIps"`
		// MaxTotalInstances is the number of instances/servers available to a tenant.
		MaxTotalInstances int `json:"maxTotalInstances"`
		// MaxTotalRAMSize is the total amount of RAM available to a tenant measured
		// in megabytes (MB).
		MaxTotalRAMSize int `json:"maxTotalRAMSize"`
		// TotalCoresUsed is the number of cores currently in use.
		TotalCoresUsed int `json:"totalCoresUsed"`
		// TotalInstancesUsed is the number of instances/servers in use.
		TotalInstancesUsed int `json:"totalInstancesUsed"`
		// TotalFloatingIpsUsed is the number of floating IPs in use.
		TotalFloatingIpsUsed int `json:"totalFloatingIpsUsed"`
		// TotalRAMUsed is the total RAM/memory in use measured in megabytes (MB).
		TotalRAMUsed int `json:"totalRAMUsed"`
		// TotalSecurityGroupsUsed is the total number of security groups in use.
		TotalSecurityGroupsUsed int `json:"totalSecurityGroupsUsed"`
		// TotalServerGroupsUsed is the total number of server groups in use.
		TotalServerGroupsUsed int `json:"totalServerGroupsUsed"`
	} `json:"absolute"`
}

package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/dihedron/cq-source-openstack/format"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func Instances() *schema.Table {
	return &schema.Table{
		Name:      "openstack_instances",
		Resolver:  fetchInstances,
		Transform: transformers.TransformWithStruct(&Instance{}, transformers.WithPrimaryKeys("ID")),
		Columns: []schema.Column{
			{
				Name: "column",
				Type: schema.TypeString,
			},
			// {
			// 	Name:     "launched_at",
			// 	Type:     schema.TypeTimestamp,
			// 	Resolver: resolveLaunchedAt,
			// },
		},
	}
}

func fetchInstances(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	compute, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := servers.ListOpts{
		AllTenants: true,
	}

	allPages, err := servers.List(compute, opts).AllPages()
	if err != nil {
		api.Logger.Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing instances with options")
		return err
	}
	allInstances := []*Instance{}
	err = servers.ExtractServersInto(allPages, &allInstances)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error extracting instances")
		return err
	}
	api.Logger.Debug().Int("count", len(allInstances)).Msg("instances retrieved")

	for _, instance := range allInstances {
		if ctx.Err() != nil {
			api.Logger.Debug().Msg("context done, exit")
			break
		}
		instance := instance
		api.Logger.Debug().Str("data", format.ToPrettyJSON(instance)).Msg("streaming instance")
		res <- instance
	}
	return nil
}

// func resolveLaunchedAt(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
// 	instance := resource.Item.(Instance)
// 	resource.Set(c.Name, instance.LaunchedAt)
// 	return nil
// }

// Instance is an internal type used to unmarshal more data from the API
// response than would usually be possible through the ordinary gophercloud
// struct. OpenStack API microversions enable more response data that is not
// taken into account by the gophercloud library, which unmarshals only what
// is available at the base level for each API version, for backward compatibility.
// This is also why there is an ExtractInto function that allows you to pass in
// an arbitrary struct to marshal the response data into.
type Instance struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt Time   `json:"created"`
	// LaunchedAt   Time   `json:"OS-SRV-USG:launched_at"`
	UpdatedAt Time `json:"updated"`
	// TerminatedAt Time   `json:"OS-SRV-USG:terminated_at"`
	HostID     string `json:"hostid"`
	Status     string `json:"status"`
	Progress   int    `json:"progress"`
	AccessIPv4 string `json:"accessIPv4"`
	AccessIPv6 string `json:"accessIPv6"`
	Image      any    `json:"image"`
	Flavor     struct {
		Disk       int `json:"disk"`
		Ephemeral  int `json:"ephemeral"`
		ExtraSpecs struct {
			CPUCores        string `json:"hw:cpu_cores"`
			CPUSockets      string `json:"hw:cpu_sockets"`
			RNGAllowed      string `json:"hw_rng:allowed"`
			WatchdogAction  string `json:"hw:watchdog_action"`
			VGPUs           string `json:"resources:VGPU"`
			TraitCustomVGPU string `json:"trait:CUSTOM_VGPU"`
		} `json:"extra_specs"`
		OriginalName string `json:"original_name"`
		RAM          int    `json:"ram"`
		Swap         int    `json:"swap"`
		VCPUs        int    `json:"vcpus"`
	} `json:"flavor"`
	Addresses map[string][]struct {
		MACAddress string `json:"OS-EXT-IPS-MAC:mac_addr"`
		IPType     string `json:"OS-EXT-IPS:type"`
		IPAddress  string `json:"addr"`
		IPVersion  int    `json:"version"`
	} `json:"addresses"`
	Metadata map[string]string `json:"metadata"`
	Links    []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	KeyName        string `json:"key_name"`
	AdminPass      string `json:"adminPass"`
	SecurityGroups []struct {
		Name string `json:"name"`
	} `json:"security_groups"`
	// AttachedVolumes []servers.AttachedVolume `json:"os-extended-volumes:volumes_attached"`
	// NO!!! Fault              servers.Fault            `json:"fault"`
	Tags         *[]string `json:"tags"`
	ServerGroups *[]string `json:"server_groups"`
	// DiskConfig         string    `json:"OS-DCF:diskConfig"`
	// AvailabilityZone   string    `json:"OS-EXT-AZ:availability_zone"`
	// Host               string    `json:"OS-EXT-SRV-ATTR:host"`
	// HostName           string    `json:"OS-EXT-SRV-ATTR:hostname"`
	// HypervisorHostname string    `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	// InstanceName       string    `json:"OS-EXT-SRV-ATTR:instance_name"`
	// KernelID           string    `json:"OS-EXT-SRV-ATTR:kernel_id"`
	// LaunchIndex        int       `json:"OS-EXT-SRV-ATTR:launch_index"`
	// RAMDiskID          string    `json:"OS-EXT-SRV-ATTR:ramdisk_id"`
	// ReservationID      string    `json:"OS-EXT-SRV-ATTR:reservation_id"`
	// RootDeviceName     string    `json:"OS-EXT-SRV-ATTR:root_device_name"`
	// UserData           string    `json:"OS-EXT-SRV-ATTR:user_data"`
	// PowerState         int       `json:"OS-EXT-STS:power_state"`
	// VMState            string    `json:"OS-EXT-STS:vm_state"`
	ConfigDrive string `json:"config_drive"`
	Description string `json:"description"`
	//	NO!!! TaskState          interface{}              `json:"OS-EXT-STS:task_state"`
}

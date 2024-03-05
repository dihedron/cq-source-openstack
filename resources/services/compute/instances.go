package compute

import (
	"context"

	"github.com/dihedron/cq-plugin-utils/utils"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func Instances() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_instances",
		Resolver: fetchInstances,
		Transform: transformers.TransformWithStruct(
			&Instance{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
		Relations: []*schema.Table{
			InstanceAddresses(),
			InstanceAttachedVolumes(),
			InstanceFlavors(),
			InstanceFlavorExtraSpecs(),
			InstanceMetadata(),
			InstanceSecurityGroups(),
			InstanceTags(),
		},
		Columns: []schema.Column{
			{
				Name:        "image_id",
				Type:        arrow.BinaryTypes.String,
				Description: "The Image image used to start the instance.",
				Resolver: transform.Apply(
					transform.OnObjectField("Image"),
					transform.GetMapEntry[string, any]("id"),
					transform.TrimString(),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "power_state_name",
				Type:        arrow.BinaryTypes.String,
				Description: "The instance power state as a string.",
				Resolver: transform.Apply(
					transform.OnObjectField("PowerState"),
					transform.RemapValue(map[int]string{
						0: "NOSTATE",
						1: "RUNNING",
						3: "PAUSED",
						4: "SHUTDOWN",
						6: "CRASHED",
						7: "SUSPENDED",
					}),
				),
			},
		},
	}
}

func fetchInstances(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	compute, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := servers.ListOpts{
		AllTenants: true,
	}

	allPages, err := servers.List(compute, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing instances with options")
		return err
	}
	allInstances := []*Instance{}
	if err = servers.ExtractServersInto(allPages, &allInstances); err != nil {
		api.Logger().Error().Err(err).Msg("error extracting instances")
		return err
	}
	api.Logger().Debug().Int("count", len(allInstances)).Msg("instances retrieved")

	for _, instance := range allInstances {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		instance := instance
		api.Logger().Debug().Str("id", instance.ID).Msg("streaming instance")
		res <- instance
	}
	return nil
}

// Instance is an internal type used to unmarshal more data from the API
// response than would usually be possible through the ordinary gophercloud
// struct. OpenStack API microversions enable more response data that is not
// taken into account by the gophercloud library, which unmarshals only what
// is available at the base level for each API version, for backward compatibility.
// This is also why there is an ExtractInto function that allows you to pass in
// an arbitrary struct to marshal the response data into.
type Instance struct {
	ID           string               `json:"id"`
	TenantID     string               `json:"tenant_id"`
	UserID       string               `json:"user_id"`
	Name         string               `json:"name"`
	CreatedAt    *utils.Time          `json:"created" cq-name:"created_at" cq-type:"timestamp"`
	LaunchedAt   *utils.Time          `json:"OS-SRV-USG:launched_at" cq-name:"launched_at" cq-type:"timestamp"`
	UpdatedAt    *utils.Time          `json:"updated" cq-name:"updated_at" cq-type:"timestamp"`
	TerminatedAt *utils.Time          `json:"OS-SRV-USG:terminated_at" cq-name:"terminated_at" cq-type:"timestamp"`
	HostID       string               `json:"hostid"`
	Status       string               `json:"status"`
	Progress     int                  `json:"progress"`
	AccessIPv4   string               `json:"accessIPv4"`
	AccessIPv6   string               `json:"accessIPv6"`
	Image        any                  `json:"image"`
	Flavor       InstanceFlavor       `json:"flavor"`
	Addresses    map[string][]Address `json:"addresses"`
	Metadata     map[string]string    `json:"metadata"`
	Links        []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	KeyName        string `json:"key_name"`
	AdminPass      string `json:"adminPass"`
	SecurityGroups []struct {
		Name string `json:"name"`
	} `json:"security_groups"`
	AttachedVolumes    []servers.AttachedVolume `json:"os-extended-volumes:volumes_attached" cq-name:"attached_volumes"`
	Tags               *[]string                `json:"tags"`
	ServerGroups       *[]string                `json:"server_groups"`
	DiskConfig         string                   `json:"OS-DCF:diskConfig" cq-name:"disk_config"`
	AvailabilityZone   string                   `json:"OS-EXT-AZ:availability_zone" cq-name:"availability_zone"`
	Host               string                   `json:"OS-EXT-SRV-ATTR:host" cq-name:"host"`
	HostName           string                   `json:"OS-EXT-SRV-ATTR:hostname" cq-name:"hostname"`
	HypervisorHostname string                   `json:"OS-EXT-SRV-ATTR:hypervisor_hostname" cq-name:"hypervisor_hostname"`
	InstanceName       string                   `json:"OS-EXT-SRV-ATTR:instance_name" cq-name:"instance_name"`
	KernelID           string                   `json:"OS-EXT-SRV-ATTR:kernel_id" cq-name:"kernel_id"`
	LaunchIndex        int                      `json:"OS-EXT-SRV-ATTR:launch_index" cq-name:"launch_index"`
	RAMDiskID          string                   `json:"OS-EXT-SRV-ATTR:ramdisk_id" cq-name:"ramdisk_id"`
	ReservationID      string                   `json:"OS-EXT-SRV-ATTR:reservation_id" cq-name:"reservation_id"`
	RootDeviceName     string                   `json:"OS-EXT-SRV-ATTR:root_device_name" cq-name:"root_device_name"`
	UserData           string                   `json:"OS-EXT-SRV-ATTR:user_data" cq-name:"user_data"`
	PowerState         int                      `json:"OS-EXT-STS:power_state" cq-name:"power_state_id"`
	VMState            string                   `json:"OS-EXT-STS:vm_state" cq-name:"vm_state"`
	ConfigDrive        string                   `json:"config_drive"`
	Description        string                   `json:"description"`
	// NO!!! Fault              servers.Fault            `json:"fault"`
	// NO!!! TaskState          interface{}              `json:"OS-EXT-STS:task_state"`
}

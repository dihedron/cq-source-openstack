package client

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/dihedron/cq-source-openstack/format"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/rs/zerolog"
)

type Client struct {
	Logger   zerolog.Logger
	Client   *gophercloud.ProviderClient
	Services map[ServiceType]*gophercloud.ServiceClient
	Specs    *Spec
}

func (c *Client) ID() string {
	return "github.com/dihedron/cq-source-openstack"
}

func New(ctx context.Context, logger zerolog.Logger, s specs.Source, opts source.Options) (schema.ClientMeta, error) {
	var pluginSpec Spec

	if err := s.UnmarshalSpec(&pluginSpec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin spec: %w", err)
	}

	// TODO: Add your client initialization here
	logger.Debug().Str("spec", format.ToJSON(pluginSpec)).Msg("Plugin Spec")

	auth := gophercloud.AuthOptions{
		AllowReauth: true,
	}

	if pluginSpec.EndpointUrl != nil {
		auth.IdentityEndpoint = *pluginSpec.EndpointUrl
	}
	if pluginSpec.UserID != nil {
		auth.UserID = *pluginSpec.UserID
	}
	if pluginSpec.Username != nil {
		auth.Username = *pluginSpec.Username
	}
	if pluginSpec.Password != nil {
		auth.Password = *pluginSpec.Password
	}
	if pluginSpec.ProjectID != nil {
		auth.TenantID = *pluginSpec.ProjectID
	}
	if pluginSpec.ProjectName != nil {
		auth.TenantName = *pluginSpec.ProjectName
	}
	if pluginSpec.DomainID != nil {
		auth.DomainID = *pluginSpec.DomainID
	}
	if pluginSpec.DomainName != nil {
		auth.DomainName = *pluginSpec.DomainName
	}
	if pluginSpec.AccessToken != nil {
		auth.TokenID = *pluginSpec.AccessToken
	}
	if pluginSpec.AppCredentialID != nil {
		auth.ApplicationCredentialID = *pluginSpec.AppCredentialID
	}
	if pluginSpec.AppCredentialSecret != nil {
		auth.ApplicationCredentialSecret = *pluginSpec.AppCredentialSecret
	}
	if pluginSpec.AllowReauth != nil {
		auth.AllowReauth = *pluginSpec.AllowReauth
	}

	//
	// IMPORTANT NOTE: when using App Credentials, it is necessary
	// that all other fields except the endpoint URL be left blank!
	//

	client, err := openstack.AuthenticatedClient(auth)
	if err != nil {
		logger.Error().Err(err).Msg("error creating authenticated client")
		return nil, fmt.Errorf("error creating authenticated client: %w", err)
	}

	logger.Info().Msg("openstack client created")

	return &Client{
		Logger:   logger,
		Client:   client,
		Services: map[ServiceType]*gophercloud.ServiceClient{},
		Specs:    &pluginSpec,
	}, nil
}

func (c *Client) GetServiceClient(key ServiceType) (*gophercloud.ServiceClient, error) {

	if service, ok := c.Services[key]; ok && service != nil {
		c.Logger.Info().Str("type", string(key)).Msg("returning existing service client")
		return service, nil
	}

	// no existing service client, need to initialise one

	if _, ok := serviceConfigMap[key]; !ok {
		c.Logger.Error().Str("type", string(key)).Msg("invalid service client type")
		panic(fmt.Sprintf("invalid service type: %q", key))
	}

	c.Logger.Info().Str("type", string(key)).Msg("creating new service client")

	region := ""
	if c.Specs.Region != nil {
		region = *c.Specs.Region
	}

	client, err := serviceConfigMap[key].newClient(c.Client, gophercloud.EndpointOpts{Region: region})

	if err != nil {
		c.Logger.Error().Str("type", string(key)).Err(err).Msg("error creating service client")
		return nil, err
	}
	client.Microversion = serviceConfigMap[key].getMicroversion(c.Specs)

	// save to object
	c.Services[key] = client

	c.Logger.Info().Str("type", string(key)).Msg("new service client ready")

	return client, nil
}

const (
	// defaults currently referring to Train
	DefaultComputeV2Microversion      = "2.79"
	DefaultIdentityV3Microversion     = "3.13"
	DefaultBlockStorageV3Microversion = "3.59"
	DefaultImageServiceV2Microversion = "2.9"
)

type ServiceType string

const (
	// IdentityV3 identifies the OpenStack Identity V3 service (Keystone).
	IdentityV3 ServiceType = "openstack_identity_v3"
	// Compute identifies the penStack Compute V2 service (Nova).
	ComputeV2 = "openstack_compute_v2"
	// NetworkV2 identifies the OpenStack Network V2 service (Neutron).
	NetworkV2 = "openstack_network_v2"
	// BlockStorageV3 identifies the OpenStack Block Storage V3 service (Cinder).
	BlockStorageV3 = "openstack_blockstorage_v3"
	// ImageServiceV2 identifies the OpenStack Image Service V2 service (Glance).
	ImageServiceV2 = "openstack_imageservice_v2"
)

type serviceConfig struct {
	newClient       func(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)
	getMicroversion func(spec *Spec) string
}

var serviceConfigMap = map[ServiceType]serviceConfig{
	IdentityV3: {
		newClient: openstack.NewIdentityV3,
		getMicroversion: func(spec *Spec) string {
			microversion := DefaultIdentityV3Microversion
			if spec.IdentityV3Microversion != nil {
				microversion = *spec.IdentityV3Microversion
			}
			return microversion
		},
	},
	ComputeV2: {
		newClient: openstack.NewComputeV2,
		getMicroversion: func(spec *Spec) string {
			microversion := DefaultComputeV2Microversion
			if spec.ComputeV2Microversion != nil {
				microversion = *spec.ComputeV2Microversion
			}
			return microversion
		},
	},
	NetworkV2: {
		newClient: openstack.NewNetworkV2,
		getMicroversion: func(spec *Spec) string {
			// TODO: check if we need to leverage/support micro-versions
			return ""
		},
	},
	BlockStorageV3: {
		newClient: openstack.NewBlockStorageV3,
		getMicroversion: func(spec *Spec) string {
			microversion := DefaultBlockStorageV3Microversion
			if spec.BlockStorageV3Microversion != nil {
				microversion = *spec.BlockStorageV3Microversion
			}
			return microversion
		},
	},
	ImageServiceV2: {
		newClient: openstack.NewImageServiceV2,
		getMicroversion: func(spec *Spec) string {
			microversion := DefaultImageServiceV2Microversion
			if spec.ImageServiceV2Microversion != nil {
				microversion = *spec.ImageServiceV2Microversion
			}
			return microversion
		},
	},
}

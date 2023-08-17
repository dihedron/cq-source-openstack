package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/rs/zerolog"
)

type Client struct {
	Client *gophercloud.ProviderClient
	Spec   Spec

	logger zerolog.Logger
	// tables   schema.Tables
	mutex    sync.RWMutex
	services map[ServiceType]*gophercloud.ServiceClient
}

func (c *Client) ID() string {
	return "github.com/dihedron/cq-source-openstack"
}

func (c *Client) Logger() *zerolog.Logger {
	return &c.logger
}

func New(ctx context.Context, logger zerolog.Logger, spec *Spec) (*Client, error) {

	logger.Debug().Str("spec", format.ToJSON(spec)).Msg("plugin configuration")

	auth := gophercloud.AuthOptions{
		AllowReauth: true,
	}

	if spec.EndpointUrl != nil {
		auth.IdentityEndpoint = *spec.EndpointUrl
	}
	if spec.UserID != nil {
		auth.UserID = *spec.UserID
	}
	if spec.Username != nil {
		auth.Username = *spec.Username
	}
	if spec.Password != nil {
		auth.Password = *spec.Password
	}
	if spec.ProjectID != nil {
		auth.TenantID = *spec.ProjectID
	}
	if spec.ProjectName != nil {
		auth.TenantName = *spec.ProjectName
	}
	if spec.DomainID != nil {
		auth.DomainID = *spec.DomainID
	}
	if spec.DomainName != nil {
		auth.DomainName = *spec.DomainName
	}
	if spec.AccessToken != nil {
		auth.TokenID = *spec.AccessToken
	}
	if spec.AppCredentialID != nil {
		auth.ApplicationCredentialID = *spec.AppCredentialID
	}
	if spec.AppCredentialSecret != nil {
		auth.ApplicationCredentialSecret = *spec.AppCredentialSecret
	}
	if spec.AllowReauth != nil {
		auth.AllowReauth = *spec.AllowReauth
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
		Spec:     *spec,
		logger:   logger,
		Client:   client,
		services: map[ServiceType]*gophercloud.ServiceClient{},
	}, nil
}

func (c *Client) GetServiceClient(key ServiceType) (*gophercloud.ServiceClient, error) {

	c.mutex.RLock()
	if service, ok := c.services[key]; ok && service != nil {
		c.Logger().Info().Str("type", string(key)).Msg("returning existing service client")
		c.mutex.RUnlock()
		return service, nil
	}
	c.mutex.RUnlock()
	return c.initServiceClient(key)
}

func (c *Client) initServiceClient(key ServiceType) (*gophercloud.ServiceClient, error) {

	// no existing service client, need to initialise one
	if _, ok := serviceConfigMap[key]; !ok {
		c.Logger().Error().Str("type", string(key)).Msg("invalid service client type")
		panic(fmt.Sprintf("invalid service type: %q", key))
	}

	c.Logger().Info().Str("type", string(key)).Msg("creating new service client")

	region := ""
	if c.Spec.Region != nil {
		region = *c.Spec.Region
	}

	client, err := serviceConfigMap[key].newClient(c.Client, gophercloud.EndpointOpts{Region: region})

	if err != nil {
		c.Logger().Error().Str("type", string(key)).Err(err).Msg("error creating service client")
		return nil, err
	}
	client.Microversion = serviceConfigMap[key].getMicroversion(&c.Spec)

	// save to object
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.services[key] = client

	c.Logger().Info().Str("type", string(key)).Msg("new service client ready")

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

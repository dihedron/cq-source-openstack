package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"

	// "github.com/dihedron/cq-source-openstack/resources/services/ironic"
	"github.com/dihedron/cq-source-openstack/resources/services/cinder"
	"github.com/dihedron/cq-source-openstack/resources/services/glance"
	"github.com/dihedron/cq-source-openstack/resources/services/keystone"
	"github.com/dihedron/cq-source-openstack/resources/services/neutron"
	"github.com/dihedron/cq-source-openstack/resources/services/nova"
	"github.com/rs/zerolog"
)

type Client struct {
	logger     zerolog.Logger
	config     client.Spec
	tables     schema.Tables
	syncClient *client.Client
	scheduler  *scheduler.Scheduler

	plugin.UnimplementedDestination
}

func Configure(ctx context.Context, logger zerolog.Logger, spec []byte, opts plugin.NewClientOptions) (plugin.Client, error) {
	if opts.NoConnection {
		return &Client{
			logger: logger,
			tables: getTables(),
		}, nil
	}

	config := &client.Spec{}
	if err := json.Unmarshal(spec, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}

	syncClient, err := client.New(ctx, logger, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &Client{
		logger:     logger,
		config:     *config,
		tables:     getTables(),
		syncClient: syncClient,
		scheduler:  scheduler.NewScheduler(scheduler.WithLogger(logger)),
	}, nil
}

func (c *Client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return err
	}

	return c.scheduler.Sync(ctx, c.syncClient, tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
}

func (c *Client) Tables(_ context.Context, options plugin.TableOptions) (schema.Tables, error) {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

func (*Client) Close(_ context.Context) error {
	// TODO: Add your client cleanup here
	return nil
}

func getTables() schema.Tables {
	tables := schema.Tables{
		// ironic.Allocations(),
		// ironic.Drivers(),
		// ironic.Nodes(),
		// ironic.Ports(),
		cinder.Attachments(),
		cinder.AvailabilityZones(),
		cinder.Limits(),
		cinder.QoS(),
		cinder.QuotaSets(),
		cinder.QuotaSetsUsage(),
		cinder.Services(),
		cinder.Snapshots(),
		cinder.Volumes(),
		nova.Aggregates(),
		nova.Flavors(),
		nova.Hypervisors(),
		// nova.Images(),
		nova.Instances(),
		// nova.Networks(),
		// nova.SecGroups(),
		nova.ServerUsage(),
		// nova.TenantNetworks(),
		keystone.Domains(),
		keystone.Projects(),
		// keystone.Tenants(),
		keystone.Users(),
		keystone.Services(),
		glance.Images(),
		neutron.Networks(),
		neutron.Ports(),
		neutron.SecurityGroups(),
		neutron.SecurityGroupRules(),
	}
	if err := transformers.TransformTables(tables); err != nil {
		panic(err)
	}
	for _, t := range tables {
		schema.AddCqIDs(t)
	}
	return tables
}

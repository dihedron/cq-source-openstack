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
	"github.com/dihedron/cq-source-openstack/resources/internal/pattern_matcher"

	"github.com/dihedron/cq-source-openstack/resources/services/ironic"
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
	config := &client.Spec{}
	if err := json.Unmarshal(spec, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}

	if opts.NoConnection {
		return &Client{
			logger: logger,
			tables: getTables(config),
		}, nil
	}

	syncClient, err := client.New(ctx, logger, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	tables := getTables(config)

	return &Client{
		logger:     logger,
		config:     *config,
		tables:     tables,
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

func getTables(spec *client.Spec) schema.Tables {
	available_tables := schema.Tables{
		ironic.Allocations(),
		ironic.Drivers(),
		ironic.Nodes(),
		ironic.Ports(),
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
		nova.Images(),
		nova.Instances(),
		nova.Networks(),
		nova.SecGroups(),
		nova.ServerUsage(),
		nova.TenantNetworks(),
		keystone.Domains(),
		keystone.Projects(),
		keystone.Tenants(),
		keystone.Users(),
		keystone.Services(),
		glance.Images(),
		neutron.Networks(),
		neutron.Ports(),
		neutron.SecurityGroups(),
		neutron.SecurityGroupRules(),
	}

	// must compile these patterns to be included
	includesDaSpec := spec.IncludedTables
	// must compile these patterns to be excluded
	excludesDaSpec := spec.ExcludedTables

	// if includesDaSpec is empty, include everything
	var pm *pattern_matcher.PatternMatcher
	if len(includesDaSpec) == 0 {
		pm = pattern_matcher.New(
			pattern_matcher.WithExclude(excludesDaSpec),
		)
	} else if len(excludesDaSpec) == 0 {
		pm = pattern_matcher.New(
			pattern_matcher.WithInclude(includesDaSpec),
		)
	} else {
		pm = pattern_matcher.New(
			pattern_matcher.WithInclude(includesDaSpec),
			pattern_matcher.WithExclude(excludesDaSpec),
		)
	}

	tables := make(schema.Tables, 0, len(available_tables))
	for _, t := range available_tables {
		if pm.Match(t.Name) {
			tables = append(tables, t)
		}
	}

	if err := transformers.TransformTables(tables); err != nil {
		panic(err)
	}
	for _, t := range tables {
		schema.AddCqIDs(t)
	}
	return tables
}

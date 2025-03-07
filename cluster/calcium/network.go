package calcium

import (
	"context"

	enginetypes "github.com/projecteru2/core/engine/types"
	"github.com/projecteru2/core/log"
	"github.com/projecteru2/core/types"

	"github.com/pkg/errors"
)

// ListNetworks by podname
// get one node from a pod
// and list networks
// only get those driven by network driver
func (c *Calcium) ListNetworks(ctx context.Context, podname string, driver string) ([]*enginetypes.Network, error) {
	logger := log.WithField("Calcium", "ListNetworks").WithField("podname", podname).WithField("driver", driver)
	networks := []*enginetypes.Network{}
	nodes, err := c.store.GetNodesByPod(ctx, &types.NodeFilter{Podname: podname})
	if err != nil {
		return networks, logger.ErrWithTracing(ctx, errors.WithStack(err))
	}

	if len(nodes) == 0 {
		return networks, logger.ErrWithTracing(ctx, errors.WithStack(types.NewDetailedErr(types.ErrPodNoNodes, podname)))
	}

	drivers := []string{}
	if driver != "" {
		drivers = append(drivers, driver)
	}

	node := nodes[0]

	networks, err = node.Engine.NetworkList(ctx, drivers)
	return networks, logger.ErrWithTracing(ctx, errors.WithStack(err))
}

// ConnectNetwork connect to a network
func (c *Calcium) ConnectNetwork(ctx context.Context, network, target, ipv4, ipv6 string) ([]string, error) {
	logger := log.WithField("Calcium", "ConnectNetwork").WithField("network", network).WithField("target", target).WithField("ipv4", ipv4).WithField("ipv6", ipv6)
	workload, err := c.GetWorkload(ctx, target)
	if err != nil {
		return nil, logger.ErrWithTracing(ctx, errors.WithStack(err))
	}

	networks, err := workload.Engine.NetworkConnect(ctx, network, target, ipv4, ipv6)
	return networks, logger.ErrWithTracing(ctx, errors.WithStack(err))
}

// DisconnectNetwork connect to a network
func (c *Calcium) DisconnectNetwork(ctx context.Context, network, target string, force bool) error {
	logger := log.WithField("Calcium", "DisconnectNetwork").WithField("network", network).WithField("target", target).WithField("force", force)
	workload, err := c.GetWorkload(ctx, target)
	if err != nil {
		return logger.ErrWithTracing(ctx, errors.WithStack(err))
	}

	return logger.ErrWithTracing(ctx, errors.WithStack(workload.Engine.NetworkDisconnect(ctx, network, target, force)))
}

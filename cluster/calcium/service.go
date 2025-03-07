package calcium

import (
	"context"
	"sync"
	"time"

	"github.com/projecteru2/core/log"
	"github.com/projecteru2/core/types"
	"github.com/projecteru2/core/utils"

	"github.com/pkg/errors"
)

// WatchServiceStatus returns chan of available service address
func (c *Calcium) WatchServiceStatus(ctx context.Context) (<-chan types.ServiceStatus, error) {
	id, ch := c.watcher.Subscribe(ctx)
	_ = c.pool.Invoke(func() {
		<-ctx.Done()
		c.watcher.Unsubscribe(id)
	})
	return ch, nil
}

// RegisterService writes self service address in store
func (c *Calcium) RegisterService(ctx context.Context) (unregister func(), err error) {
	serviceAddress, err := utils.GetOutboundAddress(c.config.Bind)
	if err != nil {
		log.Errorf(ctx, "[RegisterService] failed to get outbound address: %v", err)
		return
	}

	var (
		expiry            <-chan struct{}
		unregisterService func()
	)
	for {
		if expiry, unregisterService, err = c.registerService(ctx, serviceAddress); err == nil {
			break
		}
		if errors.Is(err, types.ErrKeyExists) {
			log.Debugf(ctx, "[RegisterService] service key exists: %v", err)
			time.Sleep(time.Second)
			continue
		}
		log.Errorf(ctx, "[RegisterService] failed to first register service: %+v", err)
		return nil, errors.WithStack(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(ctx)
	_ = c.pool.Invoke(func() {
		defer func() {
			unregisterService()
			wg.Done()
		}()

		for {
			select {
			case <-expiry:
				// The original one had been expired, we're going to register again.
				if ne, us, err := c.registerService(ctx, serviceAddress); err != nil {
					log.Errorf(ctx, "[RegisterService] failed to re-register service: %v", err)
					time.Sleep(c.config.GRPCConfig.ServiceHeartbeatInterval)
				} else {
					expiry = ne
					unregisterService = us
				}

			case <-ctx.Done():
				log.Infof(ctx, "[RegisterService] heartbeat done: %v", ctx.Err())
				return
			}
		}
	})
	return func() {
		cancel()
		wg.Wait()
	}, nil
}

func (c *Calcium) registerService(ctx context.Context, addr string) (<-chan struct{}, func(), error) {
	return c.store.RegisterService(ctx, addr, c.config.GRPCConfig.ServiceHeartbeatInterval)
}

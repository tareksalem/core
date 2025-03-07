package models

import (
	"github.com/projecteru2/core/store/etcdv3/meta"
	coretypes "github.com/projecteru2/core/types"

	"github.com/sirupsen/logrus"
)

// CPUMem manages cpu and memory
type CPUMem struct {
	Config coretypes.Config
	store  meta.KV
}

func NewCPUMem(config coretypes.Config) (*CPUMem, error) {
	c := &CPUMem{Config: config}
	var err error
	if len(config.Etcd.Machines) > 0 {
		c.store, err = meta.NewETCD(config.Etcd, nil)
		if err != nil {
			logrus.Errorf("[NewCPUMem] failed to create etcd client, err: %v", err)
			return nil, err
		}
	}
	return c, nil
}

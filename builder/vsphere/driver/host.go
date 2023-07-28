// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package driver

import (
	"fmt"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type Host struct {
	driver *VCenterDriver
	host   *object.HostSystem
}

func (d *VCenterDriver) NewHost(ref *types.ManagedObjectReference) *Host {
	return &Host{
		host:   object.NewHostSystem(d.client.Client, *ref),
		driver: d,
	}
}

func (d *VCenterDriver) FindHost(name string) (*Host, error) {
	h, err := d.finder.HostSystem(d.ctx, name)
	if err != nil {
		return nil, err
	}
	return &Host{
		host:   h,
		driver: d,
	}, nil
}

func (h *Host) Info(params ...string) (*mo.HostSystem, error) {
	var p []string
	if len(params) == 0 {
		p = []string{"*"}
	} else {
		p = params
	}
	var info mo.HostSystem
	err := h.host.Properties(h.driver.ctx, h.host.Reference(), p, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// ParentCluster returns the parent cluster for the host. If the host is not part of a cluster,
// then the return value will be nil.
func (h *Host) ParentCluster() (*mo.ClusterComputeResource, error) {
	hostInfo, err := h.Info()
	if err != nil {
		return nil, err
	}

	switch hostInfo.Parent.Type {
	case "ComputeResource":
		// not a cluster
		return nil, nil
	case "ClusterComputeResource":
		cluster := object.NewClusterComputeResource(h.driver.vimClient, hostInfo.Parent.Reference())
		var clusterMo mo.ClusterComputeResource
		err = property.DefaultCollector(h.driver.vimClient).RetrieveOne(h.driver.ctx, cluster.Reference(), nil, &clusterMo)
		if err != nil {
			return nil, err
		}
		return &clusterMo, nil
	default:
		return nil, fmt.Errorf("unexpected parent type for host: %s", hostInfo.Parent.Type)
	}
}

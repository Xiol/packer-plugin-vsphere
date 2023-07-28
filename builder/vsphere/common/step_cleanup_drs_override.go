// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-vsphere/builder/vsphere/driver"
)

type StepCleanupDRSOverride struct{}

func (s *StepCleanupDRSOverride) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	if _, ok := state.GetOk("drs_override"); !ok {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packersdk.Ui)
	vm := state.Get("vm").(*driver.VirtualMachineDriver)
	cluster := state.Get("cluster").(string)
	driver := state.Get("driver").(driver.Driver)

	ui.Say("Cleaning up DRS override...")
	vmMo, err := vm.Info()
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	err = driver.CleanupDRSRule(cluster, vmMo.Reference())
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepCleanupDRSOverride) Cleanup(state multistep.StateBag) {}

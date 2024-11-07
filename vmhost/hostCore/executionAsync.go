package hostCore

import (
	"github.com/kalyan3104/k-chain-vm-go/vmhost"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	runtime := host.Runtime()
	async := host.Async()
	runtime.SetRuntimeBreakpointValue(vmhost.BreakpointNone)

	legacyGroupID := vmhost.LegacyAsyncCallGroupID
	legacyGroup, exists := async.GetCallGroup(legacyGroupID)
	if !exists {
		return vmhost.ErrLegacyAsyncCallNotFound

	}

	if legacyGroup.IsComplete() {
		return vmhost.ErrLegacyAsyncCallInvalid
	}

	return nil
}

package wasmer2

import (
	"unsafe"

	"github.com/kalyan3104/k-chain-vm-go/executor"
)

func getVMHooksFromContextRawPtr(contextRawPtr unsafe.Pointer) executor.VMHooks {
	vmHooksPtrPtr := (*uintptr)(contextRawPtr)
	vmHooksPtr := *vmHooksPtrPtr
	return *(*executor.VMHooks)(unsafe.Pointer(vmHooksPtr))
}

func funcPointer(cFuncPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(cFuncPtr)
}

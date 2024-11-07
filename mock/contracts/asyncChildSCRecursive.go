package contracts

import (
	"math/big"

	"github.com/kalyan3104/k-chain-vm-common-go/txDataBuilder"
	mock "github.com/kalyan3104/k-chain-vm-go/mock/context"
	test "github.com/kalyan3104/k-chain-vm-go/testcommon"
	"github.com/kalyan3104/k-chain-vm-go/vmhost"
)

// RecursiveAsyncCallRecursiveChildMock is an exposed mock contract method
func RecursiveAsyncCallRecursiveChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("recursiveAsyncCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		arguments := host.Runtime().Arguments()

		err := host.Metering().UseGasBounded(testConfig.GasUsedByChild)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

		var recursiveChildCalls uint64
		if len(arguments) > 0 {
			recursiveChildCalls = big.NewInt(0).SetBytes(arguments[0]).Uint64()
		} else {
			recursiveChildCalls = 1
		}
		recursiveChildCalls = recursiveChildCalls - 1
		returnValue := big.NewInt(int64(recursiveChildCalls)).Bytes()
		if len(arguments) == 2 {
			returnValue = arguments[1]
		}
		host.Output().Finish(returnValue)
		if recursiveChildCalls == 0 {
			return instance
		}

		destination := host.Runtime().GetContextAddress()
		function := "recursiveAsyncCall"
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.BigInt(big.NewInt(int64(recursiveChildCalls)))

		async := host.Async()
		err = async.RegisterLegacyAsyncCall(destination, callData.ToBytes(), value)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointExecutionFailed)
			return instance
		}

		return instance
	})
}

// CallBackRecursiveChildMock is an exposed mock contract method
func CallBackRecursiveChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

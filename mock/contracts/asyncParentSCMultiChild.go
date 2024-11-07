package contracts

import (
	"math/big"

	"github.com/kalyan3104/k-chain-vm-common-go/txDataBuilder"
	mock "github.com/kalyan3104/k-chain-vm-go/mock/context"
	test "github.com/kalyan3104/k-chain-vm-go/testcommon"
	"github.com/kalyan3104/k-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallMultiChildMock is an exposed mock contract method
func ForwardAsyncCallMultiChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

		for childCall := 0; childCall < testConfig.ChildCalls; childCall++ {
			callData := txDataBuilder.NewBuilder()
			callData.Func(function)
			// recursiveChildCalls
			callData.BigInt(big.NewInt(1))
			// child will return this
			callData.BigInt(big.NewInt(int64(childCall)))

			async := host.Async()
			err := async.RegisterAsyncCall("myAsyncGroup", &vmhost.AsyncCall{
				Status:          vmhost.AsyncCallPending,
				Destination:     destination,
				Data:            callData.ToBytes(),
				ValueBytes:      value,
				GasLimit:        uint64(300),
				SuccessCallback: "callBack",
				ErrorCallback:   "callBack",
			})
			require.Nil(t, err)
		}

		return instance

	})
}

// CallBackMultiChildMock is an exposed mock contract method
func CallBackMultiChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

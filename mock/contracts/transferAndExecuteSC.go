package contracts

import (
	"fmt"
	"math/big"

	mock "github.com/kalyan3104/k-chain-vm-go/mock/context"
	"github.com/kalyan3104/k-chain-vm-go/testcommon"
	test "github.com/kalyan3104/k-chain-vm-go/testcommon"
	"github.com/kalyan3104/k-chain-vm-go/vmhost/vmhooks"
)

// TransferAndExecuteFuncName -
var TransferAndExecuteFuncName = "transferAndExecute"

// TransferAndExecuteReturnData -
var TransferAndExecuteReturnData = []byte{1, 2, 3}

// TransferAndExecute is an exposed mock contract method
func TransferAndExecute(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod(TransferAndExecuteFuncName, func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		_ = host.Metering().UseGasBounded(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		noOfTransfers := int(big.NewInt(0).SetBytes(arguments[0]).Int64())

		for transfer := 0; transfer < noOfTransfers; transfer++ {
			vmhooks.TransferValueExecuteWithTypedArgs(host,
				GetChildAddressForTransfer(transfer),
				big.NewInt(testConfig.TransferFromParentToChild),
				int64(testConfig.GasProvidedToChild),
				big.NewInt(int64(transfer)).Bytes(), // transfer data
				[][]byte{},
			)
		}

		host.Output().Finish(TransferAndExecuteReturnData)

		return instance
	})
}

func TransferREWAToParent(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("transferREWAToParent", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		_ = host.Metering().UseGasBounded(testConfig.GasUsedByChild)

		vmhooks.TransferValueExecuteWithTypedArgs(host,
			test.ParentAddress,
			big.NewInt(testConfig.ChildBalance/2),
			0,
			[]byte{}, // transfer data
			[][]byte{},
		)

		return instance
	})
}

// GetChildAddressForTransfer -
func GetChildAddressForTransfer(transfer int) []byte {
	return testcommon.MakeTestSCAddress(fmt.Sprintf("childSC-%d", transfer))
}

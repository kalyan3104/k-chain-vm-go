package contracts

import (
	mock "github.com/kalyan3104/k-chain-vm-go/mock/context"
	test "github.com/kalyan3104/k-chain-vm-go/testcommon"
)

// ChildAsyncMultiGroupsMock is an exposed mock contract method
func ChildAsyncMultiGroupsMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	for _, groupConfig := range AsyncGroupsConfig {
		for g := 1; g < len(groupConfig); g++ {
			functionName := groupConfig[g]
			instanceMock.AddMockMethod(functionName,
				test.WasteGasWithReturnDataMockMethod(
					instanceMock,
					testConfig.GasUsedByChild,
					[]byte(functionName+test.TestReturnDataSuffix)))
		}
	}
}

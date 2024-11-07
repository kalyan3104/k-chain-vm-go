package vmhost

import (
	"math/big"

	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
)

// MakeEmptyVMOutput creates a vmcommon.VMOutput struct with default values
func MakeEmptyVMOutput() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		ReturnData:      make([][]byte, 0),
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
	}
}

// AddFinishData appends the provided []byte to the ReturnData of the given vmOutput
func AddFinishData(vmOutput *vmcommon.VMOutput, data []byte) {
	vmOutput.ReturnData = append(vmOutput.ReturnData, data)
}

// AddNewOutputAccount creates a new vmcommon.OutputAccount from the provided arguments and adds it to OutputAccounts of the provided vmOutput
func AddNewOutputAccount(vmOutput *vmcommon.VMOutput, address []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
	return AddNewOutputAccountWithSender(vmOutput, address, nil, balanceDelta, data)
}

// AddNewOutputAccountWithSender creates a new vmcommon.OutputAccount from the provided arguments and adds it to OutputAccounts of the provided vmOutput
func AddNewOutputAccountWithSender(vmOutput *vmcommon.VMOutput, address []byte, senderAddress []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
	account := &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(balanceDelta),
		Balance:        nil,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		Code:           nil,
	}
	if data != nil {
		account.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Index:         1,
				Data:          data,
				Value:         big.NewInt(balanceDelta),
				SenderAddress: senderAddress,
			},
		}
	}
	vmOutput.OutputAccounts[string(address)] = account
	return account
}

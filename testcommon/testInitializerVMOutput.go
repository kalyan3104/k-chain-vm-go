package testcommon

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

// AddNewOutputTransfer creates a new vmcommon.OutputAccount from the provided arguments and adds it to OutputAccounts of the provided vmOutput
func AddNewOutputTransfer(vmOutput *vmcommon.VMOutput, index uint32, sender []byte, address []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
	account := vmOutput.OutputAccounts[string(address)]
	if account == nil {
		account = &vmcommon.OutputAccount{
			Address:        address,
			Nonce:          0,
			BalanceDelta:   big.NewInt(balanceDelta),
			Balance:        nil,
			StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
			Code:           nil,
		}
	}
	if data != nil {
		if account.OutputTransfers == nil {
			account.OutputTransfers = make([]vmcommon.OutputTransfer, 0)
		}
		account.OutputTransfers = append(account.OutputTransfers, vmcommon.OutputTransfer{
			Index:         index,
			Data:          data,
			Value:         big.NewInt(balanceDelta),
			SenderAddress: sender,
		},
		)
	}
	vmOutput.OutputAccounts[string(address)] = account
	return account
}

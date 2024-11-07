package vmhost

import (
	"errors"
	"math/big"
	"strings"

	"github.com/kalyan3104/k-chain-core-go/data/vm"
	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
)

// DefaultVMType is an exposed value to use in tests
var DefaultVMType = []byte{0xF, 0xF}

// ErrAccountNotFound is an exposed value to use in tests
var ErrAccountNotFound = errors.New("account not found")

// UserAddress is an exposed value to use in tests
var UserAddress = []byte("userAccount.....................")

// AddressSize is the size of an account address, in bytes.
const AddressSize = 32

// SCAddressPrefix is the prefix of any smart contract address used for testing.
var SCAddressPrefix = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0f")

// ParentAddress is an exposed value to use in tests
var ParentAddress = MakeTestSCAddress("parentSC")

// MakeTestSCAddress generates a new smart contract address to be used for
// testing based on the given identifier.
func MakeTestSCAddress(identifier string) []byte {
	return makeTestAddress(SCAddressPrefix, identifier)
}

func makeTestAddress(_ []byte, identifier string) []byte {
	numberOfTrailingDots := AddressSize - len(SCAddressPrefix) - len(identifier)
	leftBytes := SCAddressPrefix
	rightBytes := []byte(identifier + strings.Repeat(".", numberOfTrailingDots))
	return append(leftBytes, rightBytes...)
}

// MakeEmptyContractCallInput instantiates an empty ContractCallInput
func MakeEmptyContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:           nil,
			Arguments:            make([][]byte, 0),
			CallValue:            big.NewInt(0),
			CallType:             vm.DirectCall,
			GasPrice:             1,
			GasProvided:          0,
			ReturnCallAfterError: false,
		},
		RecipientAddr: nil,
		Function:      "",
	}
}

// MakeContractCallInput creates a ContractCallInput and sets the provided arguments
func MakeContractCallInput(
	caller []byte,
	recipient []byte,
	function string,
	value int,
) *vmcommon.ContractCallInput {
	input := MakeEmptyContractCallInput()
	SetCallParties(input, caller, recipient)
	input.Function = function
	input.CallValue = big.NewInt(int64(value))
	return input
}

// SetCallParties sets the caller and recipient of the given ContractCallInput
func SetCallParties(input *vmcommon.ContractCallInput, caller []byte, recipient []byte) {
	input.OriginalCallerAddr = []byte("address_original_caller")
	input.CallerAddr = caller
	input.RecipientAddr = recipient
}

// AddArgument adds the provided argument to the ContractCallInput
func AddArgument(input *vmcommon.ContractCallInput, argument []byte) {
	if input.Arguments == nil {
		input.Arguments = make([][]byte, 0)
	}
	input.Arguments = append(input.Arguments, argument)
}

// CopyTxHashes copies the tx hashes from a source ContractCallInput into another
func CopyTxHashes(input *vmcommon.ContractCallInput, sourceInput *vmcommon.ContractCallInput) {
	input.CurrentTxHash = sourceInput.CurrentTxHash
	input.PrevTxHash = sourceInput.PrevTxHash
	input.OriginalTxHash = sourceInput.OriginalTxHash
}

// DefaultTestContractCallInput creates a vmcommon.ContractCallInput struct
// with default values.
func DefaultTestContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  UserAddress,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vm.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: ParentAddress,
		Function:      "function",
	}
}

// ContractCallInputBuilder extends a ContractCallInput for extra building functionality during testing
type ContractCallInputBuilder struct {
	vmcommon.ContractCallInput
}

// CreateTestContractCallInputBuilder is a builder for ContractCallInputBuilder
func CreateTestContractCallInputBuilder() *ContractCallInputBuilder {
	return &ContractCallInputBuilder{
		ContractCallInput: *DefaultTestContractCallInput(),
	}
}

// WithRecipientAddr provides the recepient address of ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithRecipientAddr(address []byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.RecipientAddr = address
	return contractInput
}

// WithCallerAddr provides the caller address of ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithCallerAddr(address []byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.CallerAddr = address
	return contractInput
}

// WithGasProvided provides the gas of ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithGasProvided(gas uint64) *ContractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.GasProvided = gas
	return contractInput
}

// WithFunction provides the function to be called for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithFunction(function string) *ContractCallInputBuilder {
	contractInput.ContractCallInput.Function = function
	return contractInput
}

// WithArguments provides the arguments to be called for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithArguments(arguments ...[]byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.Arguments = arguments
	return contractInput
}

// WithCallType provides the arguments to be called for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithCallType(callType vm.CallType) *ContractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.CallType = callType
	return contractInput
}

// WithCurrentTxHash provides the CurrentTxHash for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithCurrentTxHash(txHash []byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.CurrentTxHash = txHash
	return contractInput
}

func (contractInput *ContractCallInputBuilder) initDCDTTransferIfNeeded() {
	if len(contractInput.DCDTTransfers) == 0 {
		contractInput.DCDTTransfers = make([]*vmcommon.DCDTTransfer, 1)
		contractInput.DCDTTransfers[0] = &vmcommon.DCDTTransfer{}
	}
}

// WithDCDTValue provides the DCDTValue for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithDCDTValue(dcdtValue *big.Int) *ContractCallInputBuilder {
	contractInput.initDCDTTransferIfNeeded()
	contractInput.ContractCallInput.DCDTTransfers[0].DCDTValue = dcdtValue
	return contractInput
}

// WithDCDTTokenName provides the DCDTTokenName for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithDCDTTokenName(dcdtTokenName []byte) *ContractCallInputBuilder {
	contractInput.initDCDTTransferIfNeeded()
	contractInput.ContractCallInput.DCDTTransfers[0].DCDTTokenName = dcdtTokenName
	return contractInput
}

// Build completes the build of a ContractCallInput
func (contractInput *ContractCallInputBuilder) Build() *vmcommon.ContractCallInput {
	return &contractInput.ContractCallInput
}

// ContractCreateInputBuilder extends a ContractCreateInput for extra building functionality during testing
type ContractCreateInputBuilder struct {
	vmcommon.ContractCreateInput
}

// WithGasProvided provides the GasProvided for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithGasProvided(gas uint64) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.GasProvided = gas
	return contractInput
}

// WithContractCode provides the ContractCode for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithContractCode(code []byte) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.ContractCode = code
	return contractInput
}

// WithCallerAddr provides the CallerAddr for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithCallerAddr(address []byte) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.CallerAddr = address
	return contractInput
}

// WithCallValue provides the CallValue for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithCallValue(callValue int64) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.CallValue = big.NewInt(callValue)
	return contractInput
}

// WithArguments provides the Arguments for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithArguments(arguments ...[]byte) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.Arguments = arguments
	return contractInput
}

// Build completes the build of a ContractCreateInput
func (contractInput *ContractCreateInputBuilder) Build() *vmcommon.ContractCreateInput {
	return &contractInput.ContractCreateInput
}

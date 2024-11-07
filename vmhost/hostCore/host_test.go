package hostCore

import (
	"math"
	"testing"

	"github.com/kalyan3104/k-chain-core-go/core"
	"github.com/kalyan3104/k-chain-scenario-go/worldmock"
	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
	"github.com/kalyan3104/k-chain-vm-common-go/builtInFunctions"
	"github.com/kalyan3104/k-chain-vm-common-go/parsers"
	"github.com/kalyan3104/k-chain-vm-go/vmhost"
	"github.com/kalyan3104/k-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

func TestNewVMHost(t *testing.T) {
	blockchainHook := worldmock.NewMockWorld()
	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	epochNotifier := &mock.EpochNotifierStub{}
	epochsHandler := &worldmock.EnableEpochsHandlerStub{}
	vmType := []byte("vmType")
	dcdtTransferParser, err := parsers.NewDCDTTransferParser(worldmock.WorldMarshalizer)
	require.Nil(t, err)

	makeHostParameters := func() *vmhost.VMHostParameters {
		return &vmhost.VMHostParameters{
			VMType:               vmType,
			DCDTTransferParser:   dcdtTransferParser,
			BuiltInFuncContainer: bfc,
			EpochNotifier:        epochNotifier,
			EnableEpochsHandler:  epochsHandler,
			Hasher:               worldmock.DefaultHasher,
		}
	}

	t.Run("NilBlockchainHook", func(t *testing.T) {
		host, err := NewVMHost(nil, makeHostParameters())
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilBlockChainHook)
	})
	t.Run("NilHostParameters", func(t *testing.T) {
		host, err := NewVMHost(blockchainHook, nil)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilHostParameters)
	})
	t.Run("NilDCDTTransferParser", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.DCDTTransferParser = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilDCDTTransferParser)
	})
	t.Run("NilBuiltInFunctionsContainer", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.BuiltInFuncContainer = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilBuiltInFunctionsContainer)
	})
	t.Run("NilEpochNotifier", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EpochNotifier = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilEpochNotifier)
	})
	t.Run("NilEnableEpochsHandler", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EnableEpochsHandler = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilEnableEpochsHandler)
	})
	t.Run("InvalidEnableEpochsHandler", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EnableEpochsHandler = &worldmock.EnableEpochsHandlerStub{
			IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
				return false
			},
		}
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, core.ErrInvalidEnableEpochsHandler)
	})
	t.Run("NilHasher", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.Hasher = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilHasher)
	})
	t.Run("NilVMType", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.VMType = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilVMType)
	})
}

func TestValidateVMInput(t *testing.T) {
	vmInput := &vmcommon.VMInput{
		GasProvided: 0,
	}

	vmInput.GasProvided = math.MaxUint64
	err := validateVMInput(vmInput)
	require.ErrorIs(t, err, vmhost.ErrInvalidGasProvided)

	vmInput.GasProvided = math.MaxInt64
	err = validateVMInput(vmInput)
	require.Nil(t, err)
}

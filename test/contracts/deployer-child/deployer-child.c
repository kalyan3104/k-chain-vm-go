#include "../kvm/context.h"
#include "../kvm/test_utils.h"

byte parentAddress[32] = {};

void init() {
	getArgument(0, parentAddress);	

	int isParentContract = isSmartContract(parentAddress);
	if (isParentContract == 0) {
		byte message[] = "[from child] parent not a contract";
		signalError(message, sizeof(message) - 1);
	}
}

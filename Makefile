.PHONY: test test-short build vmserver clean

VM_VERSION := $(shell git describe --tags --long --dirty --always)

clean:
	go clean -cache -testcache

build:
	go build ./...

vmserver:
ifndef VMSERVER_PATH
	$(error VMSERVER_PATH is undefined)
endif
	go build -o ./cmd/vmserver/vmserver ./cmd/vmserver
	cp ./cmd/vmserver/vmserver ${VMSERVER_PATH}

test:
	go clean -cache -testcache
	VMEXECUTOR="wasmer1" go test ./...
	go clean -cache -testcache
	VMEXECUTOR="wasmer2" go test ./...

test-w1: clean
	VMEXECUTOR="wasmer1" go test ./...

test-w2: clean
	VMEXECUTOR="wasmer2" go test ./...

test-v: clean
	go test ./... -v

test-serial: clean
	go test ./... -failfast -p 1

test-short: clean
	go test ./... -short

test-short-v: clean
	go test ./... -short -v

test-short-serial:
	go test ./... -short -failfast -p 1

print-api-costs:
	@echo "bigIntOps.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/bigIntOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "baseOps.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/baseOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "managedei.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/managedei.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "manBufOps.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/manBufOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "smallIntOps.go:"
	@grep "func v1_5\|GasSchedule" vmhost/vmhooks/smallIntOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"


build-test-contracts: build-test-contracts-moapy build-test-contracts-wat

build-test-contracts-wat:
	cd test/contracts/init-simple-popcnt && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/data-drop/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-copy/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-fill/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-init/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/simd/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/imported-global/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-exceeded-max-pages/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-exceeded-pages/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-grow/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-min-pages-greater-than-max-pages/output && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/mem-multiple-max-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-multiple-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-no-max-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-no-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-single-page/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/memoryless/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/multiple-memories/output/ && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/multiple-mutable/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/noglobals/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/single-immutable/output/ && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/single-mutable/output/ && wat2wasm *.wat

build-test-contracts-moapy:
	moapy contract build --no-optimization ./test/contracts/answer
	moapy contract build ./test/contracts/async-call-builtin
	moapy contract build ./test/contracts/async-call-child
	moapy contract build ./test/contracts/async-call-parent
	moapy contract build ./test/contracts/breakpoint
	moapy contract build ./test/contracts/big-floats
	moapy contract build ./test/contracts/counter
	moapy contract build ./test/contracts/deployer
	moapy contract build ./test/contracts/deployer-child
	moapy contract build ./test/contracts/deployer-fromanother-contract
	moapy contract build ./test/contracts/deployer-parent
	moapy contract build ./test/contracts/vmhooks
	moapy contract build ./test/contracts/erc20
	moapy contract build ./test/contracts/exchange
	
	moapy contract build ./test/contracts/exec-dest-ctx-builtin
	moapy contract build ./test/contracts/exec-dest-ctx-by-caller/child
	moapy contract build ./test/contracts/exec-dest-ctx-by-caller/parent
	moapy contract build ./test/contracts/exec-dest-ctx-child
	moapy contract build ./test/contracts/exec-dest-ctx-dcdt/basic
	moapy contract build ./test/contracts/exec-dest-ctx-parent
	moapy contract build ./test/contracts/exec-dest-ctx-recursive
	moapy contract build ./test/contracts/exec-dest-ctx-recursive-child
	moapy contract build ./test/contracts/exec-dest-ctx-recursive-parent
	
	moapy contract build ./test/contracts/exec-same-ctx-child
	moapy contract build ./test/contracts/exec-same-ctx-parent
	moapy contract build ./test/contracts/exec-same-ctx-recursive
	moapy contract build ./test/contracts/exec-same-ctx-recursive-child
	moapy contract build ./test/contracts/exec-same-ctx-recursive-parent
	moapy contract build ./test/contracts/exec-same-ctx-simple-child
	moapy contract build ./test/contracts/exec-same-ctx-simple-parent
	
	moapy contract build ./test/contracts/exec-sync-ctx-multiple/alpha
	moapy contract build ./test/contracts/exec-sync-ctx-multiple/beta
	moapy contract build ./test/contracts/exec-sync-ctx-multiple/delta
	moapy contract build ./test/contracts/exec-sync-ctx-multiple/gamma
	
	moapy contract build ./test/contracts/init-correct
	moapy contract build ./test/contracts/init-simple
	moapy contract build ./test/contracts/init-wrong
	moapy contract build ./test/contracts/managed-buffers
	moapy contract build ./test/contracts/misc
	moapy contract build --no-optimization ./test/contracts/num-with-fp
	moapy contract build ./test/contracts/promises
	moapy contract build ./test/contracts/promises-train
	moapy contract build ./test/contracts/promises-tracking
	moapy contract build ./test/contracts/signatures
	moapy contract build ./test/contracts/timelocks
	moapy contract build ./test/contracts/upgrader-fromanother-contract

build-delegation:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-delegation-rs
	git clone --depth=1 --branch=master https://github.com/kalyan3104/sc-delegation-rs.git ${SANDBOX}/sc-delegation-rs
	rm -rf ${SANDBOX}/sc-delegation-rs/.git
	moapy contract build ${SANDBOX}/sc-delegation-rs
	moapy contract test --directory="tests" ${SANDBOX}/sc-delegation-rs
	cp ${SANDBOX}/sc-delegation-rs/output/delegation.wasm ./test/delegation/delegation.wasm


build-dns:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-dns-rs
	git clone --depth=1 --branch=master https://github.com/kalyan3104/sc-dns-rs.git ${SANDBOX}/sc-dns-rs
	rm -rf ${SANDBOX}/sc-dns-rs/.git
	moapy contract build ${SANDBOX}/sc-dns-rs
	moapy contract test --directory="tests" ${SANDBOX}/sc-dns-rs
	cp ${SANDBOX}/sc-dns-rs/output/dns.wasm ./test/dns/dns.wasm


build-sc-examples:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples

	moapy contract new --template=erc20-c --directory ${SANDBOX}/sc-examples erc20-c
	moapy contract build ${SANDBOX}/sc-examples/erc20-c
	cp ${SANDBOX}/sc-examples/erc20-c/output/wrc20.wasm ./test/erc20/contracts/erc20-c.wasm


build-sc-examples-rs:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples-rs
	
	moapy contract new --template=simple-coin --directory ${SANDBOX}/sc-examples-rs simple-coin
	moapy contract new --template=adder --directory ${SANDBOX}/sc-examples-rs adder
	moapy contract build ${SANDBOX}/sc-examples-rs/adder
	moapy contract build ${SANDBOX}/sc-examples-rs/simple-coin
	moapy contract test ${SANDBOX}/sc-examples-rs/adder
	moapy contract test ${SANDBOX}/sc-examples-rs/simple-coin
	cp ${SANDBOX}/sc-examples-rs/adder/output/adder.wasm ./test/adder/adder.wasm
	cp ${SANDBOX}/sc-examples-rs/simple-coin/output/simple-coin.wasm ./test/erc20/contracts/simple-coin.wasm

lint-install:
ifeq (,$(wildcard test -f bin/golangci-lint))
	@echo "Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s
endif

run-lint:
	@echo "Running golint"
	bin/golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 --timeout=2m

lint: lint-install run-lint

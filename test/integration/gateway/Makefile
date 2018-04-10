.PHONY: gen-go

gen-go:
	cd ../../.. && go install
	rm -rf ./gen-go
	frugal --gen go gateway_test.frugal && frugal --gen gateway gateway_test.frugal

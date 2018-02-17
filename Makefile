.PHONY: mocks

build:
	go build

test:
	go test --timeout 1s ./...

mocks:
	find . -name "mocks" | xargs rm -rf
	mockery -dir capability -output capability/mocks -all
	mockery -dir device/rthrottle -output device/rthrottle/mocks -all
	mockery -dir device -output device/mocks -all
	mockery -dir protocol -output protocol/mocks -name Protocol
	mockery -dir protocol/packet -output protocol/packet/mocks -all
	mockery -dir protocol/transport -output protocol/transport/mocks -all
	mockery -dir subscription/common -output subscription/common/mocks -all

cover:
	gocov test --timeout 1s ./... | gocov report

cover-html:
	gocov test --timeout 1s ./... | gocov-html > coverage.html && open coverage.html

mockery:
	go get -u github.com/vektra/mockery/.../

coverage:
	go get github.com/axw/gocov/gocov
	go get -u gopkg.in/matm/v1/gocov-html

tools:
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen


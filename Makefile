.PHONY: mocks

build:
	go build

test:
	go test -v --timeout 1s ./...

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

install_tools:
	go get github.com/axw/gocov/gocov
	go get -u gopkg.in/matm/v1/gocov-html
	go get -u github.com/vektra/mockery/.../

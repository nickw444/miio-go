.PHONY: mocks

build:
	go build

test:
	go test ./...

mocks:
	mockery -dir capability -all
	mockery -dir device -all
	mockery -dir protocol -name Protocol
	mockery -dir protocol/packet -all
	mockery -dir protocol/transport -all

#	mockery -dir protocol -output protocol/mocks -all
	mockery -dir subscription -output subscription/mocks -all

#	mockery -output subscription/target/mocks -dir subscription -name "Subscription" -recursive

cover:
	gocov test ./... | gocov report

cover-html:
	gocov test ./... | gocov-html > coverage.html && open coverage.html

mockery:
	go get -u github.com/vektra/mockery/.../

coverage:
	go get github.com/axw/gocov/gocov
	go get -u gopkg.in/matm/v1/gocov-html

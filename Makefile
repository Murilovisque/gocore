.PHONY: test-all test-postgres


test:
	go test -v ./...

test-postgres:
	$(MAKE) -C gcr/drivers/postgres test

test-all: test test-postgres

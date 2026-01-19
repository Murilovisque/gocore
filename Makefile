.PHONY: test-all test-postgres


test:
	go test -v ./...

test-postgres:
	$(MAKE) -C gcrepo/drivers/postgres test

test-all: test test-postgres

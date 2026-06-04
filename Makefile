.PHONY: test-all test-postgres


test:
	go test ./...

test-postgres:
	$(MAKE) -C gcrepo/drivers/postgres test

test-all: test test-postgres

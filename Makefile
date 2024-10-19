.PHONY: test
test:
	@go test ./... -count 1 -cover -covermode atomic -coverprofile cover.out -race -v

.PHONY: coverage
coverage:
	@go tool cover -html cover.out

.PHONY: fomat
format:
	@goimports -w ./..

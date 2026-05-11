.PHONY: fmt fmt-check vet test

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed on:" && gofmt -l . && exit 1)

vet:
	@packages="$$(go list ./... 2>/dev/null)"; \
	if [ -n "$$packages" ]; then \
		go vet $$packages; \
	else \
		echo "no Go packages, skip go vet"; \
	fi

test:
	@packages="$$(go list ./... 2>/dev/null)"; \
	if [ -n "$$packages" ]; then \
		go test $$packages; \
	else \
		echo "no Go packages, skip go test"; \
	fi

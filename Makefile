all:
	go build ./cmds/dime-a-tap

test:
	@for i in $$(find . -name '*_test.go' | xargs -n1 dirname); do \
		echo "testing $$i"; \
		go test ./$$i; \
	done

clean:
	rm -f dime-a-tap
	go clean -cache

vet:
	go vet --shadow ./...

fmt:
	go fmt ./...

.PHONY: all read clean test fmt

fmt:
	gofmt -w $$(gofmt -l ./)
vet:
	go vet -x $$(go list ./...)
test:
	go test -v $$(go list ./...)
bench:
	./bench.sh

.PHONY: test
test: build
	cloudquery sync config.yaml --log-level=debug

.PHONY: build
build:
	cd .. && make && cd _test
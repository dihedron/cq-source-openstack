.PHONY: lab
lab: build
	cloudquery sync lab.yaml --log-level=debug --telemetry-level=none

.PHONY: silver
silver: build
	cloudquery sync silver.yaml --log-level=debug --telemetry-level=none

.PHONY: build
build:
	cd .. && make && cd _test

.PHONY: reset
reset:
	rm *.log *.sqlite
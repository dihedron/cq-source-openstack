SHELL := /bin/bash

.PHONY: lab
lab: build 
	source .env/lab && cloudquery sync openstack.yaml --log-level=warn --telemetry-level=none

.PHONY: slv
slv: build
	source .env/silver && cloudquery sync openstack.yaml --log-level=warn --telemetry-level=none

.PHONY: build
build:
	cd .. && make && cd _test

.PHONY: reset
reset:
	rm -rf *.log *.sqlite

.PHONY: connect
connect:
	psql --dbname=superset --host=localhost --port=15432 --username=superset

# create database cloudquery
# create schema cloudquery
# create user cloudquery with encrypted password 'cloudquery'
# grant all privileges on database cloudquery to cloudquery
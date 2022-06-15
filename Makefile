GOPATH:=$(shell go env GOPATH)


.PHONY: api
api:
	@goctl api go -api ./api/trans.api -dir ./api


.PHONY: rpc
rpc:
	@goctl rpc proto -src ./trans/trans.proto -dir ./trans

.PHONY: model
model:
	@goctl model mysql ddl -src ./scripts/database.sql -dir ./model -c



.PHONY: run_api
run_api:
	@go run api/trans.go

.PHONY: run_rpc
run_rpc:
	@go run trans/trans.go
# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show list of make targets and their description
	@grep -E '^[/%.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:= help


.PHONY: build
build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./js_in_go ./main.go

# eg: make deploy game_svr.  svr_list: game_svr game_admin_svr ws_svr
.PHONY: deploy
deploy:
	scp -r build/${svr_name}/bin/* root@120.79.210.245:/opt/theMGame/${svr_name}/bin/

PROJECT := sinonimos

GOPATH := $(CURDIR)
GOPATHCMD=GOPATH=$(GOPATH)
GOCMD=$(GOPATHCMD) go

PROJECT_PATH=$(GOPATH)/src/$(PROJECT)

DEP=cd $(PROJECT_PATH) && GOPATH=$(GOPATH) dep

.DEFAULT_GOAL: install

.PHONY: dep-ensure dep-update dep-add dep-status install run

dep-ensure:
	@$(DEP) ensure -v

dep-update:
	@$(DEP) ensure -v -update $(PACKAGE)

dep-add:
ifdef PACKAGE
	@$(DEP) ensure -v -add $(PACKAGE)
else
	@echo "Usage: PACKAGE=<package url> make dep-add"
	@echo "The environment variable \`PACKAGE\` is not defined."
endif

dep-status:
	@$(DEP) status

install: dep-ensure
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -a -installsuffix cgo -ldflags="-w -s" -o ./bin/$(PROJECT) ./src/$(PROJECT)

run:
	@GOPATH=$(GOPATH) $(GOCMD) run $(PROJECT_PATH)/main.go $(WORD)

vet:
	@GOPATH=$(GOPATH) $(GOCMD) vet ./src/$(PROJECT)

lint: vet
	@GOPATH=$(GOPATH) golint -set_exit_status ./src/$(PROJECT)
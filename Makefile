DEV_DEPS = github.com/kardianos/govendor \
github.com/mitchellh/gox

BINNAME = casecmp
BINARY = bin/${BINNAME}
DOCKERREPO = jimeh/casecmp
BINDIR = $(shell dirname ${BINARY})
SOURCES = $(shell find . -name '*.go' -o -name 'VERSION')
VERSION = $(shell cat VERSION)
OSARCH = "darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm"
RELEASEDIR = releases

$(BINARY): $(SOURCES)
	go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"

.PHONY: build
build: $(BINARY)

.PHONY: clean
clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi; \
	if [ -d ${BINDIR} ]; then rmdir ${BINDIR}; fi

.PHONY: run
run: $(BINARY)
	$(BINARY)

.PHONY: deps
deps:
	@govendor sync

.PHONY: dev-deps
dev-deps:
	@$(foreach DEP,$(DEV_DEPS),go get $(DEP);)

.PHONY: update-dev-deps
update-dev-deps:
	@$(foreach DEP,$(DEV_DEPS),go get -u $(DEP);)

.PHONY: release-build
release-build: deps
	gox -output "${RELEASEDIR}/${BINNAME}_${VERSION}_{{.OS}}_{{.Arch}}" \
		-osarch=${OSARCH} \
		-ldflags "-X main.Version=${VERSION}"

.SILENT: release
.PHONY: release
release: release-build
	$(eval BINS := $(shell cd ${RELEASEDIR} && find . \
		-name "${BINNAME}_${VERSION}_*" -not -name "*.tar.gz"))
	cd $(RELEASEDIR); \
		$(foreach BIN,$(BINS),tar -cvzf $(BIN).tar.gz $(BIN) && rm $(BIN);)

.PHONY: docker
docker: clean deps
	docker build -t "${DOCKERREPO}:latest" . \
		&& docker tag "${DOCKERREPO}:latest" "${DOCKERREPO}:${VERSION}"

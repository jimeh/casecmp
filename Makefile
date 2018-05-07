NAME = casecmp
BINARY = bin/${NAME}

VERSION ?= $(shell cat VERSION)
COMMIT = $(shell git show --format="%h" --no-patch)
DATE = $(shell date +%Y-%m-%dT%T%z)

SOURCES = $(shell find . \
	-name '*.go' \
	-o -name 'LICENSE' \
	-o -name 'Makefile' \
	-o -name 'README.md' \
	-o -name 'VERSION')

RELEASE_DIR = releases
RELEASE_TARGETS = \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_darwin_386.tar.gz \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_darwin_amd64.tar.gz \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_freebsd_386.tar.gz \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_freebsd_amd64.tar.gz \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_freebsd_arm.tar.gz \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_linux_386.tar.gz \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_linux_amd64.tar.gz \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_linux_arm.tar.gz \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_windows_386.zip \
	$(RELEASE_DIR)/$(NAME)-$(VERSION)_windows_amd64.zip
RELEASE_ASSETS = \
	README.md \
	LICENSE

$(BINARY): $(SOURCES)
	go build -o ${BINARY} -ldflags \ "\
		-X main.version=${VERSION} \
		-X main.commit=${COMMIT} \
		-X main.date=${DATE}"

.PHONY: build
build: $(BINARY)

.PHONY: run
run: $(BINARY)
	$(BINARY)

.PHONY: clean
clean:
	$(eval BIN_DIR := $(shell dirname ${BINARY}))
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi
	if [ -d ${BIN_DIR} ]; then rmdir ${BIN_DIR}; fi

.PHONY: docker
docker:
	docker build -t "$(shell whoami)/$(NAME)" .

.PHONY: release
release: $(RELEASE_TARGETS)

$(RELEASE_DIR)/$(NAME)-$(VERSION)_%.tar.gz: $(SOURCES)
	$(eval OS := $(word 1, $(subst _, ,$*)))
	$(eval ARCH := $(word 2, $(subst _, ,$*)))
	$(eval TARGET := $(NAME)-$(VERSION)_$*)
	mkdir -p "$(TARGET)" \
		&& env GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -a \
			-o "$(TARGET)/$(NAME)" -ldflags "-X main.Version=$(VERSION)" \
		&& cp $(RELEASE_ASSETS) "$(TARGET)/" \
		&& tar -cvzf "$@" "$(TARGET)" \
		&& cd "$(TARGET)" && rm "$(NAME)" $(RELEASE_ASSETS) && cd .. \
		&& rmdir "$(TARGET)"

$(RELEASE_DIR)/$(NAME)-$(VERSION)_windows_%.zip: $(SOURCES)
	$(eval TARGET := $(NAME)-$(VERSION)_windows_$*)
	mkdir -p "$(TARGET)" \
		&& env GOOS=windows GOARCH=$* CGO_ENABLED=0 go build -a \
			-o "$(TARGET)/$(NAME).exe" -ldflags "-X main.Version=$(VERSION)" \
		&& cp $(RELEASE_ASSETS) "$(TARGET)/" \
		&& zip -r "$@" "$(TARGET)" \
		&& cd "$(TARGET)" && rm "$(NAME).exe" $(RELEASE_ASSETS) && cd .. \
		&& rmdir "$(TARGET)"

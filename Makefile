# reference https://github.com/Dreamacro/clash/blob/master/Makefile

NAME=clash
BINDIR=bin
VERSION=$(shell git describe --tags || echo "unknown version")
CLASHCORENAME=123
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s -X "github.com/csj8520/clash-for-flutter-service/constant.Version=$(VERSION)"'

PLATFORM_LIST = \
	darwin-amd64 \
	linux-amd64

WINDOWS_ARCH_LIST = \
	windows-amd64

all: darwin-amd64 linux-amd64 windows-amd64

# darwin-arm64:
# 	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-service

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-service

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-service

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-service.exe



gz_releases=$(addsuffix .gz, $(PLATFORM_LIST))
zip_releases=$(addsuffix .zip, $(WINDOWS_ARCH_LIST))

$(gz_releases): %.gz : %
	# chmod +x $(BINDIR)/$(NAME)-$(basename $@)-service
	gzip -f -S -$(VERSION).gz $(BINDIR)/$(NAME)-$(basename $@)-service

$(zip_releases): %.zip : %
	zip -m -j $(BINDIR)/$(NAME)-$(basename $@)-service-$(VERSION).zip $(BINDIR)/$(NAME)-$(basename $@)-service.exe

releases: $(gz_releases) $(zip_releases)

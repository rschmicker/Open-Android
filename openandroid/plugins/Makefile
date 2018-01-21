PLUGINS := Sha256/Sha256.go Apis/Apis.go Intents/Intents.go \
			IsMalicious/IsMalicious.go Md5/Md5.go PackageName/PackageName.go \
			PackageVersion/PackageVersion.go Permissions/Permissions.go \
			Sha1/Sha1.go Strings/Strings.go
PREFIX=go build -buildmode=plugin

.PHONY: all clean

all: $(PLUGINS)
	@for plugin in $(PLUGINS); do echo "Building Plugin $$plugin"; $(PREFIX) $$plugin; done

clean:
	@rm -f *.so
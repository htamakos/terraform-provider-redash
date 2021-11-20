.PHONY: all format vet tidy test clean

# -----------------------------------------------------------------------------
#  CONSTANTS
# -----------------------------------------------------------------------------

version = `cat VERSION`

src_dir := terraform-provider-redash

build_dir = build

coverage_dir  = $(build_dir)/coverage
coverage_out  = $(coverage_dir)/coverage.out
coverage_html = $(coverage_dir)/coverage.html

output_dir    = $(build_dir)/output

linux_dir     = $(output_dir)/linux
darwin_dir    = $(output_dir)/darwin
windows_dir   = $(output_dir)/windows

bin_name      = terraform-provider-redash_v$(version)
bin_linux     = $(linux_dir)/$(bin_name)
bin_darwin    = $(darwin_dir)/$(bin_name)
bin_windows   = $(windows_dir)/$(bin_name)

# -----------------------------------------------------------------------------
#  BUILDING
# -----------------------------------------------------------------------------

all:
	go install github.com/mitchellh/gox@latest
	gox -osarch=linux/amd64 -output=$(bin_linux) ./$(src_dir)
	gox -osarch=darwin/amd64 -output=$(bin_darwin) ./$(src_dir)
	gox -osarch=windows/amd64 -output=$(bin_windows) ./$(src_dir)

# -----------------------------------------------------------------------------
#  FORMATTING
# -----------------------------------------------------------------------------

format:
	go fmt ./$(src_dir)
	gofmt -s -w ./$(src_dir)

vet:
	go vet ./$(src_dir)

tidy:
	go mod tidy

# -----------------------------------------------------------------------------
#  TESTING
# -----------------------------------------------------------------------------

test:
	mkdir -p $(coverage_dir)
	go install golang.org/x/tools/cmd/cover/...@latest
	go test ./$(src_dir) -tags test -v -covermode=count -coverprofile=$(coverage_out)
	go tool cover -html=$(coverage_out) -o $(coverage_html)

# -----------------------------------------------------------------------------
#  CLEANUP
# -----------------------------------------------------------------------------

clean:
	rm -rf $(build_dir)
	
TARGETS     = logshipper
BUILD_DIR   = ./build

all: dependencies $(TARGETS)

$(BUILD_DIR):
	mkdir -p "$(BUILD_DIR)"

dependencies:
	go get -v ./...

$(TARGETS): $(BUILD_DIR)
	go build -v -o $(BUILD_DIR)/$@ cmd/$@/main.go

clean:
	[[ -d "$(BUILD_DIR)" ]] && rm -rf "$(BUILD_DIR)" || true

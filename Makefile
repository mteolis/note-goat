VERSION := 1.0.0
APP_NAME := NoteGoat
BUILD_DIR := build

build: clean
	mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-$(VERSION).exe .

clean:
	rm -rf $(BUILD_DIR)

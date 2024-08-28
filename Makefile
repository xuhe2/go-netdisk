# set the go project name
BINARY_NAME=go-netdisk

# build the go binary
$(BINARY_NAME): 
	@echo "Building ${BINARY_NAME}..."
	@go build -o bin/${BINARY_NAME}

.PHONY: run run-build clean-bin test
# run the go binary
run: $(BINARY_NAME)
	@bin/${BINARY_NAME} $(ARGS)

# clean up the build artifacts
clean-bin:
	rm -f bin/${BINARY_NAME}

# clean the data
# 删除所有在data文件夹下并且用.part/.info结尾的文件
clean-data:
	rm -f data/*.part
	rm -f data/*.info

# clean up the build artifacts and data
clean: clean-bin clean-data

# test the go code
test:
	@go test ./...
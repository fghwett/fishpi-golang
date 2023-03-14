# 声明编译器变量
GO=go
# 声明编译参数变量
GOFLAGS=

# 声明编译目标变量
BINARY=/Users/fghwett/go/bin/fishpi-golang

# 声明源文件变量
SOURCES=$(wildcard *.go)

# 定义默认目标
all: $(BINARY)

# 定义编译目标
$(BINARY): $(SOURCES)
	$(GO) build $(GOFLAGS) -o $(BINARY)

# 定义清理目标
clean:
	rm -f $(BINARY)

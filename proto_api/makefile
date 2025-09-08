# Каталог, где лежат .proto файлы
PROTO_DIR = grpc/v1

# Протобаф-файлы
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

# Флаги генерации
GO_OUT_FLAGS = \
	--go_out=. \
	--go-grpc_out=. 

# Цель по умолчанию
.PHONY: all
all: generate

# Генерация protobuf
.PHONY: generate
generate:
	protoc $(GO_OUT_FLAGS) $(PROTO_FILES)

# Очистка сгенерированных файлов (если нужно)
.PHONY: clean
clean:
	rm -f $(PROTO_DIR)/*.pb.go
	rm -f $(PROTO_DIR)/*_grpc.pb.go

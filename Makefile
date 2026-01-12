gen-account:
	mkdir -p ./account/pb
	cd account && \
		protoc \
			--go_out=./pb --go_opt=paths=source_relative \
			--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
			account.proto

gen-catalog:
	mkdir -p ./catalog/pb
	cd catalog && \
		protoc \
			--go_out=./pb --go_opt=paths=source_relative \
			--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
			catalog.proto

gen-order:
	mkdir -p ./order/pb
	cd order && \
		protoc \
			--go_out=./pb --go_opt=paths=source_relative \
			--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
			order.proto
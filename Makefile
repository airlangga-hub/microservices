gen-account:
	mkdir -p ./services/account/pb
	cd services/account && \
		protoc \
			--go_out=./pb --go_opt=paths=source_relative \
			--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
			account.proto
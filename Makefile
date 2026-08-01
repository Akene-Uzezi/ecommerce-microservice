.PHONY: gen clean

gen:
	@mkdir -p api/gen/order api/gen/stock api/gen/payment

	@protoc -I=api/proto \
		--go_out=api/gen/order --go_opt=paths=source_relative \
		--go-grpc_out=api/gen/order --go-grpc_opt=paths=source_relative \
		order.proto

	@protoc -I=api/proto \
		--go_out=api/gen/stock --go_opt=paths=source_relative \
		--go-grpc_out=api/gen/stock --go-grpc_opt=paths=source_relative \
		stock.proto

	@protoc -I=api/proto \
		--go_out=api/gen/payment --go_opt=paths=source_relative \
		--go-grpc_out=api/gen/payment --go-grpc_opt=paths=source_relative \
		payment.proto

	@echo "Protobuf stubs generated cleanly!"

clean:
	rm -rf api/gen/*
	rm -f api/proto/*.pb.go

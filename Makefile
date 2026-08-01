.PHONY: gen

gen:
		@protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			api/proto/order.proto --go_out=api/gen/order --go-grpc_out=api/gen/order
		@protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			api/proto/stock.proto --go_out=api/gen/stock --go-grpc_out=api/gen/stock
		@protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			api/proto/payment.proto --go_out=api/gen/payment --go-grpc_out=api/gen/payment

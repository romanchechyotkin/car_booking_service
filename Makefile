swag:
	$(GOPATH)/bin/swag init;

bombardier:
	cd cmd/main && $(GOPATH)/bin/bombardier -c 100 -d 60s http://localhost:5000/users;

run_kafka-zookeeper:
	cd /var/lib/kafka && bin/zookeeper-server-start.sh config/zookeeper.properties;

run_kafka-server:
	cd /var/lib/kafka && bin/kafka-server-start.sh config/server.properties;

create_kafka-topics:
	cd /var/lib/kafka && bin/kafka-topics.sh --create --topic emails --bootstrap-server localhost:9092 && bin/kafka-topics.sh --create --topic payments --bootstrap-server localhost:9092;

run_postgresql:
	cd /var/lib/postgresql/bin/ && ./pg_ctl -D /var/lib/postgresql/main start;

run_email_sender_microservice:
	cd ../email_sender_microservice && $(GOPATH)/go run cmd/main/main.go;

run_payment_microservice:
	cd ../payment_microservice && $(GOPATH)/go run cmd/main/main.go;

test-monorepo:
	cd /home/chechyotka/projects/golang_projects/car_booking_service/monorepo/ && $(GOPATH)/go test ./...;

test-email_sender_microservice:
	cd /home/chechyotka/projects/golang_projects/car_booking_service/email_sender_microservice/ && $(GOPATH)/go test ./...;

run_postman:
	cd /var/lib/Postman/ && ./Postman;

proto:
	protoc -I ./internal/user/proto ./internal/user/proto/proto.proto --go_out=./internal/user/proto/pb --go-grpc_out=./internal/user/proto/pb

prometheus:
	cd /var/lib/prometheus && ./prometheus --config.file=/home/chechyotka/projects/golang_projects/car_booking_service/monorepo/prometheus.yml

build_project:
	echo $(GOPATH) && $(GOPATH)/bin/swag init && $(GOPATH)/go build -o ./build/bin ./main.go && ./build/bin
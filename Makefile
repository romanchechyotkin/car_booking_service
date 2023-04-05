swag:
	cd cmd/main && $(GOPATH)/bin/swag init -o "../../docs";

bombardier:
	cd cmd/main && $(GOPATH)/bin/bombardier -c 100 -d 60s http://localhost:5000/users;

kafka-zookeeper:
	cd /var/lib/kafka && bin/zookeeper-server-start.sh config/zookeeper.properties;

kafka-server:
	cd /var/lib/kafka && bin/kafka-server-start.sh config/server.properties;

kafka-topic:
	cd /var/lib/kafka && bin/kafka-topics.sh --create --topic emails --bootstrap-server localhost:9092;

postgresql:
	cd /var/lib/postgresql/bin/ && ./pg_ctl -D /var/lib/postgresql/main start;

email_sender_microservice:
	cd ../email_sender_microservice && $(GOPATH)/go run cmd/main/main.go;

test-monorepo:
	cd /home/chechyotka/projects/golang_projects/car_booking_service/monorepo/ && $(GOPATH)/go test ./...;

test-email_sender_microservice:
	cd /home/chechyotka/projects/golang_projects/car_booking_service/email_sender_microservice/ && $(GOPATH)/go test ./...;

postman:
	cd /var/lib/Postman/ && ./Postman;
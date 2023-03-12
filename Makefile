swag:
	cd cmd/main && $(GOPATH)/bin/swag init -o "../../docs";

bombardier:
	cd cmd/main && $(GOPATH)/bin/bombardier -c 100 -d 60s http://localhost:5000/users;

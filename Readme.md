1) run local postgres on port 5432
2) run local prometheus using  
 --config.file=/home/chechyotka/projects/golang_projects/car_booking_service/prometheus.yml on port 9090
3) run service on port 5000 (on port 5001 running metrics server) 
4) $GOPATH/bin/bombardier -c 100 -d 60s http://localhost:5000/users  run performance test
5) $GOPATH/bin/swag init -o "../../docs" generate swagger docs
6) TODO set grafana for postgresql
 

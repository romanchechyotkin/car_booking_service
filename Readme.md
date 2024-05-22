```env
    MINIO_HOST=minio
    MINIO_PORT=9000
    MINIO_USER=minio
    MINIO_PASSWORD=minio123
    POSTGRES_HOST=postgres
    POSTGRES_PORT=5432
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=5432
    POSTGRES_DATABASE=car_booking_service
```


Database scheme: https://dbdiagram.io/d/63ea01e6296d97641d80732a
// invalid now, database was increased 

working with microservice through Kafka
https://github.com/romanchechyotkin/email_sender_microservice

2) run local prometheus using  
 --config.file=/home/chechyotka/projects/golang_projects/car_booking_service/prometheus.yml on port 9090
3) run service on port 5000 (on port 5001 running metrics server) 
4) TODO set grafana for postgresql

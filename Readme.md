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

run
```env
    docker compose up --build
    go run main.go
```


Database scheme: https://dbdiagram.io/d/car-booking-service-6650e95bf84ecd1d2216d37a

working with microservice through Kafka
https://github.com/romanchechyotkin/email_sender_microservice

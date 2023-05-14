FROM postgres

ENV POSTGRES_DB car_booking_service
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD 5432

COPY psql_dump.sql /docker-entrypoint-initdb.d/

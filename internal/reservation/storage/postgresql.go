package reservation

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	reservation "github.com/romanchechyotkin/car_booking_service/internal/reservation/model"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"
	"log"
	"time"
)

const (
	DDMMYYYY = "02.01.2006"
)

type Repository struct {
	client *pgxpool.Pool
}

func NewRepository(client *pgxpool.Pool) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) CreateReservation(ctx context.Context, dto reservation.Dto) error {
	query := `
		INSERT INTO public.reservations (customer_id, seller_id, car_id, start_date, end_date, total_price) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	startDate, err := time.Parse(DDMMYYYY, dto.StartDate)
	if err != nil {
		return err
	}
	endDate, err := time.Parse(DDMMYYYY, dto.EndDate)
	if err != nil {
		return err
	}

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	exec, err := r.client.Exec(ctx, query, dto.CustomerId, dto.CarOwnerId, dto.Car.Id, startDate, endDate, dto.TotalPrice)
	if err != nil {
		log.Printf("error %v", err)
		return err
	}
	log.Println(exec.RowsAffected())

	return err
}

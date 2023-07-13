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
	hhDDMMYYYY = "15 02.01.2006"
)

type Repository struct {
	client *pgxpool.Pool
}

type Storage interface {
	CreateReservation(ctx context.Context, dto reservation.Dto) error
	GetReservationDates(ctx context.Context, id string) (dates []reservation.TimeFromDB, err error)
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

	startDate, err := time.Parse(hhDDMMYYYY, dto.StartDate)
	if err != nil {
		return err
	}
	endDate, err := time.Parse(hhDDMMYYYY, dto.EndDate)
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

func (r *Repository) GetReservationDates(ctx context.Context, id string) (dates []reservation.TimeFromDB, err error) {
	query := `
		SELECT start_date, end_date
		FROM public.reservations
		WHERE car_id = $1 
		ORDER BY start_date
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return dates, err
	}

	for rows.Next() {
		var dto reservation.TimeFromDB

		err = rows.Scan(&dto.StartDate, &dto.EndDate)
		if err != nil {
			return dates, err
		}

		dates = append(dates, dto)
	}

	return dates, err
}

//
//func (r *Repository) GetAllCarReservations(ctx context.Context, id string) (res []reservation.TimeFromDB, err error) {
//	query := `
//		SELECT start_date, end_date
//		FROM public.reservations
//		WHERE car_id = $1
//	`
//
//	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
//	rows, err := r.client.Query(ctx, query, id)
//	if err != nil {
//		return dates, err
//	}
//
//	for rows.Next() {
//		var dto reservation.TimeFromDB
//
//		err = rows.Scan(&dto.StartDate, &dto.EndDate)
//		if err != nil {
//			return dates, err
//		}
//
//		dates = append(dates, dto)
//	}
//
//	return dates, err
//}

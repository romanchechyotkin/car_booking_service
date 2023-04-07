package cars_storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"
	"log"
)

type Storage interface {
	CreateCar(ctx context.Context, car *car.CreateCarFormDto, id string) error
	GetCar(ctx context.Context, id string) (c *car.Car, err error)
}

type Repository struct {
	client *pgxpool.Pool
}

func NewRepository(client *pgxpool.Pool) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) CreateCar(ctx context.Context, car *car.CreateCarFormDto, userId string) error {
	conn, err := r.client.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	carsQuery := `
		INSERT INTO public.cars (id, brand, model, year, price_per_day) 
		VALUES ($1, $2, $3, $4, $5)
	`

	carsUsersQuery := `
		INSERT INTO public.cars_users (car_id, user_id) 
		VALUES ($1, $2)
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(carsQuery))
	row, _ := tx.Exec(ctx, carsQuery, car.Id, car.Brand, car.Model, car.Year, car.PricePerDay)
	fmt.Println(row.RowsAffected())
	if row.RowsAffected() == 0 {
		return errors.New("wrong cars numbers")
	}

	log.Printf("SQL query: %s", postgresql.FormatQuery(carsUsersQuery))
	row, _ = tx.Exec(ctx, carsUsersQuery, car.Id, userId)
	fmt.Println(row.RowsAffected())

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetCar(ctx context.Context, id string) (c car.Car, err error) {
	carQuery := `
		SELECT id, brand, model, price_per_day, year, is_available, rating
		FROM public.cars
		WHERE id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(carQuery))
	err = r.client.QueryRow(ctx, carQuery, id).Scan(&c.Id, &c.Brand, &c.Model, &c.PricePerDay, &c.Year, &c.IsAvailable, &c.Rating)
	if err != nil {
		log.Println(err)
		return c, err
	}

	imagesQuery := `
		SELECT url FROM car_images WHERE car_id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(imagesQuery))
	rows, err := r.client.Query(ctx, imagesQuery, c.Id)
	if err != nil {
		log.Println(err)
		return c, err
	}
	defer rows.Close()

	images := make([]string, 0)
	for rows.Next() {
		var i car.Image

		err = rows.Scan(&i.Url)
		if err != nil {
			log.Println(err)
			return c, err
		}

		images = append(images, i.Url)
	}

	c.Images = images

	return c, nil
}

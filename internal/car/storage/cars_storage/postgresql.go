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
	GetAllCars(ctx context.Context) error
	GetCar(ctx context.Context, id string) (c *car.Car, err error)
	GetCarOwner(ctx context.Context, id string)
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

func (r *Repository) GetAllCars(ctx context.Context) ([]car.GetCarDto, error) {
	var query = `
		SELECT cars.id, cars.brand, cars.model, cars.price_per_day, cars.year, cars.is_available, cars.rating, cu.user_id
		FROM public.cars
		INNER JOIN cars_users cu on cars.id = cu.car_id
	`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	cars := make([]car.GetCarDto, 0)
	for rows.Next() {
		var c car.GetCarDto
		err = rows.Scan(&c.Car.Id, &c.Car.Brand, &c.Car.Model, &c.Car.PricePerDay, &c.Car.Year, &c.Car.IsAvailable, &c.Car.Rating, &c.UserId)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		imgQuery := `
			SELECT url FROM car_images WHERE car_id = $1
		`
		fmt.Println(c.Car)

		r, err := r.client.Query(ctx, imgQuery, c.Car.Id)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		images := make([]string, 0)
		for r.Next() {
			var i car.Image

			err = r.Scan(&i.Url)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			images = append(images, i.Url)
		}

		c.Car.Images = images
		cars = append(cars, c)
	}

	return cars, err
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

func (r *Repository) GetCarOwner(ctx context.Context, id string) (userId string, err error) {
	query := `
		SELECT user_id
		FROM public.cars_users
		WHERE car_id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	err = r.client.QueryRow(ctx, query, id).Scan(&userId)
	if err != nil {
		log.Println(err)
		return userId, err
	}

	return userId, nil
}

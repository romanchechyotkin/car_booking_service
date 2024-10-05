package cars_storage

import (
	"context"
	"fmt"
	"log"

	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	SORT_BY_ASC_PRICE  = "prc.a"
	SORT_BY_DESC_PRICE = "prc.d"
	SORT_BY_ASC_YEAR   = "y.a"
	SORT_BY_DESC_YEAR  = "y.d"
)

type Storage interface {
	CreateCar(ctx context.Context, car *car.CreateCarFormDto, userId string) error
	GetAllUserCars(ctx context.Context, userID string) ([]car.Car, error)
	GetAllUserCarsAmount(ctx context.Context, userID string) (int, error)
	GetAllCars(ctx context.Context, opt ...string) ([]car.GetCarDto, error)
	GetCar(ctx context.Context, id string) (c car.Car, err error)
	GetCarOwner(ctx context.Context, id string) (userId string, err error)
	ChangeIsAvailable(ctx context.Context, id string) error
	GetAllCarRatings(ctx context.Context, id string) ([]car.GetAllCarRatingsDto, error)
	CreateRating(ctx context.Context, dto user.RateDto, carId, ratedBy string) error
	GetAmountCarRatings(ctx context.Context, carId string) (amount float32, sum float32, err error)
	ChangeCarRating(ctx context.Context, id string, rating float32) error
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
		INSERT INTO public.cars (id, brand, model, year, price_per_day, is_automatic, seats) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	carsUsersQuery := `
		INSERT INTO public.cars_users (car_id, user_id) 
		VALUES ($1, $2)
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(carsQuery))
	row, err := tx.Exec(ctx, carsQuery, car.Id, car.Brand, car.Model, car.Year, car.PricePerDay, car.IsAutomatic, car.Seats)
	if err != nil {
		return err
	}

	log.Printf("SQL query: %s", postgresql.FormatQuery(carsUsersQuery))
	row, _ = tx.Exec(ctx, carsUsersQuery, car.Id, userId)
	fmt.Println(row.RowsAffected())

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllCars(ctx context.Context, opt ...string) ([]car.GetCarDto, error) {

	var orderBy, query string

	if len(opt) != 0 {
		log.Println(opt[0])
		orderBy = opt[0]
	}

	switch orderBy {
	case SORT_BY_ASC_PRICE:
		query = `
		SELECT cars.id, cars.brand, cars.model, cars.price_per_day, cars.year, cars.is_available, cars.rating, cars.seats, cars.is_automatic, cars.created_at, cu.user_id
		FROM public.cars
		INNER JOIN cars_users cu on cars.id = cu.car_id		
		ORDER BY price_per_day
	`
	case SORT_BY_DESC_PRICE:
		query = `
		SELECT cars.id, cars.brand, cars.model, cars.price_per_day, cars.year, cars.is_available, cars.rating, cars.seats, cars.is_automatic, cars.created_at, cu.user_id
		FROM public.cars
		INNER JOIN cars_users cu on cars.id = cu.car_id		
		ORDER BY price_per_day DESC 
	`
	case SORT_BY_ASC_YEAR:
		query = `
		SELECT cars.id, cars.brand, cars.model, cars.price_per_day, cars.year, cars.is_available, cars.rating, cars.seats, cars.is_automatic, cars.created_at, cu.user_id
		FROM public.cars
		INNER JOIN cars_users cu on cars.id = cu.car_id		
		ORDER BY year
	`
	case SORT_BY_DESC_YEAR:
		query = `
		SELECT cars.id, cars.brand, cars.model, cars.price_per_day, cars.year, cars.is_available, cars.rating, cars.seats, cars.is_automatic, cars.created_at, cu.user_id
		FROM public.cars
		INNER JOIN cars_users cu on cars.id = cu.car_id		
		ORDER BY year DESC 
	`
	default:
		query = `
		SELECT cars.id, cars.brand, cars.model, cars.price_per_day, cars.year, cars.is_available, cars.rating, cars.seats, cars.is_automatic, cars.created_at, cu.user_id
		FROM public.cars
		INNER JOIN cars_users cu on cars.id = cu.car_id		
		ORDER BY created_at 
	`
	}

	log.Println(query)

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	cars := make([]car.GetCarDto, 0)
	for rows.Next() {
		var c car.GetCarDto
		err = rows.Scan(&c.Car.Id, &c.Car.Brand, &c.Car.Model, &c.Car.PricePerDay, &c.Car.Year, &c.Car.IsAvailable, &c.Car.Rating, &c.Car.Seats, &c.Car.IsAutomatic, &c.Car.CreatedAt, &c.UserId)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		imgQuery := `
			SELECT url FROM car_images WHERE car_id = $1
		`

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

func (r *Repository) GetAllUserCars(ctx context.Context, userID string) ([]car.Car, error) {
	query := `
		select c.id, c.brand, c.model, c.price_per_day, c.year, c.is_available, c.rating, c.seats, c.is_automatic, c.created_at
		from public.cars c
		INNER JOIN cars_users u on u.car_id = c.id
	 	WHERE u.user_id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	rows, err := r.client.Query(ctx, query, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	imagesQuery := `
		SELECT url FROM car_images WHERE car_id = $1
	`

	var cars []car.Car
	for rows.Next() {
		var c car.Car

		err := rows.Scan(&c.Id, &c.Brand, &c.Model, &c.PricePerDay, &c.Year, &c.IsAvailable, &c.Rating, &c.Seats, &c.IsAutomatic, &c.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		log.Printf("SQL query: %s", postgresql.FormatQuery(imagesQuery))
		imagesRows, err := r.client.Query(ctx, imagesQuery, c.Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer imagesRows.Close()

		images := make([]string, 0)
		for imagesRows.Next() {
			var i car.Image

			err = imagesRows.Scan(&i.Url)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			images = append(images, i.Url)
		}

		c.Images = images
		cars = append(cars, c)
	}

	return cars, nil
}

func (r *Repository) GetAllUserCarsAmount(ctx context.Context, userID string) (int, error) {
	query := `
		select count(*) 
		from public.cars c
		INNER JOIN cars_users u on u.car_id = c.id
	 	WHERE u.user_id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))

	var amount int
	err := r.client.QueryRow(ctx, query, userID).Scan(&amount)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return amount, nil
}

func (r *Repository) GetCar(ctx context.Context, id string) (c car.Car, err error) {
	carQuery := `
		SELECT id, brand, model, price_per_day, year, is_available, rating, seats, is_automatic, created_at
		FROM public.cars
		WHERE id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(carQuery))
	err = r.client.QueryRow(ctx, carQuery, id).Scan(&c.Id, &c.Brand, &c.Model, &c.PricePerDay, &c.Year, &c.IsAvailable, &c.Rating, &c.Seats, &c.IsAutomatic, &c.CreatedAt)
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

func (r *Repository) ChangeIsAvailable(ctx context.Context, id string) error {
	query := `
		UPDATE public.cars
		SET is_available = false
		WHERE id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	exec, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	log.Printf("%d", exec.RowsAffected())
	return nil
}

func (r *Repository) GetAllCarRatings(ctx context.Context, id string) ([]car.GetAllCarRatingsDto, error) {
	query := `
		SELECT r.rate, r.comment, u.full_name, r.created_at
		FROM cars_ratings r
		INNER JOIN users u on u.id = r.rate_by_user
	 	WHERE r.car_id = $1
		ORDER BY r.created_at desc
`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ratings []car.GetAllCarRatingsDto
	for rows.Next() {
		var rate car.GetAllCarRatingsDto

		err = rows.Scan(&rate.Rating, &rate.Comment, &rate.User, &rate.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		ratings = append(ratings, rate)
	}

	return ratings, nil
}

func (r *Repository) CreateRating(ctx context.Context, dto user.RateDto, carId, ratedBy string) error {
	query := `
		INSERT INTO public.cars_ratings  (rate, comment, car_id, rate_by_user)
		VALUES ($1, $2, $3, $4)
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	exec, err := r.client.Exec(ctx, query, dto.Rating, dto.Comment, carId, ratedBy)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(exec.RowsAffected())

	return nil
}

func (r *Repository) GetAmountCarRatings(ctx context.Context, carId string) (amount float32, sum float32, err error) {
	query := `
		SELECT count(*), sum(rate) FROM cars_ratings WHERE car_id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	_ = r.client.QueryRow(ctx, query, carId).Scan(&amount, &sum)

	return amount, sum, nil
}

func (r *Repository) ChangeCarRating(ctx context.Context, id string, rating float32) error {
	query := `
		UPDATE public.cars
		SET rating = $1
		WHERE id = $2
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	exec, err := r.client.Exec(ctx, query, rating, id)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(exec.RowsAffected())

	return nil
}

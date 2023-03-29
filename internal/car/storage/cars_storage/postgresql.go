package cars_storage

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"

	"context"
	"fmt"
	"log"
)

type Storage interface {
	CreateCar(ctx context.Context, car *car.CreateCarFormDto, id string) error
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

	log.Printf("SQL query: %s", postgresql.FormatQuery(carsUsersQuery))
	row, _ = tx.Exec(ctx, carsUsersQuery, car.Id, userId)
	fmt.Println(row.RowsAffected())

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateUser(ctx context.Context, user *user.CreateUserDto) error {
	query := `
		INSERT INTO public.users (email, password, full_name, telephone_number) 
		VALUES ($1, $2, $3, $4)
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	_, err := r.client.Exec(ctx, query, user.Email, user.Password, user.FullName, user.TelephoneNumber)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL error: %s, Detail: %s, Where: %s, Code: %s, SQL State: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Println(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]user.GetUsersDto, error) {
	query := `
		SELECT id, email, full_name, telephone_number, is_premium, city, rating
		FROM public.users
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	rows, err := r.client.Query(ctx, query)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL error: %s, Detail: %s, Where: %s, Code: %s, SQL State: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Println(newErr)
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	users := make([]user.GetUsersDto, 0)
	for rows.Next() {
		var u user.GetUsersDto

		err = rows.Scan(&u.Id, &u.Email, &u.FullName, &u.TelephoneNumber, &u.IsPremium, &u.City, &u.Rating)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, u)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetOneUserById(ctx context.Context, id string) (u user.GetUsersDto, err error) {
	query := `
		SELECT id, email, password, full_name, telephone_number, is_premium, city, rating 
		FROM public.users
		WHERE id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	err = r.client.QueryRow(ctx, query, id).Scan(&u.Id, &u.Email, &u.Password, &u.FullName, &u.TelephoneNumber, &u.IsPremium, &u.City, &u.Rating)
	if err != nil {
		log.Println(err)
		return u, err
	}

	return u, nil
}

func (r *Repository) GetOneUserByEmail(ctx context.Context, email string) (u user.GetUsersDto, err error) {
	query := `
		SELECT id, email, password, full_name, telephone_number, is_premium, city, rating 
		FROM public.users
		WHERE email = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	err = r.client.QueryRow(ctx, query, email).Scan(&u.Id, &u.Email, &u.Password, &u.FullName, &u.TelephoneNumber, &u.IsPremium, &u.City, &u.Rating)
	if err != nil {
		log.Printf("err: %v", err)
		return u, err
	}

	return u, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id string, user *user.UpdateUserDto) error {
	query := `
		UPDATE public.users
		SET email = $1,
		    password = $2,
		    full_name = $3,
		    city = $4
		WHERE id = $5
	`

	exec, err := r.client.Exec(ctx, query, user.Email, user.Password, user.FullName, user.City, id)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected := exec.RowsAffected()
	log.Printf("rows affected: %d", rowsAffected)
	return nil
}

func (r *Repository) DeleteUserById(ctx context.Context, id string) error {
	query := `
		DELETE FROM public.users
		WHERE id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	exec, err := r.client.Exec(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected := exec.RowsAffected()
	log.Printf("after delete rows affected: %d", rowsAffected)

	return nil
}

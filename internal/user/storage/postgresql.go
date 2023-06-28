package user

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"

	"context"
	"fmt"
	"log"
)

type Storage interface {
	CreateUser(ctx context.Context, user *user.CreateUserDto) error
	GetRole(ctx context.Context, id string) (string, error)
	GetAllUsers(ctx context.Context) ([]user.GetUsersDto, error)
	GetOneUserById(ctx context.Context, id string) (u user.GetUsersDto, err error)
	GetOneUserByEmail(ctx context.Context, email string) (u user.GetUsersDto, err error)
	UpdateUser(ctx context.Context, id string, user *user.UpdateUserDto) error
	DeleteUserById(ctx context.Context, id string) error
	CreateRating(ctx context.Context, dto user.RateDto, userId, ratedBy string) error
	GetUserRatings(ctx context.Context, userId string) (amount float32, sum float32, err error)
	ChangeUserRating(ctx context.Context, id string, rating float32) error
	GetAllUserRatings(ctx context.Context, id string) ([]user.GetAllRatingsDto, error)
	CreateApplication(ctx context.Context, id string, filename string) error
	GetApplications(ctx context.Context) ([]user.ApplicationDto, error)
	ChangeUserVerify(ctx context.Context, id string) error
}

type Repository struct {
	client *pgxpool.Pool
}

func NewRepository(client *pgxpool.Pool) *Repository {
	return &Repository{
		client: client,
	}
}

// TODO: transaction below

func (r *Repository) CreateUser(ctx context.Context, user *user.CreateUserDto) error {
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

	userQuery := `
		INSERT INTO public.users (email, password, full_name, telephone_number, city) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	var id string
	log.Printf("SQL query: %s", postgresql.FormatQuery(userQuery))
	err = tx.QueryRow(ctx, userQuery, user.Email, user.Password, user.FullName, user.TelephoneNumber, user.City).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL error: %s, Detail: %s, Where: %s, Code: %s, SQL State: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Println(newErr)
			return newErr
		}
		return err
	}

	rolesQuery := `
		INSERT INTO public.roles (user_id) 
		VALUES ($1)
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(rolesQuery))
	_, err = tx.Exec(ctx, rolesQuery, id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL error: %s, Detail: %s, Where: %s, Code: %s, SQL State: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Println(newErr)
			return newErr
		}
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetRole(ctx context.Context, id string) (string, error) {
	query := `
		SELECT role from roles WHERE user_id = $1
	`

	var role string
	err := r.client.QueryRow(ctx, query, id).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
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
		SELECT id, email, password, full_name, telephone_number, is_premium, city, rating, is_verified 
		FROM public.users
		WHERE id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	err = r.client.QueryRow(ctx, query, id).Scan(&u.Id, &u.Email, &u.Password, &u.FullName, &u.TelephoneNumber, &u.IsPremium, &u.City, &u.Rating, &u.IsVerified)
	if err != nil {
		log.Println(err)
		return u, err
	}

	return u, nil
}

func (r *Repository) GetOneUserByEmail(ctx context.Context, email string) (u user.GetUsersDto, err error) {
	query := `
		SELECT id, email, password, full_name, telephone_number, is_premium, city, rating, is_verified
		FROM public.users
		WHERE email = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	err = r.client.QueryRow(ctx, query, email).Scan(&u.Id, &u.Email, &u.Password, &u.FullName, &u.TelephoneNumber, &u.IsPremium, &u.City, &u.Rating, &u.IsVerified)
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

func (r *Repository) CreateRating(ctx context.Context, dto user.RateDto, userId, ratedBy string) error {
	query := `
		INSERT INTO public.users_ratings  (rate, comment, user_id, rate_by_user)
		VALUES ($1, $2, $3, $4)
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	exec, err := r.client.Exec(ctx, query, dto.Rating, dto.Comment, userId, ratedBy)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(exec.RowsAffected())

	return nil
}

func (r *Repository) GetUserRatings(ctx context.Context, userId string) (amount float32, sum float32, err error) {
	query := `
		SELECT count(*), sum(rate) FROM users_ratings WHERE user_id = $1
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	_ = r.client.QueryRow(ctx, query, userId).Scan(&amount, &sum)

	return amount, sum, nil
}

func (r *Repository) ChangeUserRating(ctx context.Context, id string, rating float32) error {
	query := `
		UPDATE public.users
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

func (r *Repository) GetAllUserRatings(ctx context.Context, id string) ([]user.GetAllRatingsDto, error) {
	query := `
		SELECT r.rate, r.comment, u.full_name, us.full_name
		FROM users_ratings r
		INNER JOIN users u on u.id = r.user_id
		INNER JOIN users us on us.id = r.rate_by_user
		WHERE r.user_id = $1;
	`

	log.Printf("SQL query: %s", postgresql.FormatQuery(query))
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ratings []user.GetAllRatingsDto
	for rows.Next() {
		var rate user.GetAllRatingsDto

		err = rows.Scan(&rate.Rating, &rate.Comment, &rate.User, &rate.RatedBy)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		ratings = append(ratings, rate)
	}

	return ratings, nil

}

func (r *Repository) CreateApplication(ctx context.Context, id string, filename string) error {
	query := `
		INSERT INTO applications (user_id, filename) VALUES ($1, $2)
	`

	exec, err := r.client.Exec(ctx, query, id, filename)
	if err != nil {
		return err
	}
	log.Println(exec.RowsAffected())
	return nil
}

func (r *Repository) GetApplications(ctx context.Context) ([]user.ApplicationDto, error) {
	query := `
		SELECT user_id, filename FROM applications WHERE is_visited = false
	`

	log.Println(query)
	rows, err := r.client.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var apps []user.ApplicationDto
	for rows.Next() {
		var a user.ApplicationDto
		err := rows.Scan(&a.UserId, &a.Filename)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		apps = append(apps, a)
	}

	return apps, nil
}

func (r *Repository) ChangeUserVerify(ctx context.Context, id string) error {
	query := `
		UPDATE users SET is_verified = true WHERE id = $1
	`

	log.Println(query)
	exec, err := r.client.Exec(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(exec.RowsAffected())

	query = `
		UPDATE applications SET is_visited = true WHERE user_id = $1
	`

	log.Println(query)
	exec, err = r.client.Exec(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(exec.RowsAffected())

	return nil
}

package images_storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Repository struct {
	client *pgxpool.Pool
}

func NewRepository(client *pgxpool.Pool) *Repository {
	return &Repository{client: client}
}

func (r *Repository) SaveImageToDB(ctx context.Context, url, carId string) error {
	query := `
		INSERT INTO public.car_images (url, car_id) 
		VALUES ($1, $2)
	`

	log.Println(query)
	_, err := r.client.Exec(ctx, query, url, carId)
	if err != nil {
		return err
	}

	return nil
}
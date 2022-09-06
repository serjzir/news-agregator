package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/serjzir/news-agregator/pkg/client/postgresql"
	"github.com/serjzir/news-agregator/pkg/logging"
)

type Repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *Repository) Create(ctx context.Context, posts []Post) error {
	for _, p := range posts {
		r.logger.Info("Вставка")
		q := "INSERT INTO news (title, content, pub_time, link) VALUES($1, $2, $3, $4) RETURNING id"
		r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
		if err := r.client.QueryRow(ctx, q, p.Title, p.Content, p.PubTime, p.Link).Scan(&p.ID); err != nil {
			var pgErr *pgconn.PgError
			if errors.Is(err, pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQL State: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return newErr
			}
			return err
		}
	}
	return nil
}

// News возвращает последние новости из БД.
func (r *Repository) News(c *gin.Context, id int) (p []Post, err error) {
	if id == 0 {
		id = 10
	}
	q := "SELECT id, title, content, pub_time, link FROM news ORDER BY pub_time DESC LIMIT $1"
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
	rows, err := r.client.Query(c, q, id)
	if err != nil {
		return nil, err
	}
	var news []Post
	for rows.Next() {
		var new Post
		err = rows.Scan(&new.ID, &new.Title, &new.Content, &new.PubTime, &new.Link)
		if err != nil {
			return nil, err
		}
		news = append(news, new)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return news, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) *Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}

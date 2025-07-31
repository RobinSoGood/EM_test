package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RobinSoGood/EM_test/internal/logger"
	"github.com/RobinSoGood/EM_test/internal/models"
	"github.com/RobinSoGood/EM_test/internal/storage/storageerror"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DBStorage struct {
	conn *pgx.Conn
}

func NewDB(ctx context.Context, addr string) (*DBStorage, error) {
	conn, err := pgx.Connect(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &DBStorage{conn: conn}, nil
}

func (dbs *DBStorage) Close() error {
	return dbs.conn.Close(context.Background())
}

func (dbs *DBStorage) GetSubs() ([]models.Sub, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := dbs.conn.Query(ctx, "SELECT * FROM subs")
	if err != nil {
		log.Error().Err(err).Msg("failed get data from table subs")
		return nil, err
	}
	var subs []models.Sub
	for rows.Next() {
		var sub models.Sub
		if err = rows.Scan(&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate); err != nil {
			log.Error().Err(err).Msg("failed scan rows data")
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (dbs *DBStorage) SaveSub(sub models.Sub) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var sid string
	row := dbs.conn.QueryRow(ctx, "SELECT sid FROM subs WHERE userID=$1 AND SeviceName=$2", sub.UserID, sub.ServiceName)
	err := row.Scan(&sid)
	if err == nil {
		return ``, storageerror.ErrSubAlredyExist
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return ``, err
	}
	id := uuid.New()
	sub.ID = id
	_, err = dbs.conn.Exec(ctx, "INSERT INTO subs (ID, userID, serviceName, price, startDate) VALUES ($1, $2, $3, $4, %5)",
		sub.ID, sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate)
	if err != nil {
		return ``, err
	}
	return sub.ID.String(), nil
}

func (dbs *DBStorage) GetSub(sid string) (models.Sub, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sub models.Sub
	row := dbs.conn.QueryRow(ctx, "SELECT * FROM subs WHERE bid=$1", sid)
	err := row.Scan(&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Sub{}, storageerror.ErrSubNoFound
		}
		return models.Sub{}, err
	}
	return sub, nil
}

func (dbs *DBStorage) SetDeleteSubStatus(sid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := dbs.conn.Exec(ctx, "DELETE FROM subs WHERE sid=$1", sid)
	if err != nil {
		return err
	}
	return nil
}

func Migrations(dbDsn string, migratePath string) error {
	log := logger.Get()
	migrPath := fmt.Sprintf("file://%s", migratePath)
	m, err := migrate.New(migrPath, dbDsn)
	if err != nil {
		log.Error().Err(err).Msg("failed to db conntect")
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Debug().Msg("no migratons apply")
			return nil
		}
		log.Error().Err(err).Msg("run migrations failed")
		return err
	}
	log.Debug().Msg("all migrations apply")
	return nil
}

func (dbs *DBStorage) GetTotalPriceByPeriod(req models.SumRequest) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var total int
	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subs
        WHERE ($1 = 0 OR user_id = $1)
          AND ($2 = '' OR service_name = $2)
          AND start_date <= $4
          AND (end_date >= $3 OR end_date IS NULL)
    `
	err := dbs.conn.QueryRow(ctx, query, req.UserID, req.ServiceName, req.Start, req.End).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

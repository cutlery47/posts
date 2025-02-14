package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cutlery47/posts/config"
	storage "github.com/cutlery47/posts/internal/storage/user-storage"
	"github.com/google/uuid"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type pgStorage struct {
	db *sql.DB

	conf config.UserStorage
}

func NewStorage(conf config.UserStorage) (*pgStorage, error) {
	dsn := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		conf.User,
		conf.Pass,
		conf.Host,
		conf.Port,
		conf.DB,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	err = db.PingContext(timeoutCtx)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish connection with postgres: %v", err)
	}
	log.Println("[SETUP] successfully established postgres connection!")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres.WithInstance: %v", err)
	}

	migrations := fmt.Sprintf("file://%v", conf.Migrations)
	m, err := migrate.NewWithDatabaseInstance(migrations, conf.DB, driver)
	if err != nil {
		return nil, fmt.Errorf("migrate.NewWithDatabaseInstance: %v", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("[SETUP] nothing to migrate")
		} else {
			return nil, fmt.Errorf("error when migrating: %v", err)
		}
	} else {
		log.Println("[SETUP] migrated successfully!")
	}

	return &pgStorage{
		db:   db,
		conf: conf,
	}, nil
}

func (pg *pgStorage) Register(ctx context.Context, in storage.InUser) (*storage.User, error) {
	var (
		user storage.User
	)

	row := pg.db.QueryRowContext(ctx, insertUserQuery, in.Name, in.Role)
	err := row.Scan(&user.Id, &user.Name, &user.Role, &user.CreatedAt)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return nil, storage.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

func (pg *pgStorage) Login(ctx context.Context, in storage.InUser) (*storage.Session, error) {
	var (
		id uuid.UUID
	)

	err := pg.db.QueryRowContext(ctx, getUserIdQuery, in.Name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}

	var (
		sesh storage.Session
	)

	expiresAt := time.Now().Add(pg.conf.SessionDuration)

	row := pg.db.QueryRowContext(ctx, insertSessionQuery, id, expiresAt)
	err = row.Scan(&sesh.Id, &sesh.UserId, &sesh.CreatedAt, &sesh.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &sesh, nil
}

func (pg *pgStorage) Logout(ctx context.Context, sesh storage.Session) error {
	_, err := pg.db.ExecContext(ctx, deleteSessionById, sesh.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrSessionNotFound
		}
		return err
	}

	return nil
}

func (pg *pgStorage) GetSession(ctx context.Context, id uuid.UUID) (*storage.Session, error) {
	var (
		sesh storage.Session
	)

	row := pg.db.QueryRowContext(ctx, getSessionById, id)
	err := row.Scan(&sesh.Id, &sesh.UserId, &sesh.CreatedAt, &sesh.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrSessionNotFound
		}
		return nil, err
	}

	return &sesh, nil
}

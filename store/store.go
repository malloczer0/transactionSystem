package store

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"time"
	"transactionSystemTestTask/logger"
	"transactionSystemTestTask/store/pg"
)

// Store contains all repositories
type Store struct {
	Pg *pg.DB // for KeepAlivePg

	Transaction TransactionRepo
	Client      ClientRepo
}

// New creates new store
func New(ctx context.Context) (*Store, error) {

	// connect to Postgres
	pgDB := func() *pg.DB {
		for {
			pgDB, err := pg.Dial()
			if err == nil {
				return pgDB
				//return nil, errors.Wrap(err, "pgdb.Dial failed")
			}
			log.Println("pgdb.Dial failed")
			t := time.Second * 3
			time.Sleep(t)
		}
	}()

	// Run Postgres migrations
	if pgDB != nil {
		log.Println("Running PostgreSQL migrations...")
		if err := runPgMigrations(); err != nil {
			return nil, errors.Wrap(err, "runPgMigrations failed")
		}
	}

	var store Store

	// Init Postgres repositories
	if pgDB != nil {
		store.Pg = pgDB
		go store.KeepAlivePg()
		store.Transaction = pg.NewTransactionRepo(pgDB)
		store.Client = pg.NewClientRepo(pgDB)
	}

	err := pg.CreateSchema(pgDB.DB)
	if err != nil {
		panic(err)
	}

	return &store, nil
}

// KeepAlivePollPeriod is a Pg keepalive check time period
const KeepAlivePollPeriod = 3

// KeepAlivePg makes sure PostgreSQL is alive and reconnects if needed
func (store *Store) KeepAlivePg() {

	log := logger.Get()
	var err error
	for {
		// Check if PostgreSQL is alive every 3 seconds
		time.Sleep(time.Second * KeepAlivePollPeriod)
		lostConnect := false
		if store.Pg == nil {
			lostConnect = true
		} else if _, err = store.Pg.Exec("SELECT 1"); err != nil {
			lostConnect = true
		}
		if !lostConnect {
			continue
		}
		log.Debug().Msg("[store.KeepAlivePg] Lost PostgreSQL connection. Restoring...")
		store.Pg, err = pg.Dial()
		if err != nil {
			log.Err(err)
			continue
		}
		log.Debug().Msg("[store.KeepAlivePg] PostgreSQL reconnected")
	}
}

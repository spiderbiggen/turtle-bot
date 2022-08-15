package migration

import (
	"context"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

type migration interface {
	Index() int
	Name() string
	// Up does the migration
	Up(context.Context, *sqlx.Tx) error
	// Down undoes the migration
	Down(context.Context, *sqlx.Tx) error
}

type completedMigration struct {
	Name      string    `db:"description"`
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

type migrationSlice []migration

func (m migrationSlice) FilterCompleted(migration []*completedMigration) migrationSlice {
	r := make(migrationSlice, 0, len(m)-len(migration))
Out:
	for _, m2 := range m {
		for _, cm := range migration {
			if m2.Name() == cm.Name {
				continue Out
			}
		}
		r = append(r, m2)
	}
	return r
}

func (m migrationSlice) Len() int           { return len(m) }
func (m migrationSlice) Less(i, j int) bool { return m[i].Index() < m[j].Index() }
func (m migrationSlice) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

// migration Registry
var migrationRegistry migrationSlice

// Register a new migration
func Register(m migration) {
	migrationRegistry = append(migrationRegistry, m)
}

// Up runs all migrationSlice in the registry
func Up(ctx context.Context, conn *sqlx.DB) error {
	if len(migrationRegistry) == 0 {
		log.Debugf("No migrations registered")
		return nil
	}

	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Make sure migration table exists
	_, err = tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS migrations (description VARCHAR PRIMARY KEY, id INTEGER NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}

	// Get completed migration
	var cm []*completedMigration
	err = tx.SelectContext(ctx, &cm, "SELECT description, id, created_at FROM migrations ORDER BY id")
	if err != nil {
		return err
	}

	p, err := tx.Preparex("INSERT INTO migrations (description, id, created_at) VALUES ($1, $2, DEFAULT)")
	if err != nil {
		log.Errorf("Failed to prepare migration update statement")
		return err
	}

	// Filter completed migration
	migrations := migrationRegistry.FilterCompleted(cm)

	if len(migrations) == 0 {
		last := cm[len(cm)-1]
		log.Debugf("Current version %d: %s", last.ID, last.Name)
		return nil
	}

	// Sort
	sort.Stable(migrations)
	for _, m := range migrations {
		log.Infof("Performing migration: %d: %s", m.Index(), m.Name())
		err := m.Up(ctx, tx)
		if err != nil {
			log.Errorf("Failed to migrate %d: %s. %v", m.Index(), m.Name(), err)
			err := tx.Rollback()
			if err != nil {
				log.Errorf("Failed to rollback migration. %v", err)
				return err
			}
			return err
		}
		_, err = p.ExecContext(ctx, m.Name(), m.Index())
		if err != nil {
			log.Errorf("Failed to save  %d: %s. %v", m.Index(), m.Name(), err)
			err := tx.Rollback()
			if err != nil {
				log.Errorf("Failed to rollback migration. %v", err)
				return err
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("Failed to commit migration. %v", err)
		return err
	}
	last := migrations[len(migrations)-1]
	log.Debugf("Current version %d: %s", last.Index(), last.Name())
	return nil
}

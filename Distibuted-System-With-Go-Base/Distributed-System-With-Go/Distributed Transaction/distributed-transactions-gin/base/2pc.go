package two_phase_commit

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Coordinator struct {
	DB *sqlx.DB
}

type Participant struct {
	DB *sqlx.DB
}

func (p *Participant) Prepare(ctx context.Context, tx *sqlx.Tx, prepareData interface{}) error {
	return nil
}

func (p *Participant) Commit(ctx context.Context, tx *sqlx.Tx) error {
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *Participant) RollBack(ctx context.Context, tx *sqlx.Tx) error {
	err := tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func (c *Coordinator) ExecuteDistributedTransaction(ctx context.Context, prepareData interface{}, participants []*Participant) error {
	var preparedTransactions []*sqlx.Tx
	for _, p := range participants {
		tx, err := p.DB.Beginx()
		if err != nil {
			return err
		}
		err = p.Prepare(ctx, tx, prepareData)
		if err != nil {
			return err
		}

		preparedTransactions = append(preparedTransactions, tx)
	}

	for i, p := range participants {
		err := p.Commit(ctx, preparedTransactions[i])
		if err != nil {
			for j := i; j > -1; j-- {
				_ = p.RollBack(ctx, preparedTransactions[j])
			}
			return err
		}
	}
	return nil
}

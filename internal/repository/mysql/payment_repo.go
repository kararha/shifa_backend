// repository/mysql/payment_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"

	// "time"

	"shifa/internal/repository"
)

type mysqlPaymentRepo struct {
	db *sql.DB
}

func NewMySQLPaymentRepo(db *sql.DB) repository.PaymentRepository {
	return &mysqlPaymentRepo{db: db}
}

func (r *mysqlPaymentRepo) Create(ctx context.Context, payment *repository.Payment) error {
	// Using named parameters with a transaction to better handle null values
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO payments (consultation_id, home_care_visit_id, amount, status, payment_date)
			  VALUES (?, ?, ?, ?, ?)`

	// Handle null foreign keys properly
	var consultationID, homeCareVisitID interface{}

	if payment.ConsultationID > 0 {
		consultationID = payment.ConsultationID
	} else {
		consultationID = nil
	}

	if payment.HomeCareVisitID > 0 {
		homeCareVisitID = payment.HomeCareVisitID
	} else {
		homeCareVisitID = nil
	}

	result, err := tx.ExecContext(ctx, query, consultationID, homeCareVisitID,
		payment.Amount, payment.Status, payment.PaymentDate)
	if err != nil {
		return err
	}

	// Get the ID of the newly inserted payment
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	payment.ID = int(id)

	return tx.Commit()
}

func (r *mysqlPaymentRepo) GetByID(ctx context.Context, id int) (*repository.Payment, error) {
	query := `SELECT id, consultation_id, home_care_visit_id, amount, status, payment_date, refund_date
			  FROM payments WHERE id = ?`

	var payment repository.Payment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID, &payment.ConsultationID, &payment.HomeCareVisitID,
		&payment.Amount, &payment.Status, &payment.PaymentDate, &payment.RefundDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

func (r *mysqlPaymentRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE payments SET status = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *mysqlPaymentRepo) GetByConsultationID(ctx context.Context, consultationID int) (*repository.Payment, error) {
	query := `SELECT id, consultation_id, home_care_visit_id, amount, status, payment_date, refund_date
			  FROM payments WHERE consultation_id = ?`

	var payment repository.Payment
	err := r.db.QueryRowContext(ctx, query, consultationID).Scan(
		&payment.ID, &payment.ConsultationID, &payment.HomeCareVisitID,
		&payment.Amount, &payment.Status, &payment.PaymentDate, &payment.RefundDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

func (r *mysqlPaymentRepo) GetByHomeCareVisitID(ctx context.Context, homeCareVisitID int) (*repository.Payment, error) {
	query := `SELECT id, consultation_id, home_care_visit_id, amount, status, payment_date, refund_date
			  FROM payments WHERE home_care_visit_id = ?`

	var payment repository.Payment
	err := r.db.QueryRowContext(ctx, query, homeCareVisitID).Scan(
		&payment.ID, &payment.ConsultationID, &payment.HomeCareVisitID,
		&payment.Amount, &payment.Status, &payment.PaymentDate, &payment.RefundDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

package models

import "time"

type Payment struct {
    ID              int        `json:"id"`
    Amount          float64    `json:"amount"`
    Status          string     `json:"status"`
    PaymentDate     *time.Time `json:"payment_date"`
    RefundDate      *time.Time `json:"refund_date"`
    ConsultationID  int        `json:"consultation_id"`
    HomeCareVisitID int        `json:"home_care_visit_id"`
}




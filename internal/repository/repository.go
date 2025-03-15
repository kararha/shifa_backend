// repository/repository.go

package repository

import (
	"context"
	"shifa/internal/models"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
}

type DoctorRepository interface {
	Create(ctx context.Context, doctor *models.Doctor) error
	GetByID(ctx context.Context, id int) (*models.Doctor, error)
	GetByUserID(ctx context.Context, userID int) (*models.Doctor, error)
	Update(ctx context.Context, doctor *models.Doctor) error
	Delete(ctx context.Context, id int) error
	GetBySpecialty(ctx context.Context, specialty string, limit, offset int) ([]*models.Doctor, error)
	GetByServiceType(ctx context.Context, serviceTypeID int, limit, offset int) ([]*models.Doctor, error)
}

type PatientRepository interface {
	Create(ctx context.Context, patient *models.Patient) error
	GetByUserID(ctx context.Context, userID int) (*models.Patient, error)
	Update(ctx context.Context, patient *models.Patient) error
	Delete(ctx context.Context, userID int) error
	List(ctx context.Context, offset, limit int) ([]*models.Patient, error)
}

type HomeCareProviderRepository interface {
	GetAll(ctx context.Context) ([]models.HomeCareProvider, error) // Add this
	Create(ctx context.Context, provider *models.HomeCareProvider) error
	// GetByID(ctx context.Context, id int) (*HomeCareProvider, error)
	GetByID(ctx context.Context, id int) (*models.HomeCareProvider, error)
	GetByUserID(ctx context.Context, userID int) (*models.HomeCareProvider, error)
	Update(ctx context.Context, provider *models.HomeCareProvider) error // Ensure this matches
	Delete(ctx context.Context, id int) error
	GetByServiceType(ctx context.Context, serviceTypeID int, limit, offset int) ([]*models.HomeCareProvider, error)
	Search(ctx context.Context, query string) ([]*models.HomeCareProvider, error)
}

type ServiceTypeRepository interface {
	Create(ctx context.Context, serviceType *ServiceType) error
	GetByID(ctx context.Context, id int) (*ServiceType, error)
	GetAll(ctx context.Context) ([]*ServiceType, error)
	Update(ctx context.Context, serviceType *ServiceType) error
	Delete(ctx context.Context, id int) error
}

type DoctorAvailabilityRepository interface {
	Create(ctx context.Context, availability *models.DoctorAvailability) error
	GetByDoctorID(ctx context.Context, doctorID int) ([]*models.DoctorAvailability, error)
	Update(ctx context.Context, availability *models.DoctorAvailability) error
	Delete(ctx context.Context, id int) error
	ListByDoctorID(ctx context.Context, doctorID int) ([]*models.DoctorAvailability, error)
	ListAllAvailability(ctx context.Context) ([]*models.DoctorAvailability, error)
}

type MedicalHistoryRepository interface {
	Create(ctx context.Context, history *models.MedicalHistory) error
	GetByPatientID(ctx context.Context, patientID int) ([]*models.MedicalHistory, error)
	Update(ctx context.Context, history *models.MedicalHistory) error
	Delete(ctx context.Context, id int) error
}

// AppointmentFilter struct should be defined in the repository package
type AppointmentFilter struct {
	StartDate    *time.Time
	EndDate      *time.Time
	Status       string
	ProviderType string
}

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *models.Appointment) error
	GetByID(ctx context.Context, id int) (*models.Appointment, error)
	Update(ctx context.Context, appointment *models.Appointment) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter AppointmentFilter, offset, limit int) ([]*models.Appointment, error)
	GetByProviderID(ctx context.Context, providerID int, providerType string) ([]*models.Appointment, error)
	GetByPatientID(ctx context.Context, patientID, limit, offset int) ([]*models.Appointment, error)
}

type ConsultationRepository interface {
	// Create inserts a new consultation into the database
	Create(ctx context.Context, consultation *models.Consultation) error

	// GetByID retrieves a consultation by its ID
	GetByID(ctx context.Context, id int) (*models.Consultation, error)

	// GetByAppointmentID retrieves a consultation by its appointment ID
	GetByAppointmentID(ctx context.Context, appointmentID int) (*models.Consultation, error)

	// Update modifies an existing consultation
	Update(ctx context.Context, consultation *models.Consultation) error

	// Delete removes a consultation from the database
	Delete(ctx context.Context, id int) error

	// List retrieves consultations based on filter criteria with pagination
	List(ctx context.Context, filter models.ConsultationFilter, offset, limit int) ([]*models.Consultation, error)
}

type HomeCareVisitRepository interface {
	Create(ctx context.Context, visit *models.HomeCareVisit) error
	Update(ctx context.Context, visit *models.HomeCareVisit) error
	GetByID(ctx context.Context, id int) (*models.HomeCareVisit, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter models.HomeCareVisitFilter) ([]models.HomeCareVisit, error)
	GetByPatientID(ctx context.Context, patientID int) ([]models.HomeCareVisit, error)
	GetByProviderID(ctx context.Context, providerID int) ([]models.HomeCareVisit, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.HomeCareVisit, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id int) (*Payment, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	GetByConsultationID(ctx context.Context, consultationID int) (*Payment, error)
	GetByHomeCareVisitID(ctx context.Context, homeCareVisitID int) (*Payment, error)
	// DeletePayment(ctx context.Context, id int) error
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	GetByID(ctx context.Context, id int) (*Review, error)
	GetByDoctorID(ctx context.Context, doctorID int, limit, offset int) ([]*Review, error)
	GetByHomeCareProviderID(ctx context.Context, providerID int, limit, offset int) ([]*Review, error)
	GetReviewByID(ctx context.Context, id int) (*models.Review, error) // Add this method
	GetReviewsByDoctorID(ctx context.Context, doctorID int) ([]models.Review, error)
	UpdateReview(ctx context.Context, review *models.Review) error // Add this method
	DeleteReview(ctx context.Context, id int) error                // Add this method
}

type ChatMessageRepository interface {
	Create(ctx context.Context, message *ChatMessage) error
	GetByConsultationID(ctx context.Context, consultationID int, limit, offset int) ([]*ChatMessage, error)
	MarkAsRead(ctx context.Context, messageID int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetByUserID(ctx context.Context, userID int, limit, offset int) ([]*models.Notification, error)
	MarkAsRead(ctx context.Context, notificationID int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
}

// Structs representing database entities
type User struct {
	ID           int
	Email        string
	PasswordHash string
	Name         string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Doctor struct {
	UserID            int
	Specialty         string
	ServiceTypeID     int
	LicenseNumber     string
	ExperienceYears   int
	Qualifications    string
	Achievements      string
	Bio               string
	ProfilePictureURL string
	ConsultationFee   float64
	Rating            float64
	IsVerified        bool
	IsAvailable       bool
	Status            string
	Latitude          float64
	Longitude         float64
	CreatedAt         time.Time // Add CreatedAt field
	UpdatedAt         time.Time // Add UpdatedAt field
}

type Patient struct {
	UserID                int
	DateOfBirth           time.Time
	Gender                string
	Phone                 string
	Address               string
	EmergencyContactName  string
	EmergencyContactPhone string
}

type HomeCareProvider struct {
	UserID            int
	ServiceTypeID     int
	ExperienceYears   int
	Qualifications    string
	Bio               string
	ProfilePictureURL string
	HourlyRate        float64
	Rating            float64
	IsVerified        bool
	IsAvailable       bool
	Status            string
	Latitude          float64
	Longitude         float64
}

type ServiceType struct {
	ID          int
	Name        string
	Description string
	IsHomeCare  bool
}

type DoctorAvailability struct {
	ID        int
	DoctorID  int
	DayOfWeek int
	StartTime time.Time
	EndTime   time.Time
}

type MedicalHistory struct {
	ID            int
	PatientID     int
	ConditionName string
	DiagnosisDate time.Time
	Treatment     string
	IsCurrent     bool
}

// repository/repository.go
type Appointment struct {
	ID                 int
	PatientID          int
	ProviderType       string
	DoctorID           *int
	HomeCareProviderID *int
	AppointmentDate    time.Time
	StartTime          time.Time
	EndTime            time.Time
	Status             string
	CancellationReason *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Consultation struct {
	ID               int       `json:"id" db:"id"`
	PatientID        int       `json:"patient_id"`
	DoctorID         int       `json:"doctor_id"`
	AppointmentID    int       `json:"appointment_id" db:"appointment_id"`
	ConsultationType string    `json:"consultation_type" db:"consultation_type"`
	Status           string    `json:"status" db:"status"`
	StartedAt        time.Time `json:"started_at" db:"started_at"`
	CompletedAt      time.Time `json:"completed_at" db:"completed_at"`
	Fee              float64   `json:"fee" db:"fee"`
}

type HomeCareVisit struct {
	ID                  int
	AppointmentID       int
	Address             string
	Latitude            float64
	Longitude           float64
	DurationHours       float64
	SpecialRequirements string
	Status              string
}

type Payment struct {
	ID              int
	Amount          float64
	Status          string
	PaymentDate     *time.Time
	RefundDate      *time.Time
	ConsultationID  int
	HomeCareVisitID int
}

type Review struct {
	ID                 int       `json:"id" db:"id"`
	PatientID          int       `json:"patient_id" db:"patient_id"`
	ReviewType         string    `json:"review_type" db:"review_type"`
	ConsultationID     *int      `json:"consultation_id,omitempty" db:"consultation_id"`
	HomeCareVisitID    *int      `json:"home_care_visit_id,omitempty" db:"home_care_visit_id"`
	DoctorID           *int      `json:"doctor_id,omitempty" db:"doctor_id"`
	HomeCareProviderID *int      `json:"home_care_provider_id,omitempty" db:"home_care_provider_id"`
	Rating             int       `json:"rating" db:"rating"`
	Comment            string    `json:"comment" db:"comment"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type ChatMessage struct {
	ID             int
	ConsultationID int
	SenderType     string
	SenderID       int
	Message        string
	SentAt         time.Time
	IsRead         bool
}

type Notification struct {
	ID               int
	UserID           int
	NotificationType string
	Message          string
	IsRead           bool
	CreatedAt        time.Time
}

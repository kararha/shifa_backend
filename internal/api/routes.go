package api

import (
	"database/sql"
	"net/http"
	"shifa/internal/api/handlers"
	"shifa/internal/api/middleware"
	"shifa/internal/repository/mysql"
	"shifa/internal/service"
	"shifa/pkg/fileutils"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// NewRouter creates and configures a new router with all application routes
func NewRouter(db *sql.DB, log *logrus.Logger, jwtSecret string) *mux.Router {
	router := mux.NewRouter()

	// Serve static files from the uploads directory
	uploadsFS := http.FileServer(http.Dir(fileutils.UploadDir))
	router.PathPrefix("/uploads/").Handler(
		http.StripPrefix("/uploads/", uploadsFS),
	)

	// Create an API subrouter
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Initialize repositories
	appointmentRepo := mysql.NewAppointmentRepository(db)
	userRepo := mysql.NewUserRepo(db)
	doctorRepo := mysql.NewDoctorRepo(db)
	serviceTypeRepo := mysql.NewServiceTypeRepository(db)
	patientRepo := mysql.NewPatientRepo(db)
	consultationRepo := mysql.NewConsultationRepo(db)
	reviewRepo := mysql.NewReviewRepository(db)
	homeCareProviderRepo := mysql.NewHomeCareProviderRepo(db) // Fixed function name
	medicalHistoryRepo := mysql.NewMedicalHistoryRepo(db)     // Add this line
	chatMessageRepo := mysql.NewMySQLChatMessageRepo(db)
	paymentRepo := mysql.NewMySQLPaymentRepo(db)
	notificationRepo := mysql.NewNotificationRepo(db, log)
	homeCareVisitRepo := mysql.NewHomeVisitRepo(db)
	systemLogRepo := mysql.NewSystemLogRepo(db)                   // Create SystemLogRepo first
	doctorAvailabilityRepo := mysql.NewDoctorAvailabilityRepo(db) // Add this line

	// Initialize services
	appointmentService := service.NewAppointmentService(
		appointmentRepo,
		doctorRepo,
		homeCareProviderRepo,
		log,
	)
	userService := service.NewUserService(userRepo)
	doctorService := service.NewDoctorService(doctorRepo, log)
	serviceTypeService := service.NewServiceTypeService(serviceTypeRepo, log)
	patientService := service.NewPatientService(patientRepo, log)
	consultationService := service.NewConsultationService(consultationRepo, log)
	reviewService := service.NewReviewService(reviewRepo, log)
	homeCareProviderService := service.NewHomeCareProviderService(homeCareProviderRepo, log)
	medicalHistoryService := service.NewMedicalHistoryService(medicalHistoryRepo, log) // Add this line
	chatMessageService := service.NewChatService(chatMessageRepo, log)                 // Where logger is an instance of your custom logger
	paymentService := service.NewPaymentService(paymentRepo, log)
	notificationService := service.NewNotificationService(notificationRepo, log)
	homeCareVisitService := service.NewHomeCareVisitService(homeCareVisitRepo, log)
	authService := service.NewAuthService(userRepo, jwtSecret)
	systemLogService := service.NewSystemLogService(systemLogRepo)                                 // Pass systemLogRepo to NewSystemLogService
	doctorAvailabilityService := service.NewDoctorAvailabilityService(doctorAvailabilityRepo, log) // Add this line

	// Initialize handlers
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)
	userHandler := handlers.NewUserHandler(userService)
	doctorHandler := handlers.NewDoctorHandler(doctorService)
	serviceTypeHandler := handlers.NewServiceTypeHandler(serviceTypeService)
	patientHandler := handlers.NewPatientHandler(patientService)
	consultationHandler := handlers.NewConsultationHandler(consultationService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	homeCareProviderHandler := handlers.NewHomeCareProviderHandler(homeCareProviderService) // Add this line
	medicalHistoryHandler := handlers.NewMedicalHistoryHandler(medicalHistoryService)       // Add this line
	chatMessageHandler := handlers.NewChatMessageHandler(chatMessageService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	homeCareVisitHandler := handlers.NewHomeCareVisitHandler(homeCareVisitService, log)
	authHandler := handlers.NewAuthHandler(authService)
	doctorAvailabilityHandler := handlers.NewDoctorAvailabilityHandler(doctorAvailabilityService) // Add this line

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)
	logMiddleware := middleware.NewSystemLogMiddleware(systemLogService)

	// Register routes
	registerAppointmentRoutes(apiRouter, appointmentHandler)
	registerUserRoutes(apiRouter, userHandler)
	registerDoctorRoutes(apiRouter, doctorHandler)
	registerServiceTypeRoutes(apiRouter, serviceTypeHandler)
	registerPatientRoutes(apiRouter, patientHandler)
	registerConsultationRoutes(apiRouter, consultationHandler)
	registerReviewRoutes(apiRouter, reviewHandler)
	registerHomeCareProviderRoutes(apiRouter, homeCareProviderHandler) // Add this line
	registerMedicalHistoryRoutes(apiRouter, medicalHistoryHandler)     // Add this line
	registerChatMessageRoutes(apiRouter, chatMessageHandler)
	registerPaymentRoutes(apiRouter, paymentHandler)
	registerNotificationRoutes(apiRouter, notificationHandler)
	registerHomeCareVisitRoutes(apiRouter, homeCareVisitHandler)
	registerDoctorAvailabilityRoutes(apiRouter, doctorAvailabilityHandler) // Add this line
	// Register public routes (no auth required)
	registerAuthRoutes(apiRouter, authHandler)

	// Create protected subrouter
	protected := apiRouter.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.RequireAuth)
	protected.Use(logMiddleware.LogSystemAction)
	// // Register all protected routes
	// registerProtectedRoutes(protected, userHandler, appointmentHandler, doctorHandler,
	//     serviceTypeHandler, patientHandler, consultationHandler, reviewHandler,
	//     homeCareProviderHandler, medicalHistoryHandler, chatMessageHandler,
	//     paymentHandler, notificationHandler, homeCareVisitHandler)

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request for path: %s", r.URL.Path)
		http.Error(w, "404 page not found", http.StatusNotFound)
	})

	return router
}

func registerProtectedRoutes(protected *mux.Router, userHandler *handlers.UserHandler, appointmentHandler *handlers.AppointmentHandler, doctorHandler *handlers.DoctorHandler, serviceTypeHandler *handlers.ServiceTypeHandler, patientHandler *handlers.PatientHandler, consultationHandler *handlers.ConsultationHandler, reviewHandler *handlers.ReviewHandler, homeCareProviderHandler *handlers.HomeCareProviderHandler, medicalHistoryHandler *handlers.MedicalHistoryHandler, chatMessageHandler *handlers.ChatMessageHandler, paymentHandler *handlers.PaymentHandler, notificationHandler *handlers.NotificationHandler, homeCareVisitHandler *handlers.HomeCareVisitHandler) {
	panic("unimplemented")
}

// registerAppointmentRoutes sets up all appointment-related routes
func registerAppointmentRoutes(router *mux.Router, handler *handlers.AppointmentHandler) {
	appointmentRouter := router.PathPrefix("/appointments").Subrouter()

	// List/Search appointments with query parameters
	appointmentRouter.HandleFunc("", handler.GetAppointmentsByProvider).
		Methods("GET").
		Queries("type", "{type}", "providerId", "{providerId}")

	// Other routes
	appointmentRouter.HandleFunc("", handler.CreateAppointment).Methods("POST")
	appointmentRouter.HandleFunc("/{id}", handler.GetAppointment).Methods("GET")
	appointmentRouter.HandleFunc("/{id}", handler.UpdateAppointment).Methods("PUT")
	appointmentRouter.HandleFunc("/{id}", handler.DeleteAppointment).Methods("DELETE")
}

// registerUserRoutes sets up all user-related routes
func registerUserRoutes(router *mux.Router, handler *handlers.UserHandler) {
	// Note: Change from router.HandleFunc to userRouter.HandleFunc
	userRouter := router.PathPrefix("/users").Subrouter()

	userRouter.HandleFunc("", handler.ListUsers).Methods("GET")          // GET /api/users
	userRouter.HandleFunc("", handler.CreateUser).Methods("POST")        // POST /api/users
	userRouter.HandleFunc("/{id}", handler.GetUser).Methods("GET")       // GET /api/users/{id}
	userRouter.HandleFunc("/{id}", handler.UpdateUser).Methods("PUT")    // PUT /api/users/{id}
	userRouter.HandleFunc("/{id}", handler.DeleteUser).Methods("DELETE") // DELETE /api/users/{id}

	// Auth routes remain at root level
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/login", handler.Login).Methods("POST")
}

// registerDoctorRoutes sets up all doctor-related routes
func registerDoctorRoutes(router *mux.Router, handler *handlers.DoctorHandler) {
	doctorRouter := router.PathPrefix("/doctors").Subrouter()

	doctorRouter.HandleFunc("", handler.CreateDoctor).Methods("POST")
	doctorRouter.HandleFunc("", handler.ListDoctors).Methods("GET")
	doctorRouter.HandleFunc("/{id}", handler.GetDoctor).Methods("GET")
	doctorRouter.HandleFunc("/{id}", handler.UpdateDoctor).Methods("PUT")
	// Add the new search endpoint
	doctorRouter.HandleFunc("/search", handler.SearchDoctors).Methods("GET")
}

// registerServiceTypeRoutes sets up all service type-related routes
func registerServiceTypeRoutes(router *mux.Router, handler *handlers.ServiceTypeHandler) {
	serviceTypeRouter := router.PathPrefix("/service-types").Subrouter()

	serviceTypeRouter.HandleFunc("", handler.CreateServiceType).Methods("POST")
	serviceTypeRouter.HandleFunc("", handler.ListServiceTypes).Methods("GET")
	serviceTypeRouter.HandleFunc("/{id}", handler.GetServiceType).Methods("GET")
	serviceTypeRouter.HandleFunc("/{id}", handler.UpdateServiceType).Methods("PUT")
	serviceTypeRouter.HandleFunc("/{id}", handler.DeleteServiceType).Methods("DELETE")
}

// registerPatientRoutes sets up all patient-related routes
func registerPatientRoutes(router *mux.Router, handler *handlers.PatientHandler) {
	patientRouter := router.PathPrefix("/patients").Subrouter()

	patientRouter.HandleFunc("", handler.CreatePatient).Methods("POST")
	patientRouter.HandleFunc("", handler.ListPatients).Methods("GET")
	patientRouter.HandleFunc("/{id}", handler.GetPatient).Methods("GET")
	patientRouter.HandleFunc("/{id}", handler.UpdatePatient).Methods("PUT")
	patientRouter.HandleFunc("/{id}", handler.DeletePatient).Methods("DELETE")
}

// registerConsultationRoutes sets up all consultation-related routes
func registerConsultationRoutes(router *mux.Router, handler *handlers.ConsultationHandler) {
	consultationRouter := router.PathPrefix("/consultations").Subrouter()

	consultationRouter.HandleFunc("", handler.StartConsultation).Methods("POST")
	consultationRouter.HandleFunc("/{id}/complete", handler.CompleteConsultation).Methods("PUT")
	consultationRouter.HandleFunc("/{id}", handler.GetConsultation).Methods("GET")
	consultationRouter.HandleFunc("", handler.ListConsultations).Methods("GET")
	consultationRouter.HandleFunc("/{id}", handler.UpdateConsultation).Methods("PUT")
	consultationRouter.HandleFunc("/{id}", handler.DeleteConsultation).Methods("DELETE")
}

// Add a new function to register review routes
func registerReviewRoutes(router *mux.Router, handler *handlers.ReviewHandler) {
	reviewRouter := router.PathPrefix("/reviews").Subrouter()

	// Core CRUD operations
	reviewRouter.HandleFunc("", handler.CreateReview).Methods("POST")
	reviewRouter.HandleFunc("/{id}", handler.GetReview).Methods("GET")
	reviewRouter.HandleFunc("", handler.ListReviews).Methods("GET")
	reviewRouter.HandleFunc("/{id}", handler.UpdateReview).Methods("PUT")
	reviewRouter.HandleFunc("/{id}", handler.DeleteReview).Methods("DELETE")

	// Doctor reviews - using doctor_id query parameter
	reviewRouter.HandleFunc("/doctor/{doctorId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		doctorID := vars["doctorId"]
		// Create a new query with the doctor_id parameter
		q := r.URL.Query()
		q.Set("doctor_id", doctorID)
		r.URL.RawQuery = q.Encode()
		handler.ListReviews(w, r)
	}).Methods("GET")

	// Provider reviews - update to match repository implementation
	reviewRouter.HandleFunc("/provider/{providerId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerID := vars["providerId"]
		// Create a new query with the home_care_provider_id parameter
		q := r.URL.Query()
		q.Set("home_care_provider_id", providerID)
		r.URL.RawQuery = q.Encode()
		handler.ListReviews(w, r)
	}).Methods("GET")
}

// Add this new function
func registerHomeCareProviderRoutes(router *mux.Router, handler *handlers.HomeCareProviderHandler) {
	providerRouter := router.PathPrefix("/providers").Subrouter()

	providerRouter.HandleFunc("", handler.CreateHomeCareProvider).Methods("POST")
	providerRouter.HandleFunc("", handler.ListHomeCareProviders).Methods("GET")
	providerRouter.HandleFunc("/{id}", handler.GetHomeCareProvider).Methods("GET")
	providerRouter.HandleFunc("/{id}", handler.UpdateHomeCareProvider).Methods("PUT")
	providerRouter.HandleFunc("/{id}", handler.DeleteHomeCareProvider).Methods("DELETE")
	providerRouter.HandleFunc("/user/{user_id}", handler.GetHomeCareProviderByUserID).Methods("GET")
	providerRouter.HandleFunc("/search", handler.SearchHomeCareProviders).Methods("GET")
}

// Add this new function to register medical history routes
func registerMedicalHistoryRoutes(router *mux.Router, handler *handlers.MedicalHistoryHandler) {
	medicalHistoryRouter := router.PathPrefix("/medical-histories").Subrouter()
	// Create a new medical history
	medicalHistoryRouter.HandleFunc("", handler.CreateMedicalHistory).Methods("POST")
	// Get medical histories by patient ID
	medicalHistoryRouter.HandleFunc("", handler.GetMedicalHistory).Methods("GET")
	// Update a medical history
	medicalHistoryRouter.HandleFunc("/{id}", handler.UpdateMedicalHistory).Methods("PUT")
	// Delete a medical history
	medicalHistoryRouter.HandleFunc("/{id}", handler.DeleteMedicalHistory).Methods("DELETE")
}

// registerChatMessageRoutes sets up all chat message-related routes
func registerChatMessageRoutes(router *mux.Router, handler *handlers.ChatMessageHandler) {
	chatRouter := router.PathPrefix("/chat").Subrouter()

	// Route to send a new message
	chatRouter.HandleFunc("/messages", handler.SendMessage).Methods("POST")

	// Route to get messages for a specific consultation
	chatRouter.HandleFunc("/messages", handler.GetMessagesByConsultation).Methods("GET")

	// Route to mark a message as read
	chatRouter.HandleFunc("/messages/{id}/read", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		messageID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}
		handler.MarkMessageAsRead(w, r, messageID)
	}).Methods("PUT")

	// Route to get unread message count
	chatRouter.HandleFunc("/unread-count", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		handler.GetUnreadMessageCount(w, r, userID)
	}).Methods("GET")
}

// registerPaymentRoutes sets up all payment-related routes
// In your router configuration file
func registerPaymentRoutes(router *mux.Router, handler *handlers.PaymentHandler) {
	paymentRouter := router.PathPrefix("/payments").Subrouter()

	// Create a new payment
	paymentRouter.HandleFunc("", handler.CreatePayment).Methods("POST")

	// Get a specific payment by ID
	paymentRouter.HandleFunc("/{id}", handler.GetPayment).Methods("GET")

	// Update payment status
	paymentRouter.HandleFunc("/{id}", handler.UpdatePayment).Methods("PUT")

	// Get payment by consultation ID
	paymentRouter.HandleFunc("/consultation/{consultationId}", handler.GetPaymentByConsultationID).Methods("GET")

	// Get payment by home care visit ID
	paymentRouter.HandleFunc("/home-care-visit/{visitId}", handler.GetPaymentByHomeCareVisitID).Methods("GET")

	// Process refund
	paymentRouter.HandleFunc("/{id}/refund", handler.ProcessRefund).Methods("POST")
}

func registerNotificationRoutes(router *mux.Router, handler *handlers.NotificationHandler) {
	notificationRouter := router.PathPrefix("/notifications").Subrouter()

	// Create a new notification
	notificationRouter.HandleFunc("", handler.CreateNotification).Methods("POST")

	// Get notifications for a user with pagination
	notificationRouter.HandleFunc("/user/{userId}", handler.GetUserNotifications).Methods("GET")

	// Mark a notification as read
	notificationRouter.HandleFunc("/{id}/read", handler.MarkNotificationAsRead).Methods("PUT")

	// Get unread notification count for a user
	notificationRouter.HandleFunc("/unread-count", handler.GetUnreadCount).Methods("GET")

	// Send appointment reminder notification
	notificationRouter.HandleFunc("/appointment-reminder/{appointmentId}", handler.SendAppointmentReminder).Methods("POST")
}

// Add the new function for registering home care visit routes
func registerHomeCareVisitRoutes(router *mux.Router, handler *handlers.HomeCareVisitHandler) {
	visitRouter := router.PathPrefix("/home-care-visits").Subrouter()

	// Core CRUD operations
	visitRouter.HandleFunc("", handler.ScheduleHomeCareVisit).Methods("POST")
	visitRouter.HandleFunc("", handler.ListHomeCareVisits).Methods("GET")
	visitRouter.HandleFunc("/{id}", handler.GetHomeCareVisit).Methods("GET")
	visitRouter.HandleFunc("/{id}", handler.UpdateHomeCareVisit).Methods("PUT")
	visitRouter.HandleFunc("/{id}", handler.DeleteHomeCareVisit).Methods("DELETE")

	// Additional operations
	// visitRouter.HandleFunc("/{id}/cancel", handler.CancelHomeCareVisit).Methods("PUT")
	// visitRouter.HandleFunc("/patient/{patientId}", handler.GetHomeCareVisitsByPatient).Methods("GET")
	// visitRouter.HandleFunc("/provider/{providerId}", handler.GetHomeCareVisitsByProvider).Methods("GET")
}

// New function to register auth routes
func registerAuthRoutes(router *mux.Router, handler *handlers.AuthHandler) {
	router.HandleFunc("/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/auth/login", handler.Login).Methods("POST")
	// router.HandleFunc("/auth/refresh-token", handler.RefreshToken).Methods("POST")
	// router.HandleFunc("/auth/forgot-password", handler.ForgotPassword).Methods("POST")
	// router.HandleFunc("/auth/reset-password", handler.ResetPassword).Methods("POST")
}

// Add this new function to register doctor availability routes
func registerDoctorAvailabilityRoutes(router *mux.Router, handler *handlers.DoctorAvailabilityHandler) {
	router.HandleFunc("/doctors/{doctorId}/availability", handler.SetAvailability).Methods("POST")
	router.HandleFunc("/doctors/{doctorId}/availability", handler.GetAvailability).Methods("GET")
	router.HandleFunc("/doctors/{doctorId}/availability/{id}", handler.UpdateAvailability).Methods("PUT")
	router.HandleFunc("/doctors/{doctorId}/availability/{id}", handler.DeleteAvailability).Methods("DELETE")
}

// Add new route for listing all availability slots
func registerAvailabilityRoutes(router *mux.Router, handler *handlers.DoctorAvailabilityHandler) {
	router.HandleFunc("/availability", handler.ListAllAvailability).Methods("GET")
}
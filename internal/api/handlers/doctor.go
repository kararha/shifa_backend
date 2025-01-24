package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "shifa/internal/models"
    "shifa/internal/service"
    "shifa/pkg/fileutils"
)

type DoctorHandler struct {
    doctorService *service.DoctorService
}

func NewDoctorHandler(doctorService *service.DoctorService) *DoctorHandler {
    return &DoctorHandler{doctorService: doctorService}
}

// CreateDoctor handles both JSON and multipart form data for doctor creation
func (h *DoctorHandler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
    var doctor models.Doctor
    var err error

    contentType := r.Header.Get("Content-Type")
    if contentType != "" && len(contentType) >= 19 && contentType[:19] == "multipart/form-data" {
        // Handle multipart form data with file upload
        if err := r.ParseMultipartForm(10 << 20); err != nil {
            http.Error(w, "Failed to parse form", http.StatusBadRequest)
            return
        }

        // Get the file from form
        file, header, err := r.FormFile("profile_picture")
        if err != nil && err != http.ErrMissingFile {
            http.Error(w, "Error retrieving file", http.StatusBadRequest)
            return
        }

        var profilePicturePath string
        if file != nil {
            defer file.Close()
            
            // Save the file and get its relative path
            profilePicturePath, err = fileutils.SaveFile(file, header.Filename, fileutils.ProfileImages)
            if err != nil {
                http.Error(w, "Failed to save profile picture", http.StatusInternalServerError)
                return
            }
        }

        // Create doctor object from form data
        doctor = models.Doctor{
            UserID:            parseInt(r.FormValue("user_id")),
            Specialty:         r.FormValue("specialty"),
            ServiceTypeID:     parseInt(r.FormValue("service_type_id")),
            LicenseNumber:     r.FormValue("license_number"),
            ExperienceYears:   parseInt(r.FormValue("experience_years")),
            Qualifications:    r.FormValue("qualifications"),
            Achievements:      r.FormValue("achievements"),
            Bio:              r.FormValue("bio"),
            ProfilePictureURL: profilePicturePath,
            ConsultationFee:  parseFloat(r.FormValue("consultation_fee")),
        }
    } else {
        // Handle JSON request
        if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    }

    // Register the doctor using the service
    err = h.doctorService.RegisterDoctor(doctor)
    if err != nil {
        // If creation fails and we uploaded a file, delete it
        if doctor.ProfilePictureURL != "" {
            fileutils.DeleteFile(doctor.ProfilePictureURL)
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) GetDoctor(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    doctorID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
        return
    }

    doctor, err := h.doctorService.GetByID(doctorID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    doctorID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
        return
    }

    var doctor models.Doctor
    if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    doctor.UserID = doctorID
    err = h.doctorService.UpdateDoctor(doctor)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) ListDoctors(w http.ResponseWriter, r *http.Request) {
    specialty := r.URL.Query().Get("specialty")
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

    if limit == 0 {
        limit = 10 // default limit
    }

    doctors, err := h.doctorService.ListDoctorsBySpecialty(specialty, limit, offset)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(doctors)
}

// Add this new method to the DoctorHandler struct
func (h *DoctorHandler) SearchDoctors(w http.ResponseWriter, r *http.Request) {
    var params models.DoctorSearchParams

    // Parse query parameters
    query := r.URL.Query()
    params.SearchTerm = query.Get("search")
    params.Specialty = query.Get("specialty")

    if rating := query.Get("min_rating"); rating != "" {
        if ratingVal, err := strconv.ParseFloat(rating, 64); err == nil {
            params.MinRating = &ratingVal
        }
    }

    if lat := query.Get("lat"); lat != "" {
        if latVal, err := strconv.ParseFloat(lat, 64); err == nil {
            params.LocationLat = &latVal
        }
    }

    if lng := query.Get("lng"); lng != "" {
        if lngVal, err := strconv.ParseFloat(lng, 64); err == nil {
            params.LocationLng = &lngVal
        }
    }

    if radius := query.Get("radius"); radius != "" {
        if radiusVal, err := strconv.Atoi(radius); err == nil {
            params.RadiusKm = &radiusVal
        }
    }

    if serviceType := query.Get("service_type_id"); serviceType != "" {
        if serviceTypeVal, err := strconv.Atoi(serviceType); err == nil {
            params.ServiceTypeID = &serviceTypeVal
        }
    }

    // Call the service
    results, err := h.doctorService.SearchDoctors(params)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

// Helper functions for parsing form values
// func parseInt(s string) int {
//     i, _ := strconv.Atoi(s)
//     return i
// }

func parseFloat(s string) float64 {
    f, _ := strconv.ParseFloat(s, 64)
    return f
}
// // RegisterRoutes registers all the doctor routes
// func (h *DoctorHandler) RegisterRoutes(router *mux.Router) {
//     router.HandleFunc("/doctors", h.CreateDoctor).Methods("POST")
//     router.HandleFunc("/doctors", h.ListDoctors).Methods("GET")
//     router.HandleFunc("/doctors/{id}", h.GetDoctor).Methods("GET")
//     router.HandleFunc("/doctors/{id}", h.UpdateDoctor).Methods("PUT")
// }
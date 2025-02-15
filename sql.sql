-- Users table (base table for all user types)
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role ENUM('doctor', 'patient', 'home_care_provider', 'admin') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- Relationship: One-to-One with doctors, patients, and home_care_providers

-- Service types table
CREATE TABLE service_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_home_care BOOLEAN DEFAULT FALSE
);
-- Relationship: One-to-Many with doctors and home_care_providers

-- Doctors table
CREATE TABLE doctors (
    user_id INT PRIMARY KEY,
    specialty VARCHAR(100) NOT NULL,
    service_type_id INT,
    license_number VARCHAR(50) UNIQUE NOT NULL,
    experience_years INT,
    qualifications TEXT,
    achievements TEXT,
    bio TEXT,
    profile_picture_url VARCHAR(255),
    consultation_fee DECIMAL(10, 2),
    rating DECIMAL(3, 2),
    is_verified BOOLEAN DEFAULT FALSE,
    is_available BOOLEAN DEFAULT TRUE,
    status ENUM('active', 'inactive') DEFAULT 'active',
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (service_type_id) REFERENCES service_types(id),
    INDEX idx_specialty (specialty),
    INDEX idx_service_type (service_type_id),
    INDEX idx_rating (rating),
    INDEX idx_status (status),
    INDEX idx_location (latitude, longitude)
);

-- Relationship: One-to-One with users, One-to-Many with doctor_availability, Many-to-One with service_types

-- Doctor availability
CREATE TABLE doctor_availability (
    id INT AUTO_INCREMENT PRIMARY KEY,
    doctor_id INT NOT NULL,
    day_of_week TINYINT NOT NULL,  -- 0 (Sunday) to 6 (Saturday)
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    FOREIGN KEY (doctor_id) REFERENCES doctors(user_id) ON DELETE CASCADE,
    INDEX idx_doctor_availability (doctor_id, day_of_week, start_time, end_time)
);
-- Relationship: Many-to-One with doctors

-- Home care providers table
CREATE TABLE home_care_providers (
    user_id INT PRIMARY KEY,
    service_type_id INT NOT NULL,
    experience_years INT,
    qualifications TEXT,
    bio TEXT,
    profile_picture_url VARCHAR(255),
    hourly_rate DECIMAL(10, 2),
    rating DECIMAL(3, 2),
    is_verified BOOLEAN DEFAULT FALSE,
    is_available BOOLEAN DEFAULT TRUE,
    status ENUM('active', 'inactive') DEFAULT 'active',
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (service_type_id) REFERENCES service_types(id),
    INDEX idx_service_type (service_type_id),
    INDEX idx_rating (rating),
    INDEX idx_status (status),
    INDEX idx_location (latitude, longitude)
);
-- Relationship: One-to-One with users, Many-to-One with service_types

-- Patients table
CREATE TABLE patients (
    user_id INT PRIMARY KEY,
    date_of_birth DATE,
    gender ENUM('male', 'female', 'other') NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    emergency_contact_name VARCHAR(100),
    emergency_contact_phone VARCHAR(20),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Relationship: One-to-One with users, One-to-Many with medical_history

-- Medical history
CREATE TABLE medical_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    patient_id INT NOT NULL,
    condition_name VARCHAR(100) NOT NULL,
    diagnosis_date DATE,
    treatment TEXT,
    is_current BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (patient_id) REFERENCES patients(user_id) ON DELETE CASCADE,
    INDEX idx_patient_condition (patient_id, condition_name)
);
-- Relationship: Many-to-One with patients

-- Appointments table
CREATE TABLE appointments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    patient_id INT NOT NULL,
    provider_type ENUM('doctor', 'home_care_provider') NOT NULL,
    doctor_id INT,
    home_care_provider_id INT,
    appointment_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status ENUM('scheduled', 'completed', 'cancelled') NOT NULL DEFAULT 'scheduled',
    cancellation_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(user_id) ON DELETE CASCADE,
    FOREIGN KEY (doctor_id) REFERENCES doctors(user_id) ON DELETE CASCADE,
    FOREIGN KEY (home_care_provider_id) REFERENCES home_care_providers(user_id) ON DELETE CASCADE,
    INDEX idx_provider_date (provider_type, doctor_id, home_care_provider_id, appointment_date),
    INDEX idx_patient_date (patient_id, appointment_date),
    INDEX idx_status (status),
    INDEX idx_appointment_date (appointment_date)
);

-- Relationship: Many-to-One with doctors/home_care_providers and patients, One-to-One with consultations/home_care_visits

-- Chat Messages table (your version)
CREATE TABLE chat_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    consultation_id INT NOT NULL,
    sender_type ENUM('doctor', 'patient') NOT NULL,
    sender_id INT NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (consultation_id) REFERENCES consultations(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_consultation_chat (consultation_id, sent_at),
    INDEX idx_sender_chat (sender_type, sender_id)
);

-- Notifications table (your version, which is a good addition)
CREATE TABLE notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    notification_type ENUM('consultation_request', 'chat_message', 'appointment_reminder'),
    message VARCHAR(255),
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);


-- Consultations table
CREATE TABLE consultations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    patient_id INT NOT NULL,
    doctor_id INT NOT NULL,
    status ENUM('requested', 'in_progress', 'completed', 'cancelled') NOT NULL DEFAULT 'requested',
    started_at DATETIME,
    completed_at DATETIME,
    fee DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES users(id),
    FOREIGN KEY (doctor_id) REFERENCES users(id),
    INDEX idx_consultation_status (status),
    INDEX idx_patient_doctor (patient_id, doctor_id),
    INDEX idx_dates (started_at, completed_at)
);
-- Relationship: One-to-One with appointments, consultation_details, and payments

-- Consultation details table
CREATE TABLE consultation_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    consultation_id INT NOT NULL,
    request_details TEXT,
    symptoms TEXT,
    diagnosis TEXT,
    prescription TEXT,
    notes TEXT,
    FOREIGN KEY (consultation_id) REFERENCES consultations(id) ON DELETE CASCADE
);
-- Relationship: One-to-One with consultations

-- Home care visits table
CREATE TABLE home_care_visits (
    id INT AUTO_INCREMENT PRIMARY KEY,
    patient_id INT NOT NULL,
    provider_id INT NOT NULL,
    address TEXT NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    duration_hours DECIMAL(4, 2) NOT NULL,
    special_requirements TEXT,
    status ENUM('scheduled', 'in_progress', 'completed', 'cancelled') DEFAULT 'scheduled',
    FOREIGN KEY (patient_id) REFERENCES patients(user_id),
    FOREIGN KEY (provider_id) REFERENCES home_care_providers(user_id),
    INDEX idx_location (latitude, longitude)
);
-- Relationship: One-to-One with appointments

-- Payments table
CREATE TABLE payments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    consultation_id INT,
    home_care_visit_id INT,
    amount DECIMAL(10, 2) NOT NULL,
    status ENUM('pending', 'paid', 'refunded') DEFAULT 'pending',
    payment_date DATETIME,
    refund_date DATETIME,
    FOREIGN KEY (consultation_id) REFERENCES consultations(id) ON DELETE CASCADE,
    FOREIGN KEY (home_care_visit_id) REFERENCES home_care_visits(id) ON DELETE CASCADE,
    INDEX idx_payment_status (status)
);
-- Relationship: One-to-One with consultations or home_care_visits

-- Reviews table
CREATE TABLE reviews (
    id INT AUTO_INCREMENT PRIMARY KEY,
    patient_id INT NOT NULL,
    review_type ENUM('consultation', 'home_care') NOT NULL,
    consultation_id INT,
    home_care_visit_id INT,
    doctor_id INT,
    home_care_provider_id INT,
    rating TINYINT NOT NULL,
    comment TEXT,
    00.,
    FOREIGN KEY (patient_id) REFERENCES patients(user_id),
    FOREIGN KEY (consultation_id) REFERENCES consultations(id) ON DELETE CASCADE,
    FOREIGN KEY (home_care_visit_id) REFERENCES home_care_visits(id) ON DELETE CASCADE,
    FOREIGN KEY (doctor_id) REFERENCES doctors(user_id),
    FOREIGN KEY (home_care_provider_id) REFERENCES home_care_providers(user_id),
    INDEX idx_doctor_rating (doctor_id, rating),
    INDEX idx_home_care_provider_rating (home_care_provider_id, rating)
);
-- Relationship: Many-to-One with patients, doctors, and home_care_providers

-- System Log Table
CREATE TABLE system_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    user_id INT,
    user_type ENUM('doctor', 'patient', 'home_care_provider', 'admin', 'system') NOT NULL,
    action_type ENUM('login', 'logout', 'create', 'update', 'delete', 'view', 'book', 'cancel', 'complete', 'payment', 'refund', 'other') NOT NULL,
    action_description VARCHAR(255) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT,
    old_value TEXT,
    new_value TEXT,
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    additional_info JSON,
    INDEX idx_timestamp (timestamp),
    INDEX idx_user_id (user_id),
    INDEX idx_action_type (action_type),
    INDEX idx_entity_type (entity_type)
);

-- 3. Enhanced Audit Trail
CREATE TABLE audit_trail (
    id INT AUTO_INCREMENT PRIMARY KEY,
    table_name VARCHAR(50) NOT NULL,
    record_id INT NOT NULL,
    action ENUM('INSERT', 'UPDATE', 'DELETE') NOT NULL,
    changed_fields JSON,
    changed_by INT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (changed_by) REFERENCES users(id)
);


-- Views


-- Update views to use users table for names
CREATE OR REPLACE VIEW vw_available_doctors AS
SELECT d.user_id AS id, u.name, d.specialty, st.name AS service_type, da.day_of_week, da.start_time, da.end_time
FROM doctors d
JOIN users u ON d.user_id = u.id
JOIN service_types st ON d.service_type_id = st.id
JOIN doctor_availability da ON d.user_id = da.doctor_id
WHERE d.is_available = TRUE AND d.status = 'active';

-- Available doctors
CREATE OR REPLACE VIEW vw_available_doctors AS
SELECT 
    d.user_id AS id, 
    u.name, 
    d.specialty, 
    st.name AS service_type, 
    da.day_of_week, 
    da.start_time, 
    da.end_time
FROM 
    doctors d
JOIN 
    users u ON d.user_id = u.id
JOIN 
    service_types st ON d.service_type_id = st.id
JOIN 
    doctor_availability da ON d.user_id = da.doctor_id
WHERE 
    d.is_available = TRUE 
    AND d.status = 'active';


-- Available home care providers
CREATE OR REPLACE VIEW vw_available_home_care_providers AS
SELECT 
    hcp.user_id AS id, 
    u.name, 
    st.name AS service_type, 
    hcp.hourly_rate, 
    hcp.rating
FROM 
    home_care_providers hcp
JOIN 
    users u ON hcp.user_id = u.id
JOIN 
    service_types st ON hcp.service_type_id = st.id
WHERE 
    hcp.is_available = TRUE 
    AND hcp.status = 'active';

-- Upcoming appointments for doctors
CREATE OR REPLACE VIEW vw_doctor_upcoming_appointments AS
SELECT 
    a.id, 
    a.doctor_id, 
    u_doctor.name AS doctor_name, 
    a.patient_id, 
    u_patient.name AS patient_name, 
    a.appointment_date, 
    a.start_time, 
    a.end_time, 
    a.status
FROM 
    appointments a
JOIN 
    doctors d ON a.doctor_id = d.user_id
JOIN 
    patients p ON a.patient_id = p.user_id
JOIN 
    users u_doctor ON d.user_id = u_doctor.id
JOIN 
    users u_patient ON p.user_id = u_patient.id
WHERE 
    a.provider_type = 'doctor' 
    AND a.appointment_date >= CURDATE() 
    AND a.status = 'scheduled'
ORDER BY 
    a.appointment_date, 
    a.start_time;


-- Upcoming appointments for home care providers
CREATE OR REPLACE VIEW vw_home_care_upcoming_appointments AS
SELECT 
    a.id, 
    a.home_care_provider_id, 
    u_provider.name AS provider_name, 
    a.patient_id, 
    u_patient.name AS patient_name, 
    a.appointment_date, 
    a.start_time, 
    a.end_time, 
    a.status
FROM 
    appointments a
JOIN 
    home_care_providers hcp ON a.home_care_provider_id = hcp.user_id
JOIN 
    patients p ON a.patient_id = p.user_id
JOIN 
    users u_provider ON hcp.user_id = u_provider.id
JOIN 
    users u_patient ON p.user_id = u_patient.id
WHERE 
    a.provider_type = 'home_care_provider' 
    AND a.appointment_date >= CURDATE() 
    AND a.status = 'scheduled'
ORDER BY 
    a.appointment_date, 
    a.start_time;


-- Upcoming appointments for patients
CREATE OR REPLACE VIEW vw_patient_upcoming_appointments AS
SELECT 
    a.id, 
    a.patient_id, 
    u_patient.name AS patient_name, 
    a.provider_type,
    COALESCE(d.user_id, hcp.user_id) AS provider_id, 
    COALESCE(u_doctor.name, u_home_care.name) AS provider_name, 
    a.appointment_date, 
    a.start_time, 
    a.end_time, 
    a.status
FROM 
    appointments a
JOIN 
    patients p ON a.patient_id = p.user_id
JOIN 
    users u_patient ON p.user_id = u_patient.id
LEFT JOIN 
    doctors d ON a.doctor_id = d.user_id
LEFT JOIN 
    home_care_providers hcp ON a.home_care_provider_id = hcp.user_id
LEFT JOIN 
    users u_doctor ON d.user_id = u_doctor.id
LEFT JOIN 
    users u_home_care ON hcp.user_id = u_home_care.id
WHERE 
    a.appointment_date >= CURDATE() 
    AND a.status = 'scheduled'
ORDER BY 
    a.appointment_date, 
    a.start_time;


-- Completed consultations with pending payments
CREATE OR REPLACE VIEW vw_pending_payments AS
SELECT 
    CASE 
        WHEN c.id IS NOT NULL THEN c.id
        ELSE hcv.id
    END AS visit_id,
    CASE 
        WHEN c.id IS NOT NULL THEN 'consultation'
        ELSE 'home_care_visit'
    END AS visit_type,
    a.id AS appointment_id,
    CASE 
        WHEN c.id IS NOT NULL THEN c.fee
        ELSE hcv.duration_hours * hcp.hourly_rate
    END AS amount_due,
    a.patient_id, 
    p.name AS patient_name,
    COALESCE(d.id, hcp.id) AS provider_id, 
    COALESCE(d.name, hcp.name) AS provider_name
FROM 
    appointments a
JOIN 
    patients p ON a.patient_id = p.user_id
LEFT JOIN 
    consultations c ON a.id = c.appointment_id
LEFT JOIN 
    home_care_visits hcv ON a.id = hcv.appointment_id
LEFT JOIN 
    doctors d ON a.doctor_id = d.user_id
LEFT JOIN 
    home_care_providers hcp ON a.home_care_provider_id = hcp.user_id
LEFT JOIN 
    payments pay ON (c.id = pay.consultation_id OR hcv.id = pay.home_care_visit_id)
WHERE 
    (c.status = 'completed' OR hcv.status = 'completed') 
    AND (pay.status IS NULL OR pay.status = 'pending');


-- Stored Procedures

-- Continuing from the previous part...

-- Update stored procedures to use users table for names
DELIMITER //

CREATE OR REPLACE PROCEDURE sp_book_appointment(
    IN p_patient_id INT,
    IN p_provider_type ENUM('doctor', 'home_care_provider'),
    IN p_provider_id INT,
    IN p_appointment_date DATE,
    IN p_start_time TIME,
    IN p_end_time TIME,
    OUT p_appointment_id INT,
    OUT p_error_message VARCHAR(255)
)
BEGIN
    -- ... (rest of the procedure remains the same)
    
    -- Update the log entry to use the users table for names
    CALL sp_add_log_entry(
        p_patient_id, 'patient', 'book', 'Booked an appointment',
        'appointment', p_appointment_id, NULL, NULL,
        NULL, NULL, JSON_OBJECT('provider_type', p_provider_type, 'provider_id', p_provider_id)
    );
END //

DELIMITER ;

DELIMITER //

-- Procedure to add a log entry
CREATE PROCEDURE sp_add_log_entry(
    IN p_user_id INT,
    IN p_user_type ENUM('doctor', 'patient', 'home_care_provider', 'admin', 'system'),
    IN p_action_type ENUM('login', 'logout', 'create', 'update', 'delete', 'view', 'book', 'cancel', 'complete', 'payment', 'refund', 'other'),
    IN p_action_description VARCHAR(255),
    IN p_entity_type VARCHAR(50),
    IN p_entity_id INT,
    IN p_old_value TEXT,
    IN p_new_value TEXT,
    IN p_ip_address VARCHAR(45),
    IN p_user_agent VARCHAR(255),
    IN p_additional_info JSON
)
BEGIN
    INSERT INTO system_logs (
        user_id, user_type, action_type, action_description, 
        entity_type, entity_id, old_value, new_value, 
        ip_address, user_agent, additional_info
    ) VALUES (
        p_user_id, p_user_type, p_action_type, p_action_description, 
        p_entity_type, p_entity_id, p_old_value, p_new_value, 
        p_ip_address, p_user_agent, p_additional_info
    );
END //
---------------

DELIMITER //

CREATE PROCEDURE sp_book_appointment(
    IN p_patient_id INT,
    IN p_provider_type ENUM('doctor', 'home_care_provider'),
    IN p_provider_id INT,
    IN p_appointment_date DATE,
    IN p_start_time TIME,
    IN p_end_time TIME,
    OUT p_appointment_id INT,
    OUT p_error_message VARCHAR(255)
)
BEGIN
    DECLARE v_conflict INT DEFAULT 0;
    DECLARE v_provider_available INT DEFAULT 0;
    
    -- Check if the provider is available at the given time
    IF p_provider_type = 'doctor' THEN
        SELECT COUNT(*) INTO v_provider_available
        FROM doctor_availability
        WHERE doctor_id = p_provider_id
          AND day_of_week = WEEKDAY(p_appointment_date)
          AND start_time <= p_start_time
          AND end_time >= p_end_time;
    ELSE
        -- Assume home care providers are always available, or implement a similar availability check
        SET v_provider_available = 1;
    END IF;
    
    IF v_provider_available = 0 THEN
        SET p_error_message = 'Provider is not available at the selected time.';
    ELSE
        -- Check for conflicting appointments
        SELECT COUNT(*) INTO v_conflict
        FROM appointments
        WHERE ((provider_type = 'doctor' AND doctor_id = p_provider_id) 
            OR (provider_type = 'home_care_provider' AND home_care_provider_id = p_provider_id))
          AND appointment_date = p_appointment_date
          AND ((start_time <= p_start_time AND end_time > p_start_time)
            OR (start_time < p_end_time AND end_time >= p_end_time)
            OR (start_time >= p_start_time AND end_time <= p_end_time));
        
        IF v_conflict > 0 THEN
            SET p_error_message = 'There is a conflicting appointment at the selected time.';
        ELSE
            -- Book the appointment
            INSERT INTO appointments (patient_id, provider_type, doctor_id, home_care_provider_id, appointment_date, start_time, end_time)
            VALUES (p_patient_id, p_provider_type, 
                    IF(p_provider_type = 'doctor', p_provider_id, NULL),
                    IF(p_provider_type = 'home_care_provider', p_provider_id, NULL),
                    p_appointment_date, p_start_time, p_end_time);
            
            SET p_appointment_id = LAST_INSERT_ID();
            SET p_error_message = NULL;
            
            -- Log the appointment booking
            CALL sp_add_log_entry(
                p_patient_id, 'patient', 'book', 'Booked an appointment',
                'appointment', p_appointment_id, NULL, NULL,
                NULL, NULL, JSON_OBJECT('provider_type', p_provider_type, 'provider_id', p_provider_id)
            );
        END IF;
    END IF;
END //

DELIMITER ;


DELIMITER //

CREATE PROCEDURE sp_cancel_appointment(
    IN p_appointment_id INT,
    IN p_cancellation_reason TEXT,
    IN p_user_id INT,
    IN p_user_type ENUM('patient', 'doctor', 'home_care_provider', 'admin'),
    OUT p_error_message VARCHAR(255)
)
BEGIN
    DECLARE v_appointment_status VARCHAR(20);
    DECLARE v_patient_id INT;
    DECLARE v_provider_type VARCHAR(20);
    DECLARE v_provider_id INT;
    
    -- Check if the appointment exists and get its status
    SELECT status, patient_id, provider_type, 
           CASE WHEN provider_type = 'doctor' THEN doctor_id ELSE home_care_provider_id END
    INTO v_appointment_status, v_patient_id, v_provider_type, v_provider_id
    FROM appointments
    WHERE id = p_appointment_id;
    
    IF v_appointment_status IS NULL THEN
        SET p_error_message = 'Appointment not found.';
    ELSEIF v_appointment_status != 'scheduled' THEN
        SET p_error_message = 'Only scheduled appointments can be cancelled.';
    ELSEIF p_user_type = 'patient' AND p_user_id != v_patient_id THEN
        SET p_error_message = 'You can only cancel your own appointments.';
    ELSEIF p_user_type IN ('doctor', 'home_care_provider') AND p_user_id != v_provider_id THEN
        SET p_error_message = 'You can only cancel appointments assigned to you.';
    ELSE
        -- Cancel the appointment
        UPDATE appointments
        SET status = 'cancelled', cancellation_reason = p_cancellation_reason
        WHERE id = p_appointment_id;
        
        SET p_error_message = NULL;
        
        -- Log the cancellation
        CALL sp_add_log_entry(
            p_user_id, p_user_type, 'cancel', 'Cancelled an appointment',
            'appointment', p_appointment_id, NULL, NULL,
            NULL, NULL, JSON_OBJECT('reason', p_cancellation_reason)
        );
    END IF;
END //

DELIMITER //

-- Procedure to find nearby doctors
DELIMITER //

-- Procedure to find nearby doctors
CREATE PROCEDURE sp_find_nearby_doctors(
    IN p_latitude DECIMAL(10, 8),
    IN p_longitude DECIMAL(11, 8),
    IN p_max_distance INT,  -- in kilometers
    IN p_specialty VARCHAR(100)
)
BEGIN
    SELECT d.user_id, d.specialty, d.rating,  -- Removed 'd.name'
           ST_Distance_Sphere(
               POINT(p_longitude, p_latitude),
               POINT(d.longitude, d.latitude)
           ) / 1000 AS distance_km
    FROM doctors d
    WHERE d.status = 'active' AND d.is_available = TRUE
      AND (p_specialty IS NULL OR d.specialty = p_specialty)
      AND ST_Distance_Sphere(
          POINT(p_longitude, p_latitude),
          POINT(d.longitude, d.latitude)
      ) / 1000 <= p_max_distance
    ORDER BY distance_km;
END //

DELIMITER ;


DELIMITER //

-- Procedure to find nearby home care providers
CREATE PROCEDURE sp_find_nearby_home_care_providers(
    IN p_latitude DECIMAL(10, 8),
    IN p_longitude DECIMAL(11, 8),
    IN p_max_distance INT,  -- in kilometers
    IN p_service_type_id INT
)
BEGIN
    SELECT hcp.user_id AS provider_id,  -- Changed 'hcp.id' to 'hcp.user_id'
           -- hcp.name,  -- Remove this line unless you add a 'name' column
           st.name AS service_type, 
           hcp.hourly_rate, 
           hcp.rating,
           ST_Distance_Sphere(
               POINT(p_longitude, p_latitude),
               POINT(hcp.longitude, hcp.latitude)
           ) / 1000 AS distance_km
    FROM home_care_providers hcp
    JOIN service_types st ON hcp.service_type_id = st.id
    WHERE hcp.status = 'active' AND hcp.is_available = TRUE
      AND (p_service_type_id IS NULL OR hcp.service_type_id = p_service_type_id)
      AND ST_Distance_Sphere(
          POINT(p_longitude, p_latitude),
          POINT(hcp.longitude, hcp.latitude)
      ) / 1000 <= p_max_distance
    ORDER BY distance_km;
END //

DELIMITER ;
CALL sp_find_nearby_home_care_providers(34.052235, -118.243683, 10, 1);





-- Soft Deletes
ALTER TABLE users
ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;

ALTER TABLE doctors
ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;

ALTER TABLE patients
ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;

ALTER TABLE home_care_providers
ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;


ALTER TABLE medical_history
ADD COLUMN version INT DEFAULT 1;

-- Create a trigger to increment version on update
DELIMITER //
CREATE TRIGGER trg_medical_history_version
BEFORE UPDATE ON medical_history
FOR EACH ROW
BEGIN
    SET NEW.version = OLD.version + 1;
END //
DELIMITER ;



-- 1. Enhance the doctors and home_care_providers tables with additional searchable fields

ALTER TABLE doctors
ADD FULLTEXT INDEX ft_idx_doctor_search (specialty, qualifications, bio);

ALTER TABLE home_care_providers
ADD FULLTEXT INDEX ft_idx_provider_search (qualifications, bio);

-- 2. Create a unified view for searching both doctors and home care providers

CREATE OR REPLACE VIEW vw_healthcare_providers AS
SELECT 
    'doctor' AS provider_type,
    d.user_id AS provider_id,
    u.name AS provider_name,
    d.specialty AS specialization,
    d.service_type_id,
    st.name AS service_type_name,
    d.qualifications,
    d.bio,
    d.rating,
    d.is_verified,
    d.consultation_fee AS fee,
    d.latitude,
    d.longitude
FROM doctors d
JOIN users u ON d.user_id = u.id
JOIN service_types st ON d.service_type_id = st.id
WHERE d.is_available = TRUE AND d.status = 'active' AND d.deleted_at IS NULL

UNION ALL

SELECT 
    'home_care_provider' AS provider_type,
    hcp.user_id AS provider_id,
    u.name AS provider_name,
    st.name AS specialization,
    hcp.service_type_id,
    st.name AS service_type_name,
    hcp.qualifications,
    hcp.bio,
    hcp.rating,
    hcp.is_verified,
    hcp.hourly_rate AS fee,
    hcp.latitude,
    hcp.longitude
FROM home_care_providers hcp
JOIN users u ON hcp.user_id = u.id
JOIN service_types st ON hcp.service_type_id = st.id
WHERE hcp.is_available = TRUE AND hcp.status = 'active' AND hcp.deleted_at IS NULL;

-- 3. Create a stored procedure for advanced searching

DELIMITER //

CREATE PROCEDURE sp_search_healthcare_providers(
    IN p_provider_type VARCHAR(20),
    IN p_search_term VARCHAR(255),
    IN p_service_type_id INT,
    IN p_min_rating DECIMAL(3,2),
    IN p_max_fee DECIMAL(10,2),
    IN p_is_verified BOOLEAN,
    IN p_latitude DECIMAL(10,8),
    IN p_longitude DECIMAL(11,8),
    IN p_max_distance INT,
    IN p_sort_by VARCHAR(20),
    IN p_sort_order VARCHAR(4),
    IN p_limit INT,
    IN p_offset INT
)
BEGIN
    SET @sql = CONCAT('
        SELECT 
            provider_type,
            provider_id,
            provider_name,
            specialization,
            service_type_name,
            qualifications,
            bio,
            rating,
            is_verified,
            fee,
            latitude,
            longitude,
            CASE 
                WHEN latitude IS NOT NULL AND longitude IS NOT NULL THEN
                    ST_Distance_Sphere(
                        POINT(', p_longitude, ', ', p_latitude, '),
                        POINT(longitude, latitude)
                    ) / 1000
                ELSE NULL
            END AS distance_km
        FROM vw_healthcare_providers
        WHERE 1=1
    ');

    IF p_provider_type IS NOT NULL THEN
        SET @sql = CONCAT(@sql, ' AND provider_type = ''', p_provider_type, '''');
    END IF;

    IF p_search_term IS NOT NULL THEN
        SET @sql = CONCAT(@sql, ' AND MATCH(specialization, qualifications, bio) AGAINST(''', p_search_term, ''' IN NATURAL LANGUAGE MODE)');
    END IF;

    IF p_service_type_id IS NOT NULL THEN
        SET @sql = CONCAT(@sql, ' AND service_type_id = ', p_service_type_id);
    END IF;

    IF p_min_rating IS NOT NULL THEN
        SET @sql = CONCAT(@sql, ' AND rating >= ', p_min_rating);
    END IF;

    IF p_max_fee IS NOT NULL THEN
        SET @sql = CONCAT(@sql, ' AND fee <= ', p_max_fee);
    END IF;

    IF p_is_verified IS NOT NULL THEN
        SET @sql = CONCAT(@sql, ' AND is_verified = ', p_is_verified);
    END IF;

    IF p_max_distance IS NOT NULL AND p_latitude IS NOT NULL AND p_longitude IS NOT NULL THEN
        SET @sql = CONCAT(@sql, ' AND ST_Distance_Sphere(POINT(', p_longitude, ', ', p_latitude, '), POINT(longitude, latitude)) / 1000 <= ', p_max_distance);
    END IF;

    -- Sorting
    IF p_sort_by IS NOT NULL THEN
        SET @sql = CONCAT(@sql, ' ORDER BY ');
        CASE p_sort_by
            WHEN 'rating' THEN SET @sql = CONCAT(@sql, 'rating');
            WHEN 'fee' THEN SET @sql = CONCAT(@sql, 'fee');
            WHEN 'distance' THEN SET @sql = CONCAT(@sql, 'distance_km');
            ELSE SET @sql = CONCAT(@sql, 'rating'); -- Default sort
        END CASE;
        
        IF p_sort_order IN ('ASC', 'DESC') THEN
            SET @sql = CONCAT(@sql, ' ', p_sort_order);
        ELSE
            SET @sql = CONCAT(@sql, ' DESC'); -- Default to descending
        END IF;
    END IF;

    -- Pagination
    SET @sql = CONCAT(@sql, ' LIMIT ', IFNULL(p_limit, 10), ' OFFSET ', IFNULL(p_offset, 0));

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END //

DELIMITER ;

-- 4. Example usage of the search procedure

CALL sp_search_healthcare_providers(
    'doctor',                -- p_provider_type
    'pediatric allergy',     -- p_search_term
    NULL,                    -- p_service_type_id
    4.0,                     -- p_min_rating
    200.00,                  -- p_max_fee
    TRUE,                    -- p_is_verified
    40.7128,                 -- p_latitude (New York City)
    -74.0060,                -- p_longitude
    50,                      -- p_max_distance (km)
    'rating',                -- p_sort_by
    'DESC',                  -- p_sort_order
    10,                      -- p_limit
    0                        -- p_offset
);


-----------------------------

CREATE PROCEDURE search_doctors(
    search_term VARCHAR(100),
    specialty VARCHAR(100),
    min_rating DECIMAL(3,2),
    location_lat DECIMAL(10,8),
    location_lng DECIMAL(11,8),
    radius_km INT,
    service_type_id INT
)
BEGIN
    SELECT 
        d.*,
        u.first_name,
        u.last_name,
        u.email,
        u.phone,
        -- Calculate distance if coordinates provided
        CASE 
            WHEN location_lat IS NOT NULL AND location_lng IS NOT NULL THEN
                (6371 * acos(
                    cos(radians(location_lat)) * 
                    cos(radians(d.latitude)) * 
                    cos(radians(d.longitude) - radians(location_lng)) + 
                    sin(radians(location_lat)) * 
                    sin(radians(d.latitude))
                ))
            ELSE NULL
        END as distance_km
    FROM doctors d
    INNER JOIN users u ON d.user_id = u.id
    WHERE 1=1
        -- Name search
        AND (search_term IS NULL 
             OR CONCAT(u.first_name, ' ', u.last_name) LIKE CONCAT('%', search_term, '%'))
        -- Specialty filter  
        AND (specialty IS NULL 
             OR d.specialty = specialty)
        -- Rating filter
        AND (min_rating IS NULL 
             OR d.rating >= min_rating)
        -- Service type filter
        AND (service_type_id IS NULL 
             OR d.service_type_id = service_type_id)
        -- Location filter
        AND (location_lat IS NULL 
             OR location_lng IS NULL 
             OR radius_km IS NULL
             OR (6371 * acos(
                    cos(radians(location_lat)) * 
                    cos(radians(d.latitude)) * 
                    cos(radians(d.longitude) - radians(location_lng)) + 
                    sin(radians(location_lat)) * 
                    sin(radians(d.latitude))
                )) <= radius_km)
        -- Only active and available doctors
        AND d.status = 'active'
        AND d.is_available = TRUE
    ORDER BY
        CASE 
            WHEN location_lat IS NOT NULL THEN distance_km
            ELSE d.rating
        END ASC,
        d.rating DESC
    LIMIT 100;
END;
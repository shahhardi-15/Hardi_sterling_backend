package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	godotenv.Load(".env")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbName == "" {
		dbName = "sterling_hms"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✓ Connected to database")

	// Step 1: Add columns to doctors table
	fmt.Println("\nStep 1: Adding columns to doctors table...")
	alterQuery := `
		ALTER TABLE doctors
		ADD COLUMN IF NOT EXISTS experience_years INTEGER DEFAULT 0,
		ADD COLUMN IF NOT EXISTS qualification VARCHAR(255),
		ADD COLUMN IF NOT EXISTS address TEXT;
	`

	_, err = db.Exec(alterQuery)
	if err != nil {
		log.Fatalf("Failed to alter doctors table: %v", err)
	}
	fmt.Println("✓ Columns added successfully")

	// Step 2: Update existing doctors
	fmt.Println("\nStep 2: Updating existing doctors...")
	updates := []struct {
		name            string
		experienceYears int
		qualification   string
		address         string
	}{
		{"Dr. John Smith", 15, "MD - General Medicine", "123 Medical Center, City Hospital, Street 1, Main City"},
		{"Dr. Sarah Johnson", 12, "MD - Cardiology, Board Certified", "456 Cardiology Wing, City Hospital, Heart Street, Main City"},
		{"Dr. Michael Brown", 10, "MD - Dermatology", "789 Dermatology Clinic, City Hospital, Skin Lane, Main City"},
		{"Dr. Emily Davis", 8, "MD - Neurology", "321 Neurology Department, City Hospital, Brain Avenue, Main City"},
		{"Dr. Robert Wilson", 18, "MD - Orthopedic Surgery, Board Certified", "654 Orthopedic Center, City Hospital, Bone Street, Main City"},
	}

	for _, doc := range updates {
		updateQuery := `
			UPDATE doctors 
			SET experience_years = $1, qualification = $2, address = $3
			WHERE name = $4
		`
		_, err := db.Exec(updateQuery, doc.experienceYears, doc.qualification, doc.address, doc.name)
		if err != nil {
			log.Printf("Warning: Failed to update %s: %v", doc.name, err)
		}
	}
	fmt.Println("✓ Existing doctors updated")

	// Step 3: Insert new doctors
	fmt.Println("\nStep 3: Inserting new doctors...")
	query := `
		INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available) 
		VALUES 
		  ($1, $2, $3, $4, $5, $6, $7, $8),
		  ($9, $10, $11, $12, $13, $14, $15, $16),
		  ($17, $18, $19, $20, $21, $22, $23, $24),
		  ($25, $26, $27, $28, $29, $30, $31, $32),
		  ($33, $34, $35, $36, $37, $38, $39, $40),
		  ($41, $42, $43, $44, $45, $46, $47, $48),
		  ($49, $50, $51, $52, $53, $54, $55, $56),
		  ($57, $58, $59, $60, $61, $62, $63, $64),
		  ($65, $66, $67, $68, $69, $70, $71, $72),
		  ($73, $74, $75, $76, $77, $78, $79, $80)
	`

	doctors := []interface{}{
		// Dr. Lisa Anderson - Pediatrician
		"Dr. Lisa Anderson", "Pediatrician", "lisa.anderson@hospital.com", "+1234567895", 11,
		"MD - Pediatrics", "111 Children's Hospital, Pediatric Wing, Main City", true,

		// Dr. Jennifer Martinez - Gynecologist
		"Dr. Jennifer Martinez", "Gynecologist", "jennifer.martinez@hospital.com", "+1234567896", 13,
		"MD - Obstetrics & Gynecology, Board Certified", "222 Women's Health Center, Main City", true,

		// Dr. James Wilson - Psychiatrist
		"Dr. James Wilson", "Psychiatrist", "james.wilson@hospital.com", "+1234567897", 16,
		"MD - Psychiatry, Board Certified", "333 Mental Health Institute, Main City", true,

		// Dr. Patricia Lee - Ophthalmologist
		"Dr. Patricia Lee", "Ophthalmologist", "patricia.lee@hospital.com", "+1234567898", 14,
		"MD - Ophthalmology, Board Certified", "444 Eye Care Center, Main City", true,

		// Dr. David Kim - Dentist
		"Dr. David Kim", "Dentist", "david.kim@hospital.com", "+1234567899", 9,
		"DDS - Dentistry", "555 Dental Clinic, Main City", true,

		// Dr. Thomas Garcia - Emergency Medicine Specialist
		"Dr. Thomas Garcia", "Emergency Medicine Specialist", "thomas.garcia@hospital.com", "+1234567900", 12,
		"MD - Emergency Medicine, Board Certified", "666 Emergency Department, City Hospital, Main City", true,

		// Dr. Maria Rodriguez - Critical Care Specialist
		"Dr. Maria Rodriguez", "Critical Care Specialist", "maria.rodriguez@hospital.com", "+1234567901", 15,
		"MD - Critical Care Medicine, Board Certified", "777 ICU Unit, City Hospital, Main City", true,

		// Dr. Christopher Taylor - ENT Specialist
		"Dr. Christopher Taylor", "ENT Specialist (Otolaryngologist)", "christopher.taylor@hospital.com", "+1234567902", 13,
		"MD - Otolaryngology, Board Certified", "888 ENT Center, Main City", true,

		// Dr. Susan White - General Surgeon
		"Dr. Susan White", "Surgeon (General Surgeon)", "susan.white@hospital.com", "+1234567903", 18,
		"MD - General Surgery, Board Certified", "999 Surgical Suite, City Hospital, Main City", true,

		// Dr. Andrew Miller - Plastic Surgeon
		"Dr. Andrew Miller", "Plastic Surgeon", "andrew.miller@hospital.com", "+1234567904", 14,
		"MD - Plastic Surgery, Board Certified", "1010 Cosmetic Surgery Center, Main City", true,
	}

	result, err := db.Exec(query, doctors...)
	if err != nil {
		log.Fatalf("Failed to insert doctors: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Failed to get rows affected: %v", err)
	}

	fmt.Printf("✓ Successfully inserted %d new doctors\n", rowsAffected)

	// Step 4: Query all specializations
	fmt.Println("\nStep 4: Listing all specializations...")
	specQuery := `SELECT DISTINCT specialization FROM doctors WHERE is_available = true ORDER BY specialization`
	rows, err := db.Query(specQuery)
	if err != nil {
		log.Fatalf("Failed to query specializations: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n✓ Available Specializations:")
	count := 0
	for rows.Next() {
		var spec string
		err := rows.Scan(&spec)
		if err != nil {
			log.Fatalf("Failed to scan specialization: %v", err)
		}
		count++
		fmt.Printf("  %2d. %s\n", count, spec)
	}

	fmt.Printf("\n✓ Total specializations: %d\n", count)

	// Step 5: Query total doctors count
	countQuery := `SELECT COUNT(*) FROM doctors WHERE is_available = true`
	var totalDoctors int
	err = db.QueryRow(countQuery).Scan(&totalDoctors)
	if err != nil {
		log.Fatalf("Failed to count doctors: %v", err)
	}

	fmt.Printf("✓ Total doctors available: %d\n", totalDoctors)
	fmt.Println("\n✅ Migration completed successfully!")
}

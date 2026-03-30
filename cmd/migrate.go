package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection parameters
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "admin"
	dbname := "sterling"

	// Build connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to database
	log.Println("Connecting to PostgreSQL database...")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✓ Successfully connected to PostgreSQL database")

	// Migration files to execute
	migrations := []string{
		"database/receptionist_migration.sql",
		"database/admin_patients_migration.sql",
		// "database/admin_doctors_migration.sql", // Skip for now - different schema than current doctors table
		"database/sample_doctors_with_specializations.sql",
		"database/admin_doctors_data_v3.sql",
	}

	// Execute each migration
	for _, migrationFile := range migrations {
		err := executeMigration(db, migrationFile)
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	}

	log.Println("✓ All migrations completed successfully!")
}

func executeMigration(db *sql.DB, migrationFile string) error {
	log.Printf("\nExecuting migration: %s", migrationFile)

	// Read the SQL file
	sqlPath := filepath.Join("cmd", "..", migrationFile)
	sqlBytes, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
	}

	sqlContent := string(sqlBytes)

	// Split the SQL content into individual statements
	// Handle multiple statements separated by semicolons
	statements := splitStatements(sqlContent)

	if len(statements) == 0 {
		log.Printf("  ⚠ No SQL statements found in %s", migrationFile)
		return nil
	}

	log.Printf("  Found %d SQL statements in %s", len(statements), migrationFile)

	// Execute each statement
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		log.Printf("  Executing statement %d/%d...", i+1, len(statements))
		_, err := db.Exec(statement)
		if err != nil {
			return fmt.Errorf("failed to execute statement %d in %s: %w\n%s", i+1, migrationFile, err, statement)
		}
	}

	log.Printf("✓ Successfully executed %s", migrationFile)
	return nil
}

// splitStatements splits SQL content by semicolons, handling comments and strings
func splitStatements(sql string) []string {
	var statements []string
	var currentStatement strings.Builder
	inString := false
	inLineComment := false
	inBlockComment := false
	stringChar := rune(0)

	lines := strings.Split(sql, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Handle line comments
		if strings.HasPrefix(trimmedLine, "--") {
			continue
		}

		for _, char := range line {
			// Handle block comments
			if !inString && !inLineComment {
				if char == '/' && currentStatement.Len() > 0 {
					// Check if this starts a block comment
					lastChar := []rune(currentStatement.String())[currentStatement.Len()-1]
					if lastChar == '/' {
						// This is a bit tricky, we need to look ahead
						inBlockComment = true
						currentStatement.WriteRune(char)
						continue
					}
				}
				if inBlockComment && char == '/' {
					// Check if previous was *
					str := currentStatement.String()
					if len(str) > 0 && str[len(str)-1] == '*' {
						inBlockComment = false
						currentStatement.WriteRune(char)
						continue
					}
				}
			}

			if !inString && !inBlockComment {
				if (char == '\'' || char == '"') && !inLineComment {
					inString = true
					stringChar = char
					currentStatement.WriteRune(char)
					continue
				} else if char == ';' {
					// End of statement
					currentStatement.WriteRune(char)
					statement := strings.TrimSpace(currentStatement.String())
					if statement != "" && statement != ";" {
						statements = append(statements, strings.TrimSuffix(statement, ";"))
					}
					currentStatement.Reset()
					inString = false
					inLineComment = false
					continue
				}
			} else if inString && char == stringChar {
				// Check if it's escaped
				inString = false
				stringChar = 0
				currentStatement.WriteRune(char)
				continue
			}

			currentStatement.WriteRune(char)
		}
	}

	// Add any remaining statement
	statement := strings.TrimSpace(currentStatement.String())
	if statement != "" && statement != ";" {
		statements = append(statements, strings.TrimSuffix(statement, ";"))
	}

	return statements
}

package handlers

import (
	"log"
	"net/http"
	"sterling-hms-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type AdminDepartmentHandler struct{}

func NewAdminDepartmentHandler() *AdminDepartmentHandler {
	return &AdminDepartmentHandler{}
}

// ListDepartments retrieves all active departments
func (h *AdminDepartmentHandler) ListDepartments(c *gin.Context) {
	query := `
		SELECT id, name FROM departments 
		WHERE is_active = true
		ORDER BY name ASC
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Printf("Error querying departments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch departments",
			"success": false,
		})
		return
	}
	defer rows.Close()

	type Department struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var departments []Department
	for rows.Next() {
		dept := Department{}
		err := rows.Scan(&dept.ID, &dept.Name)
		if err != nil {
			log.Printf("Error scanning department row: %v", err)
			continue
		}
		departments = append(departments, dept)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Departments retrieved successfully",
		"success": true,
		"data":    departments,
	})
}

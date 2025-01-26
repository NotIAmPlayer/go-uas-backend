package controllers

import (
	"meeting-backend/config"
	"strings"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Staff struct {
	StaffID    int    `json:"staff_id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	PositionID int    `json:"position_id"`
}

func GetStaffs(c *gin.Context) {
	staffs := []Staff{}
	var query = "SELECT Staff_ID, Full_Name, Email, Position_ID FROM Staff"

	rows, err := config.DB.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer rows.Close()

	for rows.Next() {
		var s Staff

		if err := rows.Scan(&s.StaffID, &s.FullName, &s.Email, &s.PositionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		staffs = append(staffs, s)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if len(staffs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "staffs table is empty"})
	} else {
		c.JSON(http.StatusOK, staffs)
	}
}

func GetStaffByID(c *gin.Context) {
	id := c.Param("id")

	var query = "SELECT Staff_ID, Full_Name, Email, Position_ID FROM Staff WHERE Staff_ID = ?"

	res, err := config.DB.Query(query, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	if res.Next() {
		var s Staff

		if err := res.Scan(&s.StaffID, &s.FullName, &s.Email, &s.PositionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		c.IndentedJSON(http.StatusOK, s)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "staff not found"})
	}
}

func PostStaff(c *gin.Context) {
	var newStaff Staff

	if err := c.BindJSON(&newStaff); err != nil {
		return
	}

	var query = "INSERT INTO staff (Full_Name, Email, Position_ID) VALUES (?, ?, ?)"

	ins, err := config.DB.Query(query, newStaff.FullName, newStaff.Email, newStaff.PositionID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "staff creation failed."})
	} else {
		c.IndentedJSON(http.StatusCreated, newStaff)
	}

	defer ins.Close()
}

func PutStaff(c *gin.Context) {
	var updatedStaff Staff
	updatedStaff.PositionID = -1

	if err := c.BindJSON(&updatedStaff); err != nil {
		return
	}

	updates := []string{}
	args := []interface{}{}

	id := c.Param("id")

	if updatedStaff.FullName != "" {
		updates = append(updates, "Full_Name = ?")
		args = append(args, updatedStaff.FullName)
	}

	if updatedStaff.Email != "" {
		updates = append(updates, "Email = ?")
		args = append(args, updatedStaff.Email)
	}

	if updatedStaff.PositionID != -1 {
		updates = append(updates, "Position_ID = ?")
		args = append(args, updatedStaff.PositionID)
	}

	if len(updates) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no updates provided"})
		return
	}

	args = append(args, id)

	query := "UPDATE staff SET " + strings.Join(updates, ", ") + " WHERE Staff_ID = ?"

	res, err := config.DB.Query(query, args...)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	defer res.Close()

	c.IndentedJSON(http.StatusOK, gin.H{"message": "staff " + id + " updated."})
}

func DeleteStaff(c *gin.Context) {
	id := c.Param("id")

	var query = "DELETE FROM staff WHERE Staff_ID = ?"

	res, err := config.DB.Query(query, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	defer res.Close()

	c.IndentedJSON(http.StatusOK, gin.H{"message": "staff " + id + " deleted."})
}

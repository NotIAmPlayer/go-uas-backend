package controllers

import (
	"log"
	"meeting-backend/config"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Location struct {
	LocationID   int    `json:"location_id"`
	LocationName string `json:"location_name"`
	Address      string `json:"address"`
	Floor        string `json:"floor"`
}

func GetLocations(c *gin.Context) {
	locations := []Location{}
	var query = "SELECT Location_ID, Location_Name, Address, Floor FROM Location"

	rows, err := config.DB.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer rows.Close()

	for rows.Next() {
		var l Location

		if err := rows.Scan(&l.LocationID, &l.LocationName, &l.Address, &l.Floor); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		locations = append(locations, l)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if len(locations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "staffs table is empty"})
	} else {
		c.JSON(http.StatusOK, locations)
	}
}

func GetLocationByID(c *gin.Context) {
	id := c.Param("id")

	var query = "SELECT Location_ID, Location_Name, Address, Floor FROM Location WHERE Location_ID = ?"

	res, err := config.DB.Query(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	if res.Next() {
		var l Location

		if err := res.Scan(&l.LocationID, &l.LocationName, &l.Address, &l.Floor); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		c.IndentedJSON(http.StatusOK, l)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "location not found"})
	}
}

func PostLocation(c *gin.Context) {
	var newLocation Location

	if err := c.BindJSON(&newLocation); err != nil {
		return
	}

	var query = "INSERT INTO location (Location_Name, Address, Floor) VALUES (?, ?, ?)"

	ins, err := config.DB.Query(query, newLocation.LocationName, newLocation.Address, newLocation.Floor)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "location creation failed."})
	} else {
		c.IndentedJSON(http.StatusCreated, newLocation)
	}

	defer ins.Close()
}

func PutLocation(c *gin.Context) {
	var updatedLocation Location

	if err := c.BindJSON(&updatedLocation); err != nil {
		return
	}

	updates := []string{}
	args := []interface{}{}

	id := c.Param("id")

	if updatedLocation.LocationName != "" {
		updates = append(updates, "Location_Name = ?")
		args = append(args, updatedLocation.LocationName)
	}

	if updatedLocation.Address != "" {
		updates = append(updates, "Address = ?")
		args = append(args, updatedLocation.Address)
	}

	if updatedLocation.Floor != "" {
		updates = append(updates, "Floor = ?")
		args = append(args, updatedLocation.Floor)
	}

	if len(updates) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no updates provided"})
		return
	}

	args = append(args, id)

	query := "UPDATE location SET " + strings.Join(updates, ", ") + " WHERE Location_ID = ?"

	res, err := config.DB.Query(query, args...)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "location " + id + " updated."})
	}

	defer res.Close()
}

func DeleteLocation(c *gin.Context) {
	id := c.Param("id")

	var query = "DELETE FROM location WHERE Location_ID = ?"

	res, err := config.DB.Query(query, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	c.IndentedJSON(http.StatusOK, gin.H{"message": "location " + id + " deleted."})
}

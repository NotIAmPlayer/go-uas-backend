package controllers

import (
	"meeting-backend/config"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Meeting struct {
	MeetingID   int    `json:"meeting_id"`
	LocationID  int    `json:"location_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MeetingDate string `json:"meeting_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	MeetingType string `json:"meeting_type"`
}

func GetMeetings(c *gin.Context) {
	var meetings []Meeting
	var query = "SELECT Meeting_ID, Location_ID, Title, Description, Meeting_Date, Start_Time, End_Time, Meeting_Type FROM Meeting"

	rows, err := config.DB.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer rows.Close()

	for rows.Next() {
		var m Meeting

		if err := rows.Scan(&m.MeetingID, &m.LocationID, &m.Title, &m.Description, &m.MeetingDate, &m.StartTime, &m.EndTime, &m.MeetingType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}

		meetings = append(meetings, m)
	}

	c.JSON(http.StatusOK, meetings)
}

func GetMeetingByID(c *gin.Context) {
	id := c.Param("id")

	var query = "SELECT Meeting_ID, Location_ID, Title, Description, Meeting_Date, Start_Time, End_Time, Meeting_Type FROM Meeting WHERE Meeting_ID = ?"

	res, err := config.DB.Query(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	if res.Next() {
		var m Meeting

		if err := res.Scan(&m.MeetingID, &m.LocationID, &m.Title, &m.Description, &m.MeetingDate, &m.StartTime, &m.EndTime, &m.MeetingType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}

		c.JSON(http.StatusOK, m)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "meeting not found"})
	}
}

func PostMeeting(c *gin.Context) {
	var newMeeting Meeting

	if err := c.BindJSON(&newMeeting); err != nil {
		return
	}

	var query = "INSERT INTO meeting (Location_ID, Title, Description, Meeting_Date, Start_Time, End_Time, Meeting_Type) VALUES (?, ?, ?, ?, ?, ?, ?)"

	ins, err := config.DB.Query(query, newMeeting.LocationID, newMeeting.Title, newMeeting.Description, newMeeting.MeetingDate, newMeeting.StartTime, newMeeting.EndTime, newMeeting.MeetingType)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "meeting creation failed."})
	} else {
		c.JSON(http.StatusCreated, newMeeting)
	}

	defer ins.Close()
}

func PutMeeting(c *gin.Context) {
	var updatedMeeting Meeting
	updatedMeeting.LocationID = -1

	if err := c.BindJSON(&updatedMeeting); err != nil {
		return
	}

	updates := []string{}
	args := []interface{}{}

	id := c.Param("id")

	if updatedMeeting.LocationID != -1 {
		updates = append(updates, "Location_ID = ?")
		args = append(args, updatedMeeting.LocationID)
	}

	if updatedMeeting.Title != "" {
		updates = append(updates, "Title = ?")
		args = append(args, updatedMeeting.Title)
	}

	if updatedMeeting.Description != "" {
		updates = append(updates, "Description = ?")
		args = append(args, updatedMeeting.Description)
	}

	if updatedMeeting.MeetingDate != "" {
		updates = append(updates, "Meeting_Date = ?")
		args = append(args, updatedMeeting.MeetingDate)
	}

	if updatedMeeting.StartTime != "" {
		updates = append(updates, "Start_Time = ?")
		args = append(args, updatedMeeting.StartTime)
	}

	if updatedMeeting.EndTime != "" {
		updates = append(updates, "End_Time = ?")
		args = append(args, updatedMeeting.EndTime)
	}

	if updatedMeeting.MeetingType != "" {
		updates = append(updates, "Meeting_Type = ?")
		args = append(args, updatedMeeting.MeetingType)
	}

	if len(updates) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no updates provided"})
		return
	}

	args = append(args, id)

	query := "UPDATE meeting SET " + strings.Join(updates, ", ") + " WHERE Meeting_ID = ?"

	res, err := config.DB.Query(query, args...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "meeting " + id + " updated."})
	}

	defer res.Close()
}

func DeleteMeeting(c *gin.Context) {
	id := c.Param("id")

	var query = "DELETE FROM meeting WHERE Meeting_ID = ?"

	res, err := config.DB.Query(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	c.JSON(http.StatusOK, gin.H{"message": "meeting " + id + " deleted."})
}

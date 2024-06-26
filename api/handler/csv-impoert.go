package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ImportContacts godoc
// @Security ApiKeyAuth
// @Router /contacts/import [POST]
// @Summary Import contacts from CSV
// @Description Import contacts from a CSV file
// @Tags contacts
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Success 200 {string} string "Successfully imported"
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
func (h Handler) ImportContacts(c *gin.Context) {
	// Log the start of the request handling
	fmt.Println("ImportContacts handler started")

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving file from form:", err)
		handleResponseLog(c, h.Log, "No file received", http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("File received:", file.Filename)

	src, err := file.Open()
	if err != nil {
		fmt.Println("Error opening the file:", err)
		handleResponseLog(c, h.Log, "Unable to open the file", http.StatusBadRequest, err.Error())
		return
	}
	defer src.Close()

	fmt.Println("File opened successfully")

	err = h.Services.Contacts().ImportContacts(context.Background(), src)
	if err != nil {
		fmt.Println("Error importing contacts:", err)
		handleResponseLog(c, h.Log, "Error importing contacts", http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("Contacts imported successfully")
	handleResponseLog(c, h.Log, "Successfully imported", http.StatusOK, "Successfully imported")
}

// ExportToCSV godoc
// @Security ApiKeyAuth
// @Router /contacts/export/csv [GET]
// @Summary Export contacts to CSV
// @Description Export contacts to a CSV file
// @Tags contacts
// @Produce text/csv
// @Success 200 {file} text/csv
// @Failure 500 {object} model.Response
func (h *Handler) ExportToCSV(c *gin.Context) {

	tableName := "contact"
	outputFile := "output_file.csv"

	err := h.Services.Contactcsv().ExportToCSV(context.Background(), tableName, outputFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.String(http.StatusOK, string(data))
}

package handlers

import (
	"net/http"

	"hospital-management/backend/database"
	"hospital-management/backend/models"
	"hospital-management/backend/utils"
)

var patients = []models.Patient{}

func GetPatientsHandler(w http.ResponseWriter, r *http.Request) {
    patients, err := database.GetAllPatients()
    if err != nil {
        utils.SendJSONResponse(w, http.StatusInternalServerError, Response{
            Success: false,
            Message: "Error fetching patients",
        })
        return
    }

    utils.SendJSONResponse(w, http.StatusOK, Response{
        Success: true,
        Data:    patients,
    })
}

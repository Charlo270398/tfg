package routes

import (
	"encoding/json"
	"net/http"

	models "../models"
)

//GET

func getTagsListHandler(w http.ResponseWriter, req *http.Request) {
	tagsListJSON, err := models.GetTagsList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(tagsListJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

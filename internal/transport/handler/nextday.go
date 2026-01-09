package handler

import (
	"net/http"
	"time"

	"github.com/Vafeda/TODO-List/internal/models"
	"github.com/Vafeda/TODO-List/internal/nextdate"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	date, err := nextdate.NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(date))
}

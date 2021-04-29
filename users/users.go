package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/d-ops-pro/otus-users/domain"
)

func ConfigureRouter(r *mux.Router, db *gorm.DB) {
	r.Path("/user").Methods(http.MethodPost).HandlerFunc(CreateHandler(db))
	r.Path("/user/{id}").Methods(http.MethodPut).HandlerFunc(UpdateHandler(db))
	r.Path("/user/{id}").Methods(http.MethodDelete).HandlerFunc(DeleteHandler(db))
	r.Path("/user/{id}").Methods(http.MethodGet).HandlerFunc(GetHandler(db))
}

func CreateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(domain.User)

		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("unable to decode the request body"))
			return
		}

		err = db.Create(user).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unable to save the user"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func UpdateHandler(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idRaw := vars["id"]

		id, err := strconv.Atoi(idRaw)
		if err != nil {
			logrus.WithError(err).Error("failed to parse the user id to int")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("provided id is not an int type"))
			return
		}

		user := &domain.User{ID: id}
		err = json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("unable to decode the request body"))
			return
		}

		err = db.Save(user).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unable to save the user"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	})
}

func GetHandler(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idRaw := vars["id"]

		id, err := strconv.Atoi(idRaw)
		if err != nil {
			logrus.WithError(err).Error("failed to parse the user id to int")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("provided id is not an int type"))
			return
		}

		user := new(domain.User)
		err = db.Take(user, "id=?", id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("user not found"))
				return
			}

			logrus.WithError(err).Error("failed to get the user by id")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			logrus.WithError(err).Error("unable to encode the user")
		}
	})
}

func DeleteHandler(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idRaw := vars["id"]

		id, err := strconv.Atoi(idRaw)
		if err != nil {
			logrus.WithError(err).Error("failed to parse the user id to int")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("provided id is not an int type"))
			return
		}

		err = db.Delete(domain.User{}, "id=?", id).Error
		if err != nil {
			logrus.WithError(err).Error("unable to delete the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

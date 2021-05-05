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
		logger := logrus.WithField("operation", "user_create")

		user := new(domain.User)
		{
			logger.Infof("decoding request body")
			err := json.NewDecoder(r.Body).Decode(user)
			if err != nil {
				logger.WithError(err).Error("unable to decode the request body")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("unable to decode the request body"))
				return
			}
		}

		logger.Infof("creating the user entity: %+v", user)
		err := db.Create(user).Error
		if err != nil {
			logger.WithError(err).Error("failed to create the user entity")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unable to save the user"))
			return
		}

		logger.Infof("user successfully created with id: %d", user.ID)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func UpdateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idRaw := vars["id"]
		logger := logrus.WithField("operation", "user_update").WithField("user_id", idRaw)

		logger.Infof("parsing the user id as INT")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			logger.WithError(err).Error("failed to parse the user id to int")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("provided id is not an int type"))
			return
		}

		user := &domain.User{ID: id}
		logger.Infof("checking the user exists by provided id")
		err = db.Take(user, "id=?", id).Error
		if err != nil {
			logger.WithError(err).Error("unable to find the user by provided id")

			if errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("user not found \n"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("internal server error \n"))
			return
		}

		logger.Info("user for update found")
		err = json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			logger.WithError(err).Error("unable to decode the request body")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("unable to decode the request body"))
			return
		}

		logger.Infof("updating the user with a new data: %+v", user)
		err = db.Save(user).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("unable to save the user"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func GetHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idRaw := vars["id"]
		logger := logrus.WithField("operation", "user_get").WithField("user_id", idRaw)

		logger.Infof("parsing the user id as INT")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			logger.WithError(err).Error("failed to parse the user id to int")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("provided id is not an int type"))
			return
		}

		logger.Infof("searching the user")
		user := new(domain.User)
		err = db.Take(user, "id=?", id).Error
		if err != nil {
			logger.WithError(err).Error("unable to get the user by ID")
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("user not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("user successfully found by id")
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			logrus.WithError(err).Error("unable to encode the user")
		}
	}
}

func DeleteHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idRaw := vars["id"]
		logger := logrus.WithField("operation", "user_delete").WithField("user_id", idRaw)

		logger.Info("parsing the user ID as INT")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			logger.WithError(err).Error("failed to parse the user id to int")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("provided id is not an int type"))
			return
		}

		err = db.Delete(domain.User{}, "id=?", id).Error
		if err != nil {
			logger.WithError(err).Error("unable to delete the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("user successfully deleted")
		w.WriteHeader(http.StatusOK)
	}
}

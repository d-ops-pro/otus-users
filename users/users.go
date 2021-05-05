package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	libhttp "github.com/d-ops-pro/otus-users/lib/http"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/d-ops-pro/otus-users/domain"
)

type userIDContextKey struct{}

func getUserID(r *http.Request) int {
	id := r.Context().Value(userIDContextKey{})
	return id.(int)
}

func userIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := libhttp.GetLogger(r)

		rawID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(rawID)
		if err != nil {
			logger.WithError(err).Error("failed to parse the user id to int")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("provided id is not an int type"))
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ConfigureRouter(r chi.Router, db *gorm.DB) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/", CreateHandler(db))

		r.Route("/{id}", func(r chi.Router) {
			r.Use(userIDCtx)

			r.Get("/", GetHandler(db))
			r.Put("/", UpdateHandler(db))
			r.Delete("/", DeleteHandler(db))
		})

	})
}

func CreateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := libhttp.GetLogger(r)
		logger = logrus.WithField("operation", "user_create")

		logger.Infof("decoding request body")
		user := new(domain.User)
		{
			err := json.NewDecoder(r.Body).Decode(user)
			if err != nil {
				logger.WithError(err).Error("unable to decode the request body")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("unable to decode the request body"))
				return
			}
		}

		logger.Infof("creating the user entity: %+v", user)
		{
			err := db.Create(user).Error
			if err != nil {
				logger.WithError(err).Error("failed to create the user entity")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("unable to save the user"))
				return
			}
		}

		logger.Infof("user successfully created with id: %d", user.ID)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func UpdateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := libhttp.GetLogger(r)
		userID := getUserID(r)
		logger = logrus.WithField("operation", "user_update").WithField("user_id", userID)

		user := &domain.User{ID: userID}
		logger.Infof("checking the user exists by provided id")
		{
			err := db.Take(user, "id=?", userID).Error
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
		}

		logger.Info("user for update found, decoding the request body")
		{
			err := json.NewDecoder(r.Body).Decode(user)
			if err != nil {
				logger.WithError(err).Error("unable to decode the request body")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("unable to decode the request body"))
				return
			}
		}

		logger.Infof("updating the user with a new data: %+v", user)
		{
			err := db.Save(user).Error
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("unable to save the user"))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func GetHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := libhttp.GetLogger(r)
		userID := getUserID(r)
		logger = logrus.WithField("operation", "user_get").WithField("user_id", userID)

		logger.Infof("searching the user")
		user := new(domain.User)
		err := db.Take(user, "id=?", userID).Error
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
		logger := libhttp.GetLogger(r)
		userID := getUserID(r)
		logger = logrus.WithField("operation", "user_delete").WithField("user_id", userID)

		err := db.Delete(domain.User{}, "id=?", userID).Error
		if err != nil {
			logger.WithError(err).Error("unable to delete the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("user successfully deleted")
		w.WriteHeader(http.StatusOK)
	}
}

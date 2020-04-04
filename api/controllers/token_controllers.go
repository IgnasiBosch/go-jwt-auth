package controllers

import (
	"encoding/json"
	"errors"
	"github.com/IgnasiBosch/go-jwt-auth/api/auth"
	"github.com/IgnasiBosch/go-jwt-auth/api/responses"
	"github.com/IgnasiBosch/go-jwt-auth/api/utils/formaterror"
	"github.com/IgnasiBosch/go-jwt-auth/models"
	"io/ioutil"
	"net/http"
)

func (s *Server) Token(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.JSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	currentUser, err := s.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.JSONError(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	accessToken, err := auth.CreateAccessToken(currentUser.ID)
	refreshToken, err := auth.CreateRefreshToken(currentUser.ID)

	err = responses.JSON(w, http.StatusOK, responses.JWT{AccessToken: accessToken, RefreshToken: refreshToken})
	if err != nil {
		responses.JSONError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) SignIn(email, password string) (*models.User, error) {
	user := models.User{}

	emailErr := s.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	passErr := models.VerifyPassword(user.Password, password)
	if emailErr != nil || passErr != nil {
		return &models.User{}, errors.New("invalid credentials")
	}

	return &user, nil
}

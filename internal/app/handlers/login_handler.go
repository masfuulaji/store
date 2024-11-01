package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
	"github.com/masfuulaji/store/internal/app/repositories"
	"github.com/masfuulaji/store/internal/utils"
)

type LoginHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type LoginHandlerImpl struct {
	userRepository *repositories.UserRepositoryImpl
}

func NewLoginHandler(db *sqlx.DB) *LoginHandlerImpl {
	return &LoginHandlerImpl{userRepository: repositories.NewUserRepository(db)}
}

type Claims struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (u *LoginHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := u.userRepository.GetUserByUsername(user.Username)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Password != res.Password {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(10 * time.Hour)
	claim := &Claims{
		Id:       res.Id,
		Username: res.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Login successful"})
}

func (u *LoginHandlerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Logout successful"})
}

func (u *LoginHandlerImpl) IsLogin(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "You are logged in"})
}

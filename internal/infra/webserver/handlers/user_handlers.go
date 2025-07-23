package handlers

import (
	"encoding/json"
	"github.com/go-chi/jwtauth"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/dto"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/entity"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/infra/database"
	"net/http"
	"time"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// GetJWt godoc
// @Summary      Get a user JWT
// @Description  Get a user JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.GetJWTInput  true  "user credentials"
// @Success      200 {object}  dto.GetJWTOutput
// @Failure      404 {object}  Error
// @Failure      500 {object}  Error
// @Router       /users/generate-token [post]
func (h *UserHandler) GetJWt(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)
	var user dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})
	accessToken := dto.GetJWTOutput{AccessToken: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// CreateUser godoc
// @Summary      Create User
// @Description  Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      500 {object}  Error
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

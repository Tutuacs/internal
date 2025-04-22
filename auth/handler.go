package auth

import (
	"fmt"
	"net/http"

	"github.com/Tutuacs/pkg/enums"
	jwtService "github.com/Tutuacs/pkg/jwt"
	"github.com/Tutuacs/pkg/mail"
	"github.com/Tutuacs/pkg/password"
	"github.com/Tutuacs/pkg/resolver"
	"github.com/Tutuacs/pkg/routes"
	"github.com/Tutuacs/pkg/types"
	"github.com/dgrijalva/jwt-go"
)

type Handler struct {
	subRoute string
	s        *Store
}

func NewHandler(s *Store) *Handler {
	return &Handler{
		subRoute: "/auth",
		s:        s,
	}
}

func (h *Handler) BuildRoutes(router routes.Route) {
	// TODO implement the routes call
	router.NewRoute(routes.POST, h.subRoute+"/login", h.login)
	router.NewRoute(routes.POST, h.subRoute+"/register", h.register)
	router.NewRoute(routes.POST, h.subRoute+"/forgot-request", h.forgotRequest)
	router.NewRoute(routes.POST, h.subRoute+"/forgot-update", h.forgotUpdate)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var data types.LoginDto

	if err := resolver.GetBody(r, &data); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	store := h.s

	usr, err := store.GetLogin(data.Email)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error": "Unable to retrieve user."})
		return
	}

	if !password.ValidPassword(usr.Password, data.Password) {
		resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": "Invalid credentials."})
		return
	}

	token, err := jwtService.CreateJWT(usr.Email, usr.ID, usr.Role)
	if err != nil {
		resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": fmt.Sprintf("Error creating token: %s", err)})
		return
	}

	userMap := make(map[string]any)

	userMap["id"] = usr.ID
	userMap["name"] = usr.Name
	userMap["email"] = usr.Email
	userMap["role"] = usr.Role

	resolver.WriteResponse(w, http.StatusOK, map[string]any{"token": token, "user": userMap})

}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserDto

	if err := resolver.GetBody(r, &payload); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	store := h.s

	// Check if user already exists
	_, err := store.GetUserByEmail(payload.Email)
	if err == nil {
		resolver.WriteResponse(w, http.StatusConflict, map[string]string{"Error": "User already exists."})
		return
	}

	hashedPassword, err := password.HashPassword(payload.Password)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Create new user
	err = store.CreateUser(types.User{
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     enums.ROLE_CLIENT,
	})
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	resolver.WriteResponse(w, http.StatusCreated, map[string]string{"Ok": "Created"})
}

func (h *Handler) forgotRequest(w http.ResponseWriter, r *http.Request) {
	var payload types.RequestForgotPassword

	if err := resolver.GetBody(r, &payload); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	store := h.s

	usr, err := store.GetUserByEmail(payload.Email)
	if err != nil {
		resolver.WriteResponse(w, http.StatusNotFound, map[string]string{"Error": "User not found."})
		return
	}

	token, err := jwtService.CreateForgotToken(usr.Email)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error": "Error creating token."})
		return
	}

	mailer := mail.UseMailer()

	err = mailer.SendForgotEmail(usr.Email, token)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error": fmt.Sprintf("Ocorreu um erro ao enviar um email para seu email cadastrado: %s", err.Error())})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, map[string]string{"Ok": "Email sent."})
}

func (h *Handler) forgotUpdate(w http.ResponseWriter, r *http.Request) {
	var payload types.ForgotUpdatePassword

	if err := resolver.GetBody(r, &payload); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	token, err := jwtService.ValidateJWT(payload.Token)
	if err != nil {
		resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": "Invalid token."})
		return
	}

	tokenEmail := token.Claims.(jwt.MapClaims)["email"].(string)

	store := h.s

	hashedPassword, err := password.HashPassword(payload.NewPassword)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error": "Error hashing password: " + err.Error()})
		return
	}

	err = store.UpdatePassword(tokenEmail, hashedPassword)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error": "Error updating password: " + err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, map[string]string{"Ok": "Password updated."})
}

// ! Recommended private functions
// * Create stores to get DB data like this
/*
	store, err := NewStore()
	if err != nil {
		return
	}

	* Use resolver to getParams, getBody and writeResponse

*/

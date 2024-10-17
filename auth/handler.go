package auth

import (
	"fmt"
	"net/http"

	"github.com/Tutuacs/pkg/enums"
	"github.com/Tutuacs/pkg/jwt"
	"github.com/Tutuacs/pkg/password"
	"github.com/Tutuacs/pkg/resolver"
	"github.com/Tutuacs/pkg/routes"
	"github.com/Tutuacs/pkg/types"
)

type Handler struct {
	subRoute string
}

func NewHandler() *Handler {
	return &Handler{subRoute: "/auth"}
}

func (h *Handler) BuildRoutes(router routes.Route) {
	// TODO implement the routes call
	router.NewRoute(routes.POST, h.subRoute+"/login", h.login)
	router.NewRoute(routes.POST, h.subRoute+"/register", h.register)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var data types.LoginDto

	if err := resolver.GetBody(r, &data); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	store, err := NewStore()
	if err != nil {
		return
	}

	defer store.CloseStore()

	usr, err := store.GetLogin(data.Email)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error": "Unable to retrieve user."})
		return
	}

	if !password.ValidPassword(usr.Password, data.Password) {
		resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": "Invalid credentials."})
		return
	}

	token, err := jwt.CreateJWT(usr.Email, usr.ID, usr.Role)
	if err != nil {
		resolver.WriteResponse(w, http.StatusUnauthorized, map[string]string{"Error": fmt.Sprintf("Error creating token: %s", err)})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, map[string]string{"token": token})

}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserDto

	if err := resolver.GetBody(r, &payload); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	store, err := NewStore()
	if err != nil {
		return
	}

	defer store.CloseStore()

	// Check if user already exists
	_, err = store.GetUserByEmail(payload.Email)
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

// ! Recommended private functions
// * Create stores to get DB data like this
/*
	store, err := NewStore()
	if err != nil {
		return
	}

	defer store.CloseStore()

	* Use resolver to getParams, getBody and writeResponse

*/

package user

import (
	"net/http"
	"strconv"

	"github.com/Tutuacs/pkg/enums"
	"github.com/Tutuacs/pkg/guards"
	"github.com/Tutuacs/pkg/password"
	"github.com/Tutuacs/pkg/resolver"
	"github.com/Tutuacs/pkg/routes"
	"github.com/Tutuacs/pkg/types"
	"github.com/go-playground/validator"
)

type Handler struct {
	subRoute string
}

func NewHandler() *Handler {
	return &Handler{subRoute: "/user"}
}

func (h *Handler) BuildRoutes(router routes.Route) {
	// TODO implement the routes call
	router.NewRoute(routes.POST, h.subRoute, guards.AutenticatedRoute(h.create, enums.ROLE_ADMIN))
	router.NewRoute(routes.GET, h.subRoute, guards.AutenticatedRoute(h.list, enums.ROLE_CLIENT, enums.ROLE_ADMIN))
	router.NewRoute(routes.GET, h.subRoute+"/{id}", guards.AutenticatedRoute(h.getById, enums.ROLE_CLIENT, enums.ROLE_ADMIN))
	router.NewRoute(routes.PUT, h.subRoute+"/{id}", guards.AutenticatedRoute(h.update, enums.ROLE_CLIENT, enums.ROLE_ADMIN))
	router.NewRoute(routes.DELETE, h.subRoute+"/{id}", guards.AutenticatedRoute(h.delete, enums.ROLE_CLIENT, enums.ROLE_ADMIN))
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {

	var payload types.NewUserDto

	// resolver.GetBody(r, &payload)

	if err := resolver.GetBody(r, &payload); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := resolver.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		resolver.WriteResponse(w, http.StatusBadRequest, errors.Error())
		return
	}

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	// TODO: Implement the auth validation after create

	exists, err := store.GetByEmail(payload.Email)
	if err == nil && exists.ID != 0 {
		resolver.WriteResponse(w, http.StatusConflict, map[string]string{"Error creating the user ": "User already exists"})
		return
	}

	pass, err := password.HashPassword(payload.Password)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error encrypting the users pass": err.Error()})
		return
	}
	payload.Password = pass

	usr, err := store.Create(payload)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error creating the user": err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusCreated, usr)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {

	userLogged := r.Context().Value(guards.UserKey).(*types.User)

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	if userLogged.Role == enums.ROLE_CLIENT {
		user, err := store.GetByID(userLogged.ID)
		if err != nil {
			resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error getting the user": err.Error()})
			return
		}

		resolver.WriteResponse(w, http.StatusOK, user)
		return
	}

	users, err := store.List()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error listing the users": err.Error()})
		return
	}

	if users == nil {
		resolver.WriteResponse(w, http.StatusOK, []types.User{})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, users)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {

	userLogged := r.Context().Value(guards.UserKey).(*types.User)
	param := resolver.GetParam(r, "id")

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		if err != nil {
			resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": err.Error()})
		}
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": "id is invalid"})
		return
	}

	if (userLogged.Role == enums.ROLE_CLIENT) && (userLogged.ID != id) {
		resolver.WriteResponse(w, http.StatusForbidden, map[string]string{"Error": "You are not allowed to see this user"})
		return
	}

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	user, err := store.GetByID(id)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error getting the user": err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, user)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {

	userLogged := r.Context().Value(guards.UserKey).(*types.User)
	param := resolver.GetParam(r, "id")

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		if err != nil {
			resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": err.Error()})
		}
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": "id is invalid"})
		return
	}

	var body types.UpdateUserDto
	if err := resolver.GetBody(r, body); err != nil {
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error getting the body": err.Error()})
		return
	}

	if err := resolver.Validate.Struct(body); err != nil {
		errors := err.(validator.ValidationErrors)
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error validating the body": errors.Error()})
		return
	}

	if (userLogged.Role == enums.ROLE_CLIENT) && (userLogged.ID != id) || (userLogged.Role == enums.ROLE_CLIENT) && (body.Role != enums.ROLE_CLIENT) {
		resolver.WriteResponse(w, http.StatusForbidden, map[string]string{"Error": "You are not allowed to update this users role"})
		return
	}

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	updated, err := store.Update(id, body)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error updating the user": err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, updated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {

	userLogged := r.Context().Value(guards.UserKey).(*types.User)
	param := resolver.GetParam(r, "id")

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		if err != nil {
			resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": err.Error()})
		}
		resolver.WriteResponse(w, http.StatusBadRequest, map[string]string{"Error parsing the id": "id is invalid"})
		return
	}

	if (userLogged.Role == enums.ROLE_CLIENT) && (userLogged.ID != id) {
		resolver.WriteResponse(w, http.StatusForbidden, map[string]string{"Error": "You are not allowed to see this user"})
		return
	}

	store, err := NewStore()
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error opening the store": err.Error()})
		return
	}

	defer store.CloseStore()

	deleted, err := store.Delete(id)
	if err != nil {
		resolver.WriteResponse(w, http.StatusInternalServerError, map[string]string{"Error deleting the user": err.Error()})
		return
	}

	resolver.WriteResponse(w, http.StatusOK, deleted)
}

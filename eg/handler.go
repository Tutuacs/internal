package eg

import (
	"net/http"

	"github.com/Tutuacs/pkg/routes"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) BuildRoutes(router routes.Route) {
	router.NewRoute(routes.GET, "/eg", h.getRoute)
}

func (h *Handler) getRoute(w http.ResponseWriter, r *http.Request) {

	store, err := NewStore()
	if err != nil {
		return
	}

	defer store.CloseStore()

	/*

	* Use resolver to getParams, getBody and writeResponse

	 */
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"fijoy/.gen/neondb/public/model"
	"net/http"

	. "fijoy/.gen/neondb/public/table"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type userHandler struct {
	tokenAuth *jwtauth.JWTAuth
	db        *sql.DB
}

func NewUserHandler(tokenAuth *jwtauth.JWTAuth, db *sql.DB) chi.Router {
	handler := userHandler{tokenAuth, db}

	router := chi.NewRouter()
	router.Use(jwtauth.Verifier(tokenAuth))
	router.Use(jwtauth.Authenticator(tokenAuth))

	router.Get("/", handler.GetUserData)
	return router
}

func (uh userHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	stmt := SELECT(FijoyUser.AllColumns).FROM(FijoyUser).
		WHERE(FijoyUser.ID.EQ(String(claims["user_id"].(string))))

	dest := model.FijoyUser{}

	err := stmt.QueryContext(r.Context(), uh.db, &dest)
	if err != nil {
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(dest)
}

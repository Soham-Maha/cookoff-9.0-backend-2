package controllers

import (
	"fmt"
	"net/http"

	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/go-chi/jwtauth/v5"
)

func RoleFromToken(w http.ResponseWriter, r *http.Request, user string) bool {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		httphelpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return false
	}
	role, ok := claims["role"].(string)
	if !ok {
		httphelpers.WriteError(w, http.StatusUnauthorized, "Role not found in token")
		return false
	}
	if role != user {
		msg := fmt.Sprintf("Access Denied: %s not allowed", role)
		httphelpers.WriteError(w, http.StatusForbidden, msg)
		return false
	}
	return true
}
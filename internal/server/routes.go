package server

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/controllers"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hibiken/asynq"
)

func (s *Server) RegisterRoutes(taskClient *asynq.Client) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", controllers.HealthCheck)
	r.Put("/callback", func(w http.ResponseWriter, r *http.Request) {
		controllers.CallbackUrl(w, r, taskClient)
	})
	
	r.Post("/login/user", controllers.LoginHandler)
	r.Post("/token/refresh", controllers.RefreshTokenHandler)
	r.Group(func(protected chi.Router) {
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Post("/upgrade", controllers.UpgradeUserToRound)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Post("/roast", controllers.BanUsers)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Post("/unroast", controllers.UnbanUsers)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Post("/round/{round_number}/enable", controllers.EnableRound)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Post("/round/{round_number}/disable", controllers.DisableRound)
		protected.Use(jwtauth.Verifier(auth.TokenAuth))
		protected.Use(jwtauth.Authenticator(auth.TokenAuth))
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Get("/users", controllers.GetAllUsers)
		protected.Get("/protected", controllers.ProtectedHandler)
		protected.With(middlewares.BanCheckMiddleware).Post("/submit", controllers.SubmitCode)
        protected.With(middlewares.BanCheckMiddleware).Post("/runcode", controllers.RunCode)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Post("/question/create", controllers.CreateQuestion)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Get("/questions", controllers.GetAllQuestion)
		protected.With(middlewares.BanCheckMiddleware).Get("/question/round", controllers.GetQuestionsByRound)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Get("/question/{question_id}", controllers.GetQuestionById)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Delete("/question/{question_id}", controllers.DeleteQuestion)
		protected.With(middlewares.RoleAuthorizationMiddleware("admin")).Patch("/question", controllers.UpdateQuestion)
	})

	return r
}

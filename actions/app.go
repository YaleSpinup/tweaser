package actions

import (
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/unrolled/secure"

	"github.com/YaleUniversity/tweaser/models"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var AdminToken = envy.Get("ADMIN_TOKEN", "")
var CryptToken = envy.Get("CRYPT_TOKEN", "")

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_tweaser_session",
		})
		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		app.Use(middleware.PopTransaction(models.DB))
		app.GET("/v1/tweaser/ping", PingPong)

		userAPI := app.Group("/v1/tweaser")
		userAPI.POST("/responses", ResponsesCreate)

		adminAPI := app.Group("/v1/tweaser/admin")
		adminAPI.Use(sharedTokenAuth)

		adminAPI.GET("/campaigns", CampaignsList)
		adminAPI.POST("/campaigns", CampaignsCreate)
		adminAPI.GET("/campaigns/{campaign_id}", CampaignsGet)
		adminAPI.PUT("/campaigns/{campaign_id}", CampaignsUpdate)

		adminAPI.GET("/questions", QuestionsList)
		adminAPI.GET("/questions/{question_id}", QuestionsGet)
		adminAPI.POST("/questions", QuestionsCreate)
		adminAPI.PUT("/questions/{question_id}", QuestionsUpdate)

		adminAPI.GET("/answers", AnswersList)
		adminAPI.GET("/answers/{answer_id}", AnswersGet)
		adminAPI.POST("/answers", AnswersCreate)
		adminAPI.PUT("/answers/{answer_id}", AnswersUpdate)

		adminAPI.GET("/responses", ResponsesList)
		adminAPI.GET("/responses/{response_id}", ResponsesGet)
	}

	return app
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return ssl.ForceSSL(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func sharedTokenAuth(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		headers, ok := c.Request().Header["X-Auth-Token"]
		if !ok || len(headers) == 0 || headers[0] != AdminToken {
			log.Println("Missing or bad token header for request", c.Request().URL)
			return c.Error(401, errors.New("Unauthorized!"))
		}
		log.Println("Auth looks good")
		return next(c)
	}
}

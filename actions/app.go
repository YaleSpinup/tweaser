package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"git.yale.edu/spinup/tweaser/models"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

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

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))
		api := app.Group("/v1/tweaser")
		api.GET("/ping", PingPong)
		api.GET("/campaigns", CampaignsList)
		api.POST("/campaigns", CampaignsCreate)
		api.PUT("/campaigns/{campaign_id}", CampaignsUpdate)
		api.GET("/campaigns/{campaign_id}", CampaignsGet)
		api.GET("/campaigns/{campaign_id}/questions", CampaignsGetQuestions)

		api.GET("/questions", QuestionsList)
		api.GET("/questions/{question_id}", QuestionsGet)
		api.POST("/questions", QuestionsCreate)
		api.PUT("/questions/{question_id}", QuestionsUpdate)

		api.GET("/answers", AnswersList)
		api.GET("/answers/{answer_id}", AnswersGet)
		api.POST("/answers", AnswersCreate)
		api.PUT("/answers/{answer_id}", AnswersUpdate)
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

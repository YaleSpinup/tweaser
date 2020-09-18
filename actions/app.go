package actions

import (
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/pkg/errors"
	"github.com/unrolled/secure"

	"github.com/YaleSpinup/tweaser/models"
	"github.com/YaleSpinup/tweaser/tweaser"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

var (
	app *buffalo.App

	// ENV is used to help switch settings based on where the
	// application is being run. Default is "development".
	ENV        = envy.Get("GO_ENV", "development")
	AdminToken = envy.Get("ADMIN_TOKEN", "")
	CryptToken = envy.Get("CRYPT_TOKEN", "")

	// Version is the main version number
	Version = tweaser.Version

	// VersionPrerelease is a prerelease marker
	VersionPrerelease = tweaser.VersionPrerelease

	// buildstamp is the timestamp the binary was built, it should be set at buildtime with ldflags
	buildstamp = tweaser.BuildStamp

	// githash is the git sha of the built binary, it should be set at buildtime with ldflags
	githash = tweaser.GitHash
)

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

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		app.Use(popmw.Transaction(models.DB))
		app.GET("/v1/tweaser/ping", PingPong)
		app.GET("/v1/tweaser/version", VersionHandler)

		userAPI := app.Group("/v1/tweaser")
		userAPI.POST("/responses", ResponsesCreate)

		adminAPI := app.Group("/v1/tweaser/admin")
		adminAPI.Use(sharedTokenAuth)

		adminAPI.GET("/campaigns", CampaignsList)
		adminAPI.POST("/campaigns", CampaignsCreate)
		adminAPI.GET("/campaigns/{campaign_id}", CampaignsGet)
		adminAPI.PUT("/campaigns/{campaign_id}", CampaignsUpdate)
		adminAPI.GET("/campaigns/{campaign_id}/questions", CampaignsGetQuestions)

		adminAPI.GET("/questions", QuestionsList)
		adminAPI.GET("/questions/{question_id}", QuestionsGet)
		adminAPI.POST("/questions", QuestionsCreate)
		adminAPI.PUT("/questions/{question_id}", QuestionsUpdate)
		adminAPI.GET("/questions/{question_id}/answers", QuestionsGetAnswers)
		adminAPI.GET("/questions/{question_id}/responses", QuestionsGetResponses)

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
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func sharedTokenAuth(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		headers, ok := c.Request().Header["X-Auth-Token"]
		if !ok || len(headers) == 0 || headers[0] != AdminToken {
			log.Println("Missing or bad token header for request", c.Request().URL)
			return c.Error(403, errors.New("Forbidden!"))
		}
		log.Println("Auth looks good")
		return next(c)
	}
}

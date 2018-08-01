package actions

import "github.com/gobuffalo/buffalo"

// PingPong default implementation.
func PingPong(c buffalo.Context) error {
	return c.Render(200, r.String("pong"))
}

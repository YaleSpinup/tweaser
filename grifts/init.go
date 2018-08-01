package grifts

import (
	"git.yale.edu/spinup/tweaser/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}

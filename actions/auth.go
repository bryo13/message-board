package actions

import (
	"fmt"
	"message_board/models"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/pkg/errors"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}

	// db context
	tx := c.Value("tx").(*pop.Connection)

	// find user
	query := tx.Where("provider=? and provider_id=?", user.Provider, user.UserID)
	exists, err := query.Exists("users")
	if err != nil {
		return errors.WithStack(err)
	}

	u := &models.User{}
	if exists {
		err = query.First(u)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		u.Name = user.Name
		u.Provider = user.Provider
		u.ProviderID = user.UserID
		err = tx.Save(u)
		if err != nil {
			return errors.WithStack(err)
		}

	}

	c.Session().Set("current_user_id", u.ID)
	err = c.Session().Save()
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(302, r.JSON("success signed up!"))
}

func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	if err := c.Session().Save(); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(302, r.JSON("Logged out"))
}

func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("Danger", "You must be authorized")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}

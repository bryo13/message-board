package actions

import (
	"message_board/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// MessagesAll default implementation.
func MessagesAll(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("Cant get db transaction")
	}
	messages := &models.Messages{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Blogs from the DB
	if err := q.All(messages); err != nil {
		return err
	}
	return c.Render(http.StatusOK, r.JSON(messages))
}

// MessagesCreate default implementation.
func MessagesCreate(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return c.Render(400, r.JSON("no db connection found"))
	}

	message := &models.Message{}
	if err := c.Bind(message); err != nil {
		return errors.WithStack(err)
	}

	u := c.Value("current_user").(*models.User)
	if u.ID == uuid.Nil {
		return errors.New("Please login")
	}

	message.User = u

	verrs, err := tx.Eager().ValidateAndCreate(message)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return errors.WithStack(verrs)
	}

	return c.Render(200, r.JSON(message))
}

// SetContentType on request
func SetContentType(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Set("Content-Type", "application/json")
		return next(c)
	}
}

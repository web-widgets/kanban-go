package api

import (
	"context"
	"errors"
	"net/http"

	remote "github.com/mkozhukh/go-remote"

	"web-widgets/kanban-go/data"
)

type UserID int
type DeviceID int

type CardEvent struct {
	Type   string     `json:"type"`
	From   int        `json:"-"`
	Card   *data.Card `json:"card"`
	Before int        `json:"before,omitempty"`
}

type ColumnEvent struct {
	Type   string       `json:"type"`
	From   int			`json:"-"`
	Column *data.Column `json:"column"`
	Before int        `json:"before,omitempty"`
}

type RowEvent struct {
	Type   string `json:"type"`
	From   int	  `json:"-"`
	Row *data.Row `json:"row"`
	Before int        `json:"before,omitempty"`
}


func BuildAPI(db *data.DAO) *remote.Server {
	if remote.MaxSocketMessageSize < 32000 {
		remote.MaxSocketMessageSize = 32000
	}

	api := remote.NewServer(&remote.ServerConfig{
		WebSocket: true,
	})

	api.Events.AddGuard("cards", func(m *remote.Message, c *remote.Client) bool {
		tm, ok := m.Content.(CardEvent)
		if !ok {
			return false
		}

		return int(tm.From) != c.ConnID
	})

	api.Events.AddGuard("columns", func(m *remote.Message, c *remote.Client) bool {
		tm, ok := m.Content.(ColumnEvent)
		if !ok {
			return false
		}

		return int(tm.From) != c.ConnID
	})

	api.Events.AddGuard("rows", func(m *remote.Message, c *remote.Client) bool {
		tm, ok := m.Content.(RowEvent)
		if !ok {
			return false
		}

		return int(tm.From) != c.ConnID
	})

	api.Connect = func(r *http.Request) (context.Context, error) {
		id, _ := r.Context().Value("user_id").(int)
		if id == 0 {
			return nil, errors.New("access denied")
		}
		device, _ := r.Context().Value("device_id").(int)
		if device == 0 {
			return nil, errors.New("access denied")
		}

		return context.WithValue(
			context.WithValue(r.Context(), remote.UserValue, id),
			remote.ConnectionValue, device), nil
	}

	return api
}

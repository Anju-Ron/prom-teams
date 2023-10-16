package card

import (
	"context"
    "time"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/alertmanager/notify/webhook"
)

// Office365ConnectorCard represents https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-reference#example-office-365-connector-card
type Office365ConnectorCard struct {
	Type    string        `json:"type"`
	Body    []interface{} `json:"body"`
	Actions []Action      `json:"actions"`
	Schema  string        `json:"$schema"`
	Version string        `json:"version"`
}

type Image struct {
	Type   string `json:"type"`
	URL    string `json:"url"`
	Height string `json:"height"`
	Width  string `json:"width"`
}

type TextBlock struct {
	Type                string `json:"type"`
	Text                string `json:"text"`
	Size                string `json:"size"`
	Weight              string `json:"weight"`
	HorizontalAlignment string `json:"horizontalAlignment"`
	Color               string `json:"color"`
	Wrap                bool   `json:"wrap"`
}

type Column struct {
	Type  string        `json:"type"`
	Width string        `json:"width"`
	Items []interface{} `json:"items"`
}

type ColumnSet struct {
	Type    string   `json:"type"`
	Columns []Column `json:"columns"`
}

type Container struct {
	Type  string        `json:"type"`
	Style string        `json:"style"`
	Items []interface{} `json:"items"`
}

type Action struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Converter converts an alert manager webhook message to Office365ConnectorCard.
type Converter interface {
	Convert(context.Context, webhook.Message) (Office365ConnectorCard, error)
}

type loggingMiddleware struct {
	logger log.Logger
	next   Converter
}
// NewCreatorLoggingMiddleware creates a loggingMiddleware.
func NewCreatorLoggingMiddleware(l log.Logger, n Converter) Converter {
	return loggingMiddleware{l, n}
}

func (l loggingMiddleware) Convert(ctx context.Context, a webhook.Message) (c Office365ConnectorCard, err error) {
	defer func(begin time.Time) {
		if len(c.Actions) > 5 {
			l.logger.Log(
				"warning", "There can only be a maximum of 5 actions in a potentialAction collection",
				"actions", c.Actions,
			)
		}

		
		l.logger.Log(
			"alert", a,
			"card", c,
			"took", time.Since(begin),
		)
	}(time.Now())
	return l.next.Convert(ctx, a)
}

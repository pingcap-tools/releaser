package manager

import (
	"github.com/juju/errors"
	"github.com/nlopes/slack"
)

func initSlackClient(token string) *slack.Client {
	if token == "" {
		return nil
	}
	return slack.New(token)
}

// SendMessage ...
func (m *Manager) SendMessage(message string) error {
	if m.Slack == nil {
		return nil
	}
	_, _, err := m.Slack.PostMessage(m.Config.SlackChannel,
		slack.MsgOptionText(message, true))
	return errors.Trace(err)
}

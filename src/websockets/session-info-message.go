package websockets

const (
	UserJoinEvent  = "user-joined"
	UserLeaveEvent = "user-left"
	ActionEvent    = "action"
)

type sessionInfoMessage struct {
	SessionId   string                 `json:"sessionId"`
	Event       string                 `json:"event"`
	ClientNames []string               `json:"clientNames"`
	Payload     map[string]interface{} `json:"payload"`
}

func newSessionInfoMessage(s *Session, event string, payload map[string]interface{}) *sessionInfoMessage {
	return &sessionInfoMessage{SessionId: s.sessionId, ClientNames: s.getClientNames(), Event: event, Payload: payload}
}

func NewClientJoinedMessage(session *Session, c *Client) *sessionInfoMessage {
	m := newSessionInfoMessage(session, UserJoinEvent, map[string]interface{}{"user": c.Name})
	return m
}

func NewClientLeftMessage(session *Session, c *Client) *sessionInfoMessage {
	m := newSessionInfoMessage(session, UserLeaveEvent, map[string]interface{}{"user": c.Name})
	return m
}
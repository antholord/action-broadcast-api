package websockets

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

// Session maintains the set of active clients and broadcasts messages to the
// clients.
type Manager struct {
	// Registered clients.
	sessions map[string]*Session

	MapLock sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

func (m *Manager) HandleCreate(sessionId string, clientName string, c *gin.Context) {
	session := m.createNewSession(sessionId)
	client, err := m.handleCreateClient(session, clientName, c)
	
	if err != nil {
		return
	}

	session.register <- client
}

func (m *Manager) HandleJoin(sessionId string, clientName string, c *gin.Context) {
	m.MapLock.Lock(); defer m.MapLock.Unlock()

	if session, found := m.sessions[sessionId]; found {
		client, err := m.handleCreateClient(session, clientName, c)
	
		if err != nil {
			return
		}

		session.register <- client
	} //TODO throw error if joining unfound session
}

func (m *Manager) handleCreateClient(session *Session, clientName string, c *gin.Context) (*Client, error) {
	client, err := CreateClient(session, clientName, c.Writer, c.Request)
	
	if err != nil {
		log.Println("Error Creating session and/or client", err)
		c.String(500, "Error Creating session and/or client", err)
		return nil, err
	}
	return client, nil
}

func (m *Manager) createNewSession(sessionId string) *Session {
	m.MapLock.Lock(); defer m.MapLock.Unlock()
	
	if (sessionId == "") {
		sessionId = m.createSessionId()
	}
	
	log.Println("Creating Session ", sessionId)

	session := NewSession(m, sessionId)
	m.sessions[sessionId] = session

	go session.Run()
	return session
}


func (m *Manager) createSessionId() string {
	sid := shortid.MustNew(1, shortid.DefaultABC, 2342)
	id, err := sid.Generate()
	if err != nil {
		log.Println("Error creating session id", err)
	}

	if _, found := m.sessions[id]; found {
		return m.createSessionId()
	} else {
		return id
	}
}
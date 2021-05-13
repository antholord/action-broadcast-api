package websockets

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-diceware/diceware"
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
		sessions:    make(map[string]*Session),
	}
}

func (m *Manager) HandleCreate(c *gin.Context) {
	session := m.createNewSession()
	client, err := CreateClient(session, c.Writer, c.Request)
	
	if err != nil {
		log.Println("Error Creating session and/or client", err)
		c.String(500, "Error Creating session and/or client", err)
		return
	} 

	session.register <- client
	c.String(http.StatusOK, session.sessionId)
}

func (m *Manager) HandleJoin(sessionId string, c *gin.Context) {

}

func (m *Manager) createNewSession() *Session {
	m.MapLock.Lock(); defer m.MapLock.Unlock()
	
	sessionId := m.createSessionId()
	log.Println("Creating Session ", sessionId)
	session := NewSession(sessionId)
	m.sessions[sessionId] = session
	return session
}

func (m *Manager) createSessionId() string {
	list, err := diceware.Generate(1)
	if err != nil {
		log.Fatal(err)
	}
	var sessionId = list[0];
	if _, found := m.sessions[sessionId]; found {
		//add while loop eventually
		list, err := diceware.Generate(1)
		if err != nil {
			log.Fatal(err)
		}
		sessionId = list[0];
	}
	return sessionId
}
package arenaserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	notify "github.com/bitly/go-notify"
	"github.com/bytearena/core/arenaserver/agent"
	"github.com/bytearena/core/arenaserver/comm"
	"github.com/bytearena/core/common/assert"
	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils"
	uuid "github.com/satori/go.uuid"

	bettererrors "github.com/xtuc/better-errors"
)

var (
	LISTEN_ADDR = net.IP{0, 0, 0, 0}
)

func (server *Server) listen() chan interface{} {
	serveraddress := LISTEN_ADDR.String() + ":" + strconv.Itoa(server.port)
	server.commserver = comm.NewCommServer(serveraddress)

	// Consume comm server events
	go func() {
		for {
			msg := <-server.commserver.Events()

			go func() {

				switch t := msg.(type) {
				case comm.EventLog:
					server.Log(EventLog{t.Value})

				case comm.EventWarn:
					server.Log(EventWarn{t.Err})

				case comm.EventError:
					server.Log(EventError{t.Err})

				case comm.EventRawComm:
					server.Log(EventRawComm{
						Value: t.Buffer,
						From:  t.From,
					})

				// An agent has probaly been disconnected
				// We need to remove it from our state
				case comm.EventConnDisconnected:
					server.clearAgentConn(t.Conn)
					server.Log(EventWarn{t.Err})

				default:
					msg := fmt.Sprintf("Unsupported message of type %s", reflect.TypeOf(msg))
					panic(msg)
				}
			}()
		}
	}()

	//server.events <- EventLog{"Server listening on port " + strconv.Itoa(server.port)}

	err := server.commserver.Listen(server)
	utils.Check(err, "Failed to listen on "+serveraddress)

	block := make(chan interface{})
	notify.Start("app:stopticking", block)

	return block
}

func (server *Server) clearAgentConn(conn net.Conn) {
	server.agentproxiesmutex.Lock()

	for k, agentproxy := range server.agentproxies {
		netAgent, ok := agentproxy.(agent.AgentProxyNetworkInterface)

		if ok && netAgent.GetConn() == conn {

			server.clearAgentById(k)
			break
		}

	}

	server.agentproxiesmutex.Unlock()
}

func (server *Server) clearAgentById(k uuid.UUID) {

	// Remove agent from our state
	delete(server.agentproxies, k)
	delete(server.agentimages, k)
	delete(server.agentproxieshandshakes, k)

	server.Log(EventDebug{fmt.Sprintf("Removing %s from state", k.String())})
}

/* <implementing types.AgentCommunicatorInterface> */
func (s *Server) NetSend(message []byte, conn net.Conn) error {
	return s.commserver.Send(message, conn)
}

func (server *Server) PushMutationBatch(batch types.AgentMutationBatch) {
	server.mutationsmutex.Lock()
	server.pendingmutations = append(server.pendingmutations, batch)
	server.mutationsmutex.Unlock()
}

/* </implementing types.AgentCommunicatorInterface> */

/* <implementing types.CommunicatorDispatcherInterface> */
func (server *Server) ImplementsCommDispatcherInterface() {}
func (server *Server) DispatchAgentMessage(msg types.AgentMessage) error {

	agentproxy, err := server.getAgentProxy(msg.GetAgentId().String())
	if err != nil {
		return errors.New("DispatchAgentMessage: agentid does not match any known agent in received agent message !;" + msg.GetAgentId().String())
	}

	// proto := msg.GetEmitterConn().LocalAddr().Network()
	// ip := strings.Split(msg.GetEmitterConn().RemoteAddr().String(), ":")[0]
	// if proto != "tcp" || ip != "TODO(jerome):take from agent container struct"
	// Problem here: cannot check ip against the one we get from Docker by inspecting the container
	// as the two addresses do not match

	assert.Assert(msg.GetMethod() != "", "Method is null")

	switch strings.ToLower(msg.GetMethod()) {
	case types.AgentMessageType.Handshake:
		{
			if _, found := server.agentproxieshandshakes[msg.GetAgentId()]; found {
				return errors.New("ERROR: Received duplicate handshake from agent " + agentproxy.String())
			}

			server.agentproxieshandshakes[msg.GetAgentId()] = struct{}{}

			var handshake types.AgentMessagePayloadHandshake
			err = json.Unmarshal(msg.GetPayload(), &handshake)
			if err != nil {
				return bettererrors.
					New("Failed to unmarshal agent's handshake").
					SetContext("agent", msg.GetAgentId().String())
			}

			ag, ok := agentproxy.(agent.AgentProxyNetworkInterface)
			if !ok {
				return bettererrors.
					New("Failed to cast agent to NetAgent during handshake").
					SetContext("agent", ag.String())
			}

			// Check if the agent uses a protocol version we know
			if handshake.Version == "" {
				handshake.Version = "UNKNOWN"
			}

			if !utils.IsStringInArray(types.PROTOCOL_VERSIONS, handshake.Version) {
				return bettererrors.
					New("Unsupported agent protocol").
					SetContext("agent", ag.String()).
					SetContext("protocol version", handshake.Version)
			}

			ag = ag.SetConn(msg.GetEmitterConn())
			server.setAgentProxy(ag)

			server.events <- EventDebug{"Received handshake from agent " + ag.String() + ""}

			server.nbhandshaked++

			ag.SendAgentWelcome(
				server.GetGame().GetAgentWelcome(ag.GetEntityId()),
				server,
			)

			if server.nbhandshaked == server.getNbExpectedagents() {
				server.onAgentsReady()
			}

			// TODO(sven|jerome): handle some timeout here if all agents fail to handshake

			break
		}
	case types.AgentMessageType.Actions:
		{
			var actionsMessage struct {
				Actions []types.AgentMessagePayloadActions
			}

			err = json.Unmarshal(msg.GetPayload(), &actionsMessage)
			if err != nil {

				return bettererrors.
					New("Failed to unmarshal JSON agent actions").
					SetContext("agent", agentproxy.String()).
					SetContext("payload", string(msg.GetPayload()))
			}

			mutationbatch := types.AgentMutationBatch{
				AgentProxyUUID: agentproxy.GetProxyUUID(),
				AgentEntityId:  agentproxy.GetEntityId(),
				Mutations:      actionsMessage.Actions,
			}

			server.PushMutationBatch(mutationbatch)

			break
		}
	default:
		{
			berror := bettererrors.
				New("Unknown message type").
				SetContext("method", msg.GetMethod())

			assert.AssertBE(false, berror)
		}
	}

	return nil
}

/* </implementing types.CommunicatorDispatcherInterface> */

package module

import (
	"encoding/json"
	"time"

	"github.com/edwardbrowncross/amazon-connect-simulator/event"
	"github.com/edwardbrowncross/amazon-connect-simulator/flow"
)

// CallConnector describes what a module needs to interact with the ongoing call.
type CallConnector interface {
	Send(s string, ssml bool)
	Receive(count int, timeout time.Duration, encrypt bool) *string
	Emit(event event.Event)
	GetExternal(key string) interface{}
	SetExternal(key string, value interface{})
	ClearExternal()
	GetContactData(key string) interface{}
	SetContactData(key string, value interface{})
	GetSystem(key string) interface{}
	SetSystem(key string, value interface{})
	InvokeLambda(named string, inParams json.RawMessage) (outJSON string, outErr error, err error)
	GetFlowStart(flowName string) *flow.ModuleID
}

// Runner takes a call context and returns the ID of the next block to run, or nil if the call is over.
type Runner interface {
	Run(CallConnector) (*flow.ModuleID, error)
}

// MakeRunner takes the data of a module (block) and wraps it in a type that provides the functionality of that block.
func MakeRunner(m flow.Module) Runner {
	switch m.Type {
	case flow.ModuleStoreUserInput:
		return storeUserInput(m)
	case flow.ModuleCheckAttribute:
		return checkAttribute(m)
	case flow.ModuleTransfer:
		return transfer(m)
	case flow.ModulePlayPrompt:
		return playPrompt(m)
	case flow.ModuleDisconnect:
		return disconnect(m)
	case flow.ModuleSetQueue:
		return setQueue(m)
	case flow.ModuleGetUserInput:
		return getUserInput(m)
	case flow.ModuleSetAttributes:
		return setAttributes(m)
	case flow.ModuleInvokeExternalResource:
		return invokeExternalResource(m)
	case flow.ModuleCheckHoursOfOperation:
		return checkHoursOfOperation(m)
	default:
		return passthrough(m)
	}
}

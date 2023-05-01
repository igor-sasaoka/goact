package goact

import (
	"errors"
	"fmt"
)

//Class that allow action handler functions to write data into message
type MessageWriter interface {
    Write([]byte) (int, error)
    Flush() error
}

type Action func (MessageWriter, string) error 
type ActionHandler struct {
    actions map[string]Action
}

func (h *ActionHandler) executeAction(actionId string, commWriter MessageWriter, messageBody string) error {
    action, err := h.callAction(actionId)
    if err != nil {
        return err
    }
    actionErr := action(commWriter, messageBody)
    if actionErr != nil {
        return actionErr
    }

    return nil
}

func (h *ActionHandler) Register(actionId string, action Action) {
    actions := h.getActions()
    actions[actionId] = action
}

func (h *ActionHandler) callAction(actionId string) (Action, error) {
    actions := h.getActions()
    if ac, ok := actions[actionId]; ok {
       return ac, nil
    }
    
    return nil, errors.New(fmt.Sprintf("goact: action [%s] not found", actionId)) 
}

func (h *ActionHandler) getActions() map[string]Action {
    if h.actions == nil {
        h.actions = make(map[string]Action)
    }
    return h.actions
}

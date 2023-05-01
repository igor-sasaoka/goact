package goact

import (
    "encoding/json"
    "errors"
    "fmt"
)

type Message struct {
    Id string `json:"id"`
    Action string `json:"action"`
    Body string `json:"body"`
}

func DecodeMessage(b []byte) (*Message, error) {
    m := &Message{}
    err := json.Unmarshal(b, m)
    if err != nil {
        return nil, err
    } 

    validationErr := validateMessage(m)
    if validationErr != nil {
        return nil, err
    } 

    return m, nil
}

func validateMessage(m *Message) error {
    var missingFields string
    if m.Id == "" {
        missingFields = missingFields + "[id] " 
    }
    if m.Action == "" {
        missingFields = missingFields + "[action] "
    }

    if missingFields != "" {
        return errors.New(fmt.Sprintf("goact: the following data is missing: %s", missingFields))
    }
    return nil
}

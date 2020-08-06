package xfsm

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
)

type Callback func(*tgbotapi.Message) error

type State struct {
	Name     string
	Callback func(*tgbotapi.Message) error
}

type Transition struct {
	Name        string
	Destination string
	Source      string
}

func NewFSM(initialState string,
	transitionList []Transition,
	stateList []State) *tgFSM {
	transitions := map[transitionHead]string{}
	for _, transition := range transitionList {
		transitions[transitionHead{transition.Name, transition.Source}] = transition.Destination
	}
	states := map[string]Callback{}
	for _, state := range stateList {
		states[state.Name] = state.Callback
	}
	return &tgFSM{
		current:     initialState,
		transitions: transitions,
		states:      states,
	}
}

type transitionHead struct {
	event string
	state string
}

type tgFSM struct {
	mx      sync.RWMutex
	current string

	transitions map[transitionHead]string
	states      map[string]Callback
}

func (s *tgFSM) Current() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.current
}

func (s *tgFSM) Event(event string, msg *tgbotapi.Message) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	nextSate, found := s.transitions[transitionHead{event, s.current}]
	if !found {
		if stateCallback, found := s.states[s.current]; found {
			s.current = nextSate
			return stateCallback(msg)
		}
	}

	if stateCallback, found := s.states[nextSate]; found {
		s.current = nextSate
		return stateCallback(msg)
	}

	return nil
}

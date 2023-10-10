package models

import "testing"

func TestMatchIsNotAtEndColossusStateAbandoned(t *testing.T) {
	cm := &ColossusMatch{
		Status: MatchColossusStatusAbandoned,
	}
	isNotEndState := cm.IsNotAtEndColossusState()
	if isNotEndState {
		t.Error("Expected abandoned colossus state to be the end state")
	}
}

func TestMatchIsNotAtEndColossusStateOfficial(t *testing.T) {
	cm := ColossusMatch{
		Status: MatchColossusStatusOfficial,
	}
	isNotEndState := cm.IsNotAtEndColossusState()
	if isNotEndState {
		t.Error("Expected official colossus state to be the end state")
	}
}

func TestMatchIsNotAtEndColossusStateInPlay(t *testing.T) {
	cm := ColossusMatch{
		Status: MatchColossusStatusInPlay,
	}
	isNotEndState := cm.IsNotAtEndColossusState()
	if !isNotEndState {
		t.Error("Expected in play colossus state to not be the end state")
	}
}

func TestMatchIsNotAtEndColossusStateCompleted(t *testing.T) {
	cm := ColossusMatch{
		Status: MatchColossusStatusCompleted,
	}
	isNotEndState := cm.IsNotAtEndColossusState()
	if !isNotEndState {
		t.Error("Expected completed colossus state to not be the end state")
	}
}

func TestMatchIsNotAtEndColossusStateNotStarted(t *testing.T) {
	cm := ColossusMatch{
		Status: MatchColossusStatusNotStarted,
	}
	isNotEndState := cm.IsNotAtEndColossusState()
	if !isNotEndState {
		t.Error("Expected not started colossus state to not be the end state")
	}
}

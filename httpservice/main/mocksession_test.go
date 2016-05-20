package main

import "log"

type mockSession struct {
	validateSessionValidResponse        bool
	validateSessionEmailAddressResponse string
	validateSessionWasCalled            bool
	startSessionWasCalled               bool
}

func (ms *mockSession) ValidateSession() (isValid bool, emailAddress string) {
	log.Print("The session was validated.")
	ms.validateSessionWasCalled = true
	return ms.validateSessionValidResponse, ms.validateSessionEmailAddressResponse
}

func (ms *mockSession) StartSession(emailAddress string) {
	log.Print("The session was started.")
	ms.startSessionWasCalled = true
}

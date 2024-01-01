package handlers

type testHandlers struct {
	auth    *AuthHandlers
	user    *UserHandlers
	contact *ContactHandlers
}

func newTestHandlers() *testHandlers {
	return &testHandlers{
		auth:    NewTestAuthHandlers(),
		user:    NewTestUserHandlers(),
		contact: NewTestContactHandlers(),
	}
}

var h = newTestHandlers()

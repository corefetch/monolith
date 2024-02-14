package route

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"learnt.io/core"
	"learnt.io/core/rest"
	"learnt.io/modules/accounts/service"
)

func TestAuth(t *testing.T) {

	core.TestSetup()

	account, err := service.Register(service.CreateAcountData{
		Role:          core.RoleStudent,
		Login:         "jhon@doe.com",
		Password:      "Asdasdasd1",
		PasswordMatch: "Asdasdasd1",
	})

	if err != nil {
		t.Error(err)
		return
	}

	defer account.Drop()

	var data = bytes.NewBufferString(`
		{
			"login": "jhon@doe.com",
			"password": "Asdasdasd1"
		}
	`)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", data)

	Auth(rest.NewContext(w, r))

	if w.Result().StatusCode != http.StatusOK {
		t.Error("auth expected:", w.Result().StatusCode)
	}
}

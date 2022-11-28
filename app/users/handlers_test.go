package users

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type serviceMock struct {
}

func (s *serviceMock) CreateUser(ctx context.Context, FirstName, LastName, Email, Password string) (*UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func newServiceMock() Service {
	return &serviceMock{}
}

func TestPOSTCreateUser(t *testing.T) {
	service := newServiceMock()
	router := mux.NewRouter()
	RegisterRoutes(router, service)
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/hello")
	if err != nil {
		log.Fatalf("after get  : %#v", err.Error())
	}
	rawJson, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf("reading body : %#v", err)
	}

	var sr MockUser
	err = json.Unmarshal(rawJson, &sr)
	if err != nil {
		t.Fatalf("error decoding json : %#v", err)
	}

	t.Logf("unmarshalled response : %#v", sr)

	if sr.FirstName != "test" {
		t.Fatalf("failed - expecting test, got %s", sr.FirstName)
	}
}

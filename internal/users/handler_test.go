package users

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"bitbucket.org/faceit/internal/users/usersapi"

	"bitbucket.org/faceit/app"

	"github.com/steinfletcher/apitest"
)

func TestHandler_Create(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}

	apitest.New(). // configuration
			HandlerFunc(handler).
			Post("/users").
			Expect(t).
			Status(http.StatusCreated).
			End()
}

func TestHandler_Delete(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	apitest.New(). // configuration
			HandlerFunc(handler).
			Delete("/users/test-user-id").
			Expect(t).
			Status(http.StatusNoContent).
			End()
}

func TestHandler_List(t *testing.T) {
	now := time.Now()
	handler := func(w http.ResponseWriter, r *http.Request) {
		response := usersapi.UserListResponse{}
		u := usersapi.User{
			ID:        "d2a7924e-765f-4949-bc4c-219c956d0f8b",
			FirstName: "Alice",
			LastName:  "Bob",
			Nickname:  "AB123",
			Password:  "supersecurepassword",
			Email:     "alice@bob.com",
			Country:   "UK",
			CreatedAt: now,
			UpdatedAt: now,
		}
		response.Users = append(response.Users, u)

		app.JSON(w, http.StatusOK, response)
	}

	response := usersapi.UserListResponse{}
	u := usersapi.User{
		ID:        "d2a7924e-765f-4949-bc4c-219c956d0f8b",
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "AB123",
		Password:  "supersecurepassword",
		Email:     "alice@bob.com",
		Country:   "UK",
		CreatedAt: now,
		UpdatedAt: now,
	}
	response.Users = append(response.Users, u)

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal response: %s", err)
	}
	apitest.New(). // configuration
			HandlerFunc(handler).
			Get("/users").
			Expect(t).
			Body(string(data)).
			Status(http.StatusOK).
			End()
}

func TestHandler_Update(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	apitest.New(). // configuration
			HandlerFunc(handler).
			Patch("/users/test-user-id").
			Expect(t).
			Status(http.StatusNoContent).
			End()
}

func Test_makeHash(t *testing.T) {
	got := makeHash("123456")
	want := "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

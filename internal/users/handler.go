package users

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"strings"

	"bitbucket.org/faceit/app"
	"bitbucket.org/faceit/internal/events"
	"bitbucket.org/faceit/internal/users/usersapi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type userService interface {
	Get(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID string) error
	List(ctx context.Context, filter Filter) ([]User, error)
}

type eventService interface {
	Publish(ctx context.Context, e events.Event) error
}

type Handler struct {
	Users  userService
	Events eventService
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`
}

func (c CreateUserRequest) Validate(r *http.Request) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required),
		validation.Field(&c.LastName, validation.Required),
		validation.Field(&c.Nickname, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Country, validation.Required, is.CountryCode2),
	)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateUserRequest
	if !app.BindAndValidate(w, r, &req) {
		log.Fatalln("Failed to bind and validate request CreateUserRequest")
		return
	}

	user, err := h.Users.Get(ctx, req.Email)
	if err == nil {
		if user.Email == req.Email {
			app.ErrorResponse(w, ErrEmailAlreadyExists, http.StatusBadRequest)
			return
		}
	}

	var userID = uuid.New().String()
	params := &User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  makeHash(req.Password),
		Email:     req.Email,
		Country:   strings.ToUpper(req.Country),
	}

	if err := h.Users.Create(ctx, params); err != nil {
		app.InternalError(w, r, err)
		return
	}

	event := &events.UserCreated{
		ID: userID,
	}

	if err := h.Events.Publish(ctx, event); err != nil {
		log.Fatalln("Failed to publish event UserCreated")
		app.ReportError(r, err)
	}

	log.Printf("User Successfuly Created: %s", userID)
	w.WriteHeader(http.StatusCreated)
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`
}

func (u UpdateUserRequest) Validate(r *http.Request) error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Nickname, validation.Required),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Country, validation.Required, is.CountryCode2),
	)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	var req UpdateUserRequest
	if !app.BindAndValidate(w, r, &req) {
		return
	}

	var userID = vars["id"]
	params := &User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  makeHash(req.Password),
		Email:     req.Email,
		Country:   strings.ToUpper(req.Country),
	}

	if err := h.Users.Update(ctx, params); err != nil {
		app.InternalError(w, r, err)
		return
	}

	event := &events.UserUpdated{
		ID: userID,
	}
	if err := h.Events.Publish(ctx, event); err != nil {
		log.Fatalln("Failed to publish event UserUpdated")
		app.ReportError(r, err)
	}

	log.Printf("User Successfuly Updated: %s", userID)

	w.WriteHeader(http.StatusNoContent)
}

type DeleteUserRequest struct {
	ID string `gorm:"column:id"`
}

func (d DeleteUserRequest) Validate(r *http.Request) error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ID),
	)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	var req DeleteUserRequest
	if !app.BindAndValidate(w, r, &req) {
		return
	}

	var userID = vars["id"]

	if err := h.Users.Delete(ctx, userID); err != nil {
		app.InternalError(w, r, err)
		return
	}

	event := &events.UserDeleted{
		ID: userID,
	}
	if err := h.Events.Publish(ctx, event); err != nil {
		log.Fatalln("Failed to publish event UserDeleted")
		app.ReportError(r, err)
	}

	log.Printf("User Successfuly Deleted: %s", userID)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()

	var filter = Filter{
		Page:     q.Get("page"),
		PageSize: q.Get("page_size"),
		Country:  strings.ToUpper(q.Get("country")),
	}

	users, err := h.Users.List(ctx, filter)
	if err != nil {
		app.InternalError(w, r, err)
		return
	}

	if len(users) == 0 {
		log.Printf("No users found")
	}

	app.JSON(w, http.StatusOK, userListResponse(users))
}

func userListResponse(users []User) usersapi.UserListResponse {
	response := usersapi.UserListResponse{}
	for _, user := range users {
		u := usersapi.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Country:   user.Country,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		response.Users = append(response.Users, u)
	}

	if len(response.Users) == 0 {
		response.Users = []usersapi.User{}
	}
	return response
}

func makeHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	fmt.Printf("%x", h.Sum(nil))
	return fmt.Sprintf("%x", h.Sum(nil))
}

package http

import (
	"6/internal/models"
	"6/internal/store"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"6/internal/message_broker"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"strconv"
)

type UserResource struct {
	store store.Store
	broker message_broker.MessageBroker
	cache *lru.TwoQueueCache
}

func NewUserResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *UserResource {
	return &UserResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}

func (ac *UserResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", ac.CreateUser)
	r.Get("/", ac.AllUsers)
	r.Get("/{id}", ac.ByID)
	r.Put("/", ac.UpdateUser)
	r.Delete("/{id}", ac.DeleteUser)

	return r
}

func (cr *UserResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if err := user.BeforeCreating(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	if err := cr.store.Users().Create(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if err = ur.broker.Cache().Purge(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}


func (cr *UserResource) AllUsers(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := new(models.NameFilter)

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		usersFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, usersFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	users, err := cr.store.Users().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" && len(users) > 0 {
		cr.cache.Add(searchQuery, users)
	}
	render.JSON(w, r, users)
}

func (cr *UserResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	userFromCache, ok := cr.cache.Get(id)
	if ok {
		render.JSON(w, r, userFromCache)
		return
	}

	user, err := cr.store.Users().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.cache.Add(id, user)
	render.JSON(w, r, user)
}

func (cr *UserResource) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	// err := validation.ValidateStruct(
	// 	user,
	// 	validation.Field(&user.ID, validation.Required),
	// 	validation.Field(&user.Nick, validation.Required),
	// )
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// 	fmt.Fprintf(w, "Unknown err: %v", err)
	// 	return
	// }

	if err := cr.store.Users().Update(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if err := ur.broker.Cache().Remove(user.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}

func (cr *UserResource) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Users().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if err = ur.broker.Cache().Remove(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}

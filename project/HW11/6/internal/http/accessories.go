package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	lru "github.com/hashicorp/golang-lru"
	"6/internal/models"
	"6/internal/store"
	"6/internal/message_broker"
	"net/http"
	"strconv"
)

type AccessoryResource struct {
	store store.Store
	broker message_broker.MessageBroker
	cache *lru.TwoQueueCache
}

func NewAccessoryResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *AccessoryResource {
	return &AccessoryResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}

func (ac *AccessoryResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", ac.CreateAccessory)
	r.Get("/", ac.AllAccessories)
	r.Get("/{id}", ac.ByID)
	r.Put("/", ac.UpdateAccessory)
	r.Delete("/{id}", ac.DeleteAccessory)

	return r
}

func (cr *AccessoryResource) CreateAccessory(w http.ResponseWriter, r *http.Request) {
	accessory := new(models.Accessory)
	if err := json.NewDecoder(r.Body).Decode(accessory); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Accessories().Create(r.Context(), accessory); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if err := cr.broker.Cache().Purge(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)

func (cr *AccessoryResource) AllAccessories(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.AccessoriesFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		accessoriesFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, accessoriesFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	accessories, err := cr.store.Accessories().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" && len(accessories) > 0 {
		cr.cache.Add(searchQuery, accessories)
	}
	render.JSON(w, r, accessories)
}

func (cr *AccessoryResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	accessory, err := cr.store.Accessories().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	render.JSON(w, r, accessory)
}

func (cr *AccessoryResource) UpdateAccessory(w http.ResponseWriter, r *http.Request) {
	accessory := new(models.Accessory)
	if err := json.NewDecoder(r.Body).Decode(accessory); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		accessory,
		validation.Field(&accessory.ID, validation.Required),
		validation.Field(&accessory.Name, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Accessories().Update(r.Context(), accessory); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	if err = cr.broker.Cache().Remove(accessory.Id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}

func (cr *AccessoryResource) DeleteAccessory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Accessories().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if err = cr.broker.Cache().Remove(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}
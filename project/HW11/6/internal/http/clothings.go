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

type ClothingResource struct {
	store store.Store
	broker message_broker.MessageBroker
	cache *lru.TwoQueueCache
}

func NewClothingResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *ClothingResource {
	return &CLothingResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}

func (ac *ClothingResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", ac.CreateClothing)
	r.Get("/", ac.AllClothings)
	r.Get("/{id}", ac.ByID)
	r.Put("/", ac.UpdateClothing)
	r.Delete("/{id}", ac.DeleteClothing)

	return r
}

func (cr *ClothingResource) CreateClothing(w http.ResponseWriter, r *http.Request) {
	clothing := new(models.Clothing)
	if err := json.NewDecoder(r.Body).Decode(clothing); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Clothings().Create(r.Context(), clothing); err != nil {
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
}

func (cr *ClothingResource) AllClothings(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ClothingsFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		clothingsFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, clothingsFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	clothings, err := cr.store.CLothings().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" && len(clothings) > 0 {
		cr.cache.Add(searchQuery, clothings)
	}
	render.JSON(w, r, clothings)
}

func (cr *CLothigResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}



	clothing, err := cr.store.Clothings().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, clothing)
}

func (cr *ClothingResource) UpdateClothing(w http.ResponseWriter, r *http.Request) {
	clothing := new(models.Clothing)
	if err := json.NewDecoder(r.Body).Decode(clothing); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		clothing,
		validation.Field(&clothing.ID, validation.Required),
		validation.Field(&clothing.Name, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Clothing().Update(r.Context(), clothing); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	if err = cr.broker.Cache().Remove(clothing.Id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}

func (cr *ClothingResource) DeleteClothing(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Clothing().Delete(r.Context(), id); err != nil {
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
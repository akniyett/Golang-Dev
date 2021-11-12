package http

import (
	"context"
	"encoding/json"
	"fmt"
	"6/internal/models"
	"6/internal/store"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,

		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Post("/clothings", func(w http.ResponseWriter, r *http.Request) {
		clothing := new(models.Clothing)
		if err := json.NewDecoder(r.Body).Decode(clothing); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		if err := s.store.Clothings().Create(r.Context(), clothing); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.Post("/accessories", func(w http.ResponseWriter, r *http.Request) {
		accessory := new(models.Accessory)
		if err := json.NewDecoder(r.Body).Decode(accessory); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Accessories().Create(r.Context(), accessory); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Users().Create(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})


	r.Get("/clothings", func(w http.ResponseWriter, r *http.Request) {
		clothings, err := s.store.Clothings().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, clothings)
	})

	r.Get("/accessories", func(w http.ResponseWriter, r *http.Request) {
		accessories, err := s.store.Accessories().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, accessories)
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.Users().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, users)
	})

	r.Get("/clothings/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		clothing, err := s.store.Clothings().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, clothing)
	})

	r.Get("/accessories/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		accessory, err := s.store.Accessories().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, accessory)
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		user, err := s.store.Users().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, user)
	})

	r.Put("/clothings", func(w http.ResponseWriter, r *http.Request) {
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

		if err := s.store.Clothings().Update(r.Context(), clothing); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Put("/accessories", func(w http.ResponseWriter, r *http.Request) {
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

		if err := s.store.Accessories().Update(r.Context(), accessory); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Put("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		err := validation.ValidateStruct(
			user,
			validation.Field(&user.ID, validation.Required),
			validation.Field(&user.Nick, validation.Required),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Users().Update(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Delete("/clothings/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Clothings().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Delete("/accessories/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Accessories().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Users().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // блокируемся, пока контекст приложения не отменен

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}

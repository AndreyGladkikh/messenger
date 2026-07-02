package http_server

import (
	"encoding/json"
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/application/use_cases/send_message"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	commandBus *command.Bus
}

func NewServer(
	commandBus *command.Bus,
) *Server {
	return &Server{
		commandBus: commandBus,
	}
}

func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	// r.Use(middleware.Heartbeat("/ping"))

	// r.Get("/", func (w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })

	// r.Get("/test", s.test)

	s.registerApi(r)

	err := http.ListenAndServe(":3000", r)

	return err
}

func (s *Server) registerApi(mux *chi.Mux) {
	mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.Route("/messages", func(r chi.Router) {
		r.Post("/", s.sendMessage)
	})
}

func (s *Server) success(w http.ResponseWriter, _ *http.Request, v any, httpStatus int) {
	response := map[string]any{
		"status": "success",
		"data":   v,
	}
	// encoded, e := encode(response)
	// if e != nil {
	// 	http.Error(w, e.Error(), http.StatusInternalServerError)
	// }

	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w).Encode(response)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	// w.Write(encoded)
}

func (s *Server) error(w http.ResponseWriter, _ *http.Request, err error) {
	response := map[string]any{
		"status": "error",
		"error":  err,
	}

	// encoded, e := encode(response)
	// if e != nil {
	// 	http.Error(w, e.Error(), http.StatusInternalServerError)
	// }

	w.WriteHeader(getErrorStatus(err))
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w).Encode(response)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	// w.Write(encoded)
}

func (s *Server) sendMessage(w http.ResponseWriter, r *http.Request) {
	var request Message
	json.NewDecoder(r.Body).Decode(&request)

	command := &send_message.Command{
		ChatID:           request.ChatID,
		MessageBody:      request.Body,
		ReplyToMessageID: request.ReplyToMessageID,
		Attachments:      request.Attachments,
	}
	response, err := s.commandBus.Dispatch(r.Context(), command)

	if err != nil {
		s.error(w, r, err)
	}

	s.success(w, r, response, http.StatusOK)
}

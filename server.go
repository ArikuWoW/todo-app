package todoapp

import (
	"context"
	"net/http"
	"time"
)

// Обертка над стандартным сервером go, для запуска/остановки сервера
type Server struct {
	httpServer *http.Server
}

// метод для запуска сервера
func (s *Server) Run(port string, handler http.Handler) error {
	// Создаем новый http сервер и сохраняем в поле структуры
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

// Метод при выходе из приложения
func (s *Server) Shutdown(ctx context.Context) error {
	// При таком вызове сервер перестает принимать запросы,
	// Но дожидается завершения обработки текущих
	return s.httpServer.Shutdown(ctx)
}

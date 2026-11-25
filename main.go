package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Item adalah data model yang disimpan di Firestore.
type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// Device adalah data model collection device.
type Device struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Boottime adalah data model collection boottime.
type Boottime struct {
	ID        string    `json:"id"`
	BootTime  string    `json:"boot_time"`
	CreatedAt time.Time `json:"created_at"`
}

// Server adalah HTTP server yang menggunakan Firestore client.
type Server struct {
	firestore *firestore.Client
	collection string
}

func NewServer(ctx context.Context, projectID, collection string) (*Server, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}
	return &Server{
		firestore:  client,
		collection: collection,
	}, nil
}

func (s *Server) Close() error {
	return s.firestore.Close()
}

func (s *Server) GetAllItems() ([]Item, error) {
	var items []Item
	iter := s.firestore.Collection(s.collection).Documents(context.Background())
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var item Item
		if err := doc.DataTo(&item); err != nil {
			continue
		}
		if item.ID == "" {
			item.ID = doc.Ref.ID
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Server) GetItemByID(id string) (Item, error) {
	doc, err := s.firestore.Collection(s.collection).Doc(id).Get(context.Background())
	if err != nil {
		return Item{}, err
	}
	var item Item
		if err := doc.DataTo(&item); err != nil {
			return Item{}, err
		}
	if item.ID == "" {
		item.ID = doc.Ref.ID
	}
	return item, nil
}

func (s *Server) CreateItem(name string) (Item, error) {
	item := Item{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:      name,
		Done:      false,
		CreatedAt: time.Now(),
	}
	_, err := s.firestore.Collection(s.collection).Doc(item.ID).Set(context.Background(), item)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func (s *Server) DeleteItem(id string) error {
	_, err := s.firestore.Collection(s.collection).Doc(id).Delete(context.Background())
	return err
}

func (s *Server) GetAllDevices() ([]Device, error) {
	var devices []Device
	iter := s.firestore.Collection("device").Documents(context.Background())
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var device Device
		if err := doc.DataTo(&device); err != nil {
			continue
		}
		if device.ID == "" {
			device.ID = doc.Ref.ID
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (s *Server) GetDeviceByID(id string) (Device, error) {
	doc, err := s.firestore.Collection("device").Doc(id).Get(context.Background())
	if err != nil {
		return Device{}, err
	}
	var device Device
	if err := doc.DataTo(&device); err != nil {
		return Device{}, err
	}
	if device.ID == "" {
		device.ID = doc.Ref.ID
	}
	return device, nil
}

func (s *Server) CreateDevice(name string) (Device, error) {
	device := Device{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:      name,
		CreatedAt: time.Now(),
	}
	_, err := s.firestore.Collection("device").Doc(device.ID).Set(context.Background(), device)
	if err != nil {
		return Device{}, err
	}
	return device, nil
}

func (s *Server) DeleteDevice(id string) error {
	_, err := s.firestore.Collection("device").Doc(id).Delete(context.Background())
	return err
}

func (s *Server) GetAllBoottimes() ([]Boottime, error) {
	var boottimes []Boottime
	iter := s.firestore.Collection("boottime").Documents(context.Background())
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var boottime Boottime
		if err := doc.DataTo(&boottime); err != nil {
			continue
		}
		if boottime.ID == "" {
			boottime.ID = doc.Ref.ID
		}
		boottimes = append(boottimes, boottime)
	}
	return boottimes, nil
}

func (s *Server) GetBoottimeByID(id string) (Boottime, error) {
	doc, err := s.firestore.Collection("boottime").Doc(id).Get(context.Background())
	if err != nil {
		return Boottime{}, err
	}
	var boottime Boottime
	if err := doc.DataTo(&boottime); err != nil {
		return Boottime{}, err
	}
	if boottime.ID == "" {
		boottime.ID = doc.Ref.ID
	}
	return boottime, nil
}

func (s *Server) CreateBoottime(bootTime string) (Boottime, error) {
	boottime := Boottime{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		BootTime:  bootTime,
		CreatedAt: time.Now(),
	}
	_, err := s.firestore.Collection("boottime").Doc(boottime.ID).Set(context.Background(), boottime)
	if err != nil {
		return Boottime{}, err
	}
	return boottime, nil
}

func (s *Server) DeleteBoottime(id string) error {
	_, err := s.firestore.Collection("boottime").Doc(id).Delete(context.Background())
	return err
}

func main() {
	ctx := context.Background()

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		projectID = "arduino-ee66f"
	}

	collection := os.Getenv("FIRESTORE_COLLECTION")
	if collection == "" {
		collection = "items"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server, err := NewServer(ctx, projectID, collection)
	if err != nil {
		log.Fatalf("failed to create firestore server: %v", err)
	}
	defer server.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		items, err := server.GetAllItems()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch items")
			return
		}
		writeJSON(w, http.StatusOK, items)
	})

	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		item, err := server.GetItemByID(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "item not found")
			return
		}
		writeJSON(w, http.StatusOK, item)
	})

	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}
		item, err := server.CreateItem(req.Name)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create item")
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})

	mux.HandleFunc("DELETE /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := server.DeleteItem(id); err != nil {
			writeError(w, http.StatusNotFound, "item not found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	})

	mux.HandleFunc("GET /devices", func(w http.ResponseWriter, r *http.Request) {
		devices, err := server.GetAllDevices()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch devices")
			return
		}
		writeJSON(w, http.StatusOK, devices)
	})

	mux.HandleFunc("GET /devices/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		device, err := server.GetDeviceByID(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "device not found")
			return
		}
		writeJSON(w, http.StatusOK, device)
	})

	mux.HandleFunc("POST /devices", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}
		device, err := server.CreateDevice(req.Name)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create device")
			return
		}
		writeJSON(w, http.StatusCreated, device)
	})

	mux.HandleFunc("DELETE /devices/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := server.DeleteDevice(id); err != nil {
			writeError(w, http.StatusNotFound, "device not found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	})

	mux.HandleFunc("GET /boottimes", func(w http.ResponseWriter, r *http.Request) {
		boottimes, err := server.GetAllBoottimes()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch boottimes")
			return
		}
		writeJSON(w, http.StatusOK, boottimes)
	})

	mux.HandleFunc("GET /boottimes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		boottime, err := server.GetBoottimeByID(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "boottime not found")
			return
		}
		writeJSON(w, http.StatusOK, boottime)
	})

	mux.HandleFunc("POST /boottimes", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			BootTime string `json:"boot_time"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.BootTime == "" {
			writeError(w, http.StatusBadRequest, "boot_time is required")
			return
		}
		boottime, err := server.CreateBoottime(req.BootTime)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create boottime")
			return
		}
		writeJSON(w, http.StatusCreated, boottime)
	})

	mux.HandleFunc("DELETE /boottimes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := server.DeleteBoottime(id); err != nil {
			writeError(w, http.StatusNotFound, "boottime not found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	})

	mux.HandleFunc("GET /collections", func(w http.ResponseWriter, r *http.Request) {
		names, err := listCollections(context.Background(), server.firestore)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list collections")
			return
		}
		writeJSON(w, http.StatusOK, names)
	})

	handler := loggingMiddleware(mux)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("starting backend server on :%s", port)
	log.Printf("firestore project: %s, collection: %s", projectID, collection)
	log.Printf("endpoints:")
	log.Printf("  GET  /health")
	log.Printf("  GET  /items")
	log.Printf("  GET  /items/{id}")
	log.Printf("  POST /items")
	log.Printf("  DELETE /items/{id}")
	log.Printf("  GET  /devices")
	log.Printf("  GET  /devices/{id}")
	log.Printf("  POST /devices")
	log.Printf("  DELETE /devices/{id}")
	log.Printf("  GET  /boottimes")
	log.Printf("  GET  /boottimes/{id}")
	log.Printf("  POST /boottimes")
	log.Printf("  DELETE /boottimes/{id}")
	log.Printf("  GET  /collections")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, ww.statusCode, duration.Round(time.Millisecond))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error encoding json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func listCollections(ctx context.Context, client *firestore.Client) ([]string, error) {
	iter := client.Collections(ctx)
	var names []string

	for {
		colRef, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		names = append(names, colRef.ID)
	}
	return names, nil
}

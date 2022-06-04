package server

import (
	_ "advanced-app/internal/docs"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"net/http/httputil"
)

var r *mux.Router

type user struct {
	Name  string `json:"name" validate:"required "`
	Email string `json:"email"`
}

func Start(addr string) {
	r = mux.NewRouter()
	endpoints(r)

	go func() {
		log.Error(http.ListenAndServe(addr, r))
	}()
}

func endpoints(r *mux.Router) {
	r.HandleFunc("/healthz", getHealthz).Methods(http.MethodGet)
	r.HandleFunc("/user", createUser).Methods(http.MethodPost)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
}

// @Summary Get healthz
// @Tags healthz
// @Description Get healthz
// @Accept  json
// @Success 200 {integer} integer 1
// @Router /healthz [get]
func getHealthz(w http.ResponseWriter, r *http.Request) {
	printer(w, r)

	w.WriteHeader(http.StatusOK)
}

// @Summary Create user
// @Tags user
// @Description Create user
// @Accept  json
// @Produce  json
// @Param user body user true "User"
// @Success  200  {object}  user
// @Failure  400  {object}  string
// @Router /user [post]
func createUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	printer(w, r)
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Error("Error:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if u.Name == "" || u.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	log.Info("User:", u)

	respondWithJSON(w, http.StatusOK, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Error("Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Error("Error:", err)
		return
	}
}

func printer(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Error("Error:", err)
	}
	log.Info("Request:", string(dump))
}

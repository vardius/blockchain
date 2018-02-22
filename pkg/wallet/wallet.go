package wallet

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vardius/blockchain/pkg/client"
	"github.com/vardius/gorouter"
)

// Wallet interface
type Wallet interface {
	Run(host string, port int) error
}

// TODO(vardius): add uuid to wallet
type httpWallet struct {
	client client.BlockchainClient
}

// NewHTTP creates new http wallet instance
func NewHTTP(c client.BlockchainClient) Wallet {
	return &httpWallet{c}
}

// Run runs the http server
func (wl *httpWallet) Run(host string, port int) error {
	router := wl.newRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", host, port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println(fmt.Sprintf("Wallet listening on %s:%d", host, port))

	return s.ListenAndServe()
}

// newRouter creates http handler
func (wl *httpWallet) newRouter() http.Handler {
	router := gorouter.New()
	router.GET("/", http.HandlerFunc(wl.getBlockchain))
	router.POST("/", http.HandlerFunc(wl.writeBlock))

	return router
}

func (wl *httpWallet) getBlockchain(w http.ResponseWriter, r *http.Request) {
	bc, err := wl.client.GetBlockchain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	bytes, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	io.WriteString(w, string(bytes))
}

func (wl *httpWallet) writeBlock(w http.ResponseWriter, r *http.Request) {
	var message struct {
		Data string `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)

		return
	}
	defer r.Body.Close()

	err := wl.client.AddBlock(message.Data)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, err)

		return
	}

	respondWithJSON(w, r, http.StatusCreated, "Ok")
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))

		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

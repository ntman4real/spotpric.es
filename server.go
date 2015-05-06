package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/unrolled/render"
)

var (
	rmu           = sync.RWMutex{}
	currentPrices priceHistoryResponse
)

func updatePrices() {
	p, _ := getPrices()

	rmu.Lock()
	defer rmu.Unlock()
	currentPrices = p
}

func sendOnHub(h *hub) {
	rmu.RLock()
	defer rmu.RUnlock()

	if b, err := json.Marshal(currentPrices.Items); err == nil {
		h.broadcast <- b
	}
}

func main() {
	h := newHub()
	go h.run()

	go func() {
		updatePrices()
		sendOnHub(h)

		ticker := time.Tick(30 * time.Second)
		for _ = range ticker {
			updatePrices()
			sendOnHub(h)
		}
	}()

	r := render.New(render.Options{
		IsDevelopment: true,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		rmu.RLock()
		defer rmu.RUnlock()

		r.HTML(w, http.StatusOK, "index", currentPrices.Items)
	})

	mux.Handle("/ws", wsHandler{h: h})

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), gzip.Gzip(gzip.DefaultCompression), negroni.NewStatic(http.Dir("public")))
	n.UseHandler(mux)
	n.Run(":3001")
}

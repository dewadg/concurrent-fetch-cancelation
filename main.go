package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dewadg/concurrent-fetch-cancelation/repositories"
)

func main() {
	photoRepository := repositories.NewPhotoRepository()

	server := http.NewServeMux()
	server.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		timeout := 10 * time.Second
		if timeoutQuery := request.URL.Query().Get("timeout"); timeoutQuery != "" {
			value, err := time.ParseDuration(timeoutQuery + "ms")
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}

			timeout = value
		}

		photos, totalCount, successCount, errorCount, cancelledCount, err := photoRepository.Get(request.Context(), timeout)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"total":     totalCount,
			"success":   successCount,
			"error":     errorCount,
			"cancelled": cancelledCount,
			"data":      photos,
		}
		response, err := json.Marshal(data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprint(writer, string(response))
	})

	http.ListenAndServe(":8000", server)
}

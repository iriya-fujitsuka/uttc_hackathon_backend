package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseMessage struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// クエリパラメータ "name" を取得
	name := r.URL.Query().Get("name")
	if name == "" {
		// nameパラメータが設定されていない場合はBadRequestを返す
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BadRequest(status code = 400)"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bytes, err := json.Marshal(responseMessage{
		Message: "Hello, " + name + "-san" + "!",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func main() {
	http.HandleFunc("/hello", handler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

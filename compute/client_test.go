package compute

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientAccessDenied(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "Invalid credentials.", http.StatusUnauthorized)
	}))
	defer testServer.Close()

	request, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		log.Fatal(err)

		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)

		return
	}
	defer response.Body.Close()

	fmt.Printf("%s\n", response.Status)
}

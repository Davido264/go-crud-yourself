package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func normalizeRequests(req *http.Request, reqp *http.Request, content []byte) {
	log.Println("Translating :3")
	defer func() {
		log.Printf("%v --> %v", reqp.URL.String(), req.URL.String())
	}()

	previousUrl := reqp.URL
	fragments := strings.Split(previousUrl.Path, "/")

	if req.Method == http.MethodDelete {
		id, err := strconv.Atoi(fragments[len(fragments)-1])
		if err != nil {
			// Flan
			key := getIdField(fragments[len(fragments)-1])

			var body map[string]any
			err = json.Unmarshal(content, &body)
			if err != nil {
				log.Fatalf("Error translating: %v", err)
			}

			req.URL.Path = fmt.Sprintf("%v/%v", req.URL.Path, body[key])
		} else {
			// Pao
			key := getIdField(fragments[len(fragments)-2])
			query := req.URL.Query()
			query.Add(key, strconv.Itoa(id))
			req.URL.RawQuery = query.Encode()
			req.Body = io.NopCloser(bytes.NewReader(content))
		}
		return
	}

	if req.Method == http.MethodPut {
		id, err := strconv.Atoi(fragments[len(fragments)-1])
		if err != nil {
			// Flan
			key := getIdField(fragments[len(fragments)-1])
			var body map[string]any
			err = json.Unmarshal(content, &body)
			if err != nil {
				log.Fatalf("Error translating: %v", err)
			}

			req.URL.Path = fmt.Sprintf("%v/%v", req.URL.Path, body[key])
		} else {
			// Pao
			req.URL.Path = strings.ReplaceAll(req.URL.Path, strconv.Itoa(id), "")
			req.Body = io.NopCloser(bytes.NewReader(content))
		}
	}
}

func getIdField(entity string) string {
	switch strings.TrimSpace(entity) {
	case "estudiantes", "profesores":
		return "cedula"
	case "asignaturas", "profesores_ciclo", "matriculas":
		return "id"
	default:
	}

	panic(fmt.Sprintf("Oh no!! its %v D:", entity))
}

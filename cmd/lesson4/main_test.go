package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleList(t *testing.T) {
	//req, err := http.NewRequest("GET", "/list?extension=.yaml", nil)
	req, err := http.NewRequest("GET", "/list", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := &ListHandler{"../../upload"}
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error(status, http.StatusOK)
	}

	expected := "Filename te.txt, size 0\nFilename test.yaml, size 21\n"
	if rr.Body.String() != expected {
		t.Errorf(rr.Body.String(), expected)
	}
}
func TestHandleListExtension(t *testing.T) {
	req, err := http.NewRequest("GET", "/list?extension=.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := &ListHandler{"../../upload"}
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error(status, http.StatusOK)
	}

	expected := "test.yaml"
	if rr.Body.String() != expected {
		t.Errorf(rr.Body.String(), expected)
	}
}

func TestUploadHandler(t *testing.T) {
	file, _ := os.Open("../../upload/test.yaml")
	defer file.Close()

	// действия, необходимые для того, чтобы засунуть файл в запрос
	// в качестве мультипарт-формы
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	_, err := io.Copy(part, file)
	if err != nil {
		t.Error(err)
	}
	writer.Close()

	// опять создаем запрос, теперь уже на /upload эндпоинт
	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// создаем ResponseRecorder
	rr := httptest.NewRecorder()

	// создаем заглушку файлового сервера. Для прохождения тестов
	// нам достаточно чтобы он возвращал 200 статус
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok!")
	}))
	defer ts.Close()

	uploadHandler := &UploadHandler{
		UploadDir: "../../upload",
		// таким образом мы подменим адрес файлового сервера
		// и вместо реального, хэндлер будет стучаться на заглушку
		// которая всегда будет возвращать 200 статус, что нам и нужна
		HostAddr: ts.URL,
	}

	// опять же, вызываем ServeHTTP у тестируемого обработчика
	uploadHandler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `test.yaml`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

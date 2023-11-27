package service

import (
	md52 "IB2/md5"
	"encoding/hex"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) GetHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/home", s.Home).Methods(http.MethodGet)
	router.HandleFunc("/about", s.About).Methods(http.MethodGet)
	router.HandleFunc("/home/upload", s.Upload).Methods(http.MethodPost)
	router.HandleFunc("/home/download", s.Download).Methods(http.MethodGet)
	return router
}

func (s *Service) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/home.html")
	//fmt.Fprint(w, "service/about.html")
}

func (s *Service) About(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/about.html")
}

func (s *Service) Upload(w http.ResponseWriter, r *http.Request) {
	md := md52.New()

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	hash := md.Md5Hash(fileBytes)
	//mymd5 := md5.New()
	//mymd5.Write(fileBytes)
	//hash2 := mymd5.Sum(nil)
	hashString := hex.EncodeToString(hash)
	processedFileName := handler.Filename + "_processed"
	processedFilePath := filepath.Join(".", "processed_files", processedFileName)
	if _, err := os.Stat("processed_files"); os.IsNotExist(err) {
		if err := os.Mkdir("processed_files", 0755); err != nil {
			log.Println(err)
			http.Error(w, "Error creating processed files directory", http.StatusInternalServerError)
			return
		}
	}
	processedFile, err := os.Create(processedFilePath)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error saving processed file", http.StatusInternalServerError)
		return
	}
	processedFile.Write([]byte(hashString))
	defer processedFile.Close()

	// Ваша логика обработки файла и сохранение его на сервере

	// Перенаправление на страницу скачивания
	http.Redirect(w, r, "/home/download?filename="+processedFileName, http.StatusSeeOther)
}

func (s *Service) Download(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(".", "processed_files", filename)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, filePath)
}

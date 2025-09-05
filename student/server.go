package student

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers(mux *http.ServeMux) {
	handler := StudentService{}
	mux.Handle("/students", handler)
	mux.Handle("/students/", handler)
}

type StudentService struct{}

// /students
// /students/{id}
// /students/{id}/grades
func (ss StudentService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	switch len(paths) {
	case 2:
		ss.getAll(w)
	case 3:
		id, err := strconv.Atoi(paths[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ss.getOne(w, id)
	case 4:
		id, err := strconv.Atoi(paths[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ss.addGrades(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (ss StudentService) getAll(w http.ResponseWriter) {
	stuMutex.Lock()
	defer stuMutex.Unlock()

	data, err := ss.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (ss StudentService) getOne(w http.ResponseWriter, id int) {
	stuMutex.Lock()
	defer stuMutex.Unlock()

	s, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := ss.toJSON(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (ss StudentService) addGrades(w http.ResponseWriter, r *http.Request, id int) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dec := json.NewDecoder(r.Body)
	var g Grade
	err := dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stuMutex.Lock()
	defer stuMutex.Unlock()

	s, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s.Grades = append(s.Grades, g)
	w.WriteHeader(http.StatusOK)
}

func (ss StudentService) toJSON(o interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(o)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

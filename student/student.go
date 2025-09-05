package student

import (
	"fmt"
	"sync"
)

type GradeType string

const (
	GradeQuiz = GradeType("quiz")
	GradeTest = GradeType("test")
	GradeExam = GradeType("exam")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}

type Student struct {
	Id     int
	Name   string
	Grades []Grade
}

func (s Student) AverageScore() float32 {
	var total float32 = 0.0
	for _, grade := range s.Grades {
		total += grade.Score
	}
	return total / float32(len(s.Grades))
}

func (s Student) TotalScore() float32 {
	var total float32 = 0.0
	for _, grade := range s.Grades {
		total += grade.Score
	}
	return total
}

type Students []Student

func (ss Students) GetById(id int) (*Student, error) {
	for _, s := range ss {
		if s.Id == id {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("student with id %d not found", id)
}

var (
	stuMutex sync.Mutex
	students Students
)

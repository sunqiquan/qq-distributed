package student

func init() {
	students = Students{
		{
			Id:   1,
			Name: "Alice",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 90.0,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 80.0,
				},
				{
					Title: "Exam 1",
					Type:  GradeExam,
					Score: 85.0,
				},
			},
		},
		{
			Id:   2,
			Name: "Bob",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 95.0,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 85.0,
				},
				{
					Title: "Exam 1",
					Type:  GradeExam,
					Score: 90.0,
				},
			},
		},
	}
}

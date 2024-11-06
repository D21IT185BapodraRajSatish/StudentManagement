package storage

import "github.com/D21IT185BapodraRajSatish/StudentAPI/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudent(id int64, name string, email string, age int) (int64, error)
	DeleteStudent(id int64) (int64, error)
}

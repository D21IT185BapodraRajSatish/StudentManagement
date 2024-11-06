package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/D21IT185BapodraRajSatish/StudentAPI/internal/config"
	"github.com/D21IT185BapodraRajSatish/StudentAPI/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE Table IF NOT EXISTS students(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students(name,email,age) VALUES(?,?,?)")

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastid, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var studnet types.Student

	err = stmt.QueryRow(id).Scan(&studnet.Id, &studnet.Name, &studnet.Email, &studnet.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return studnet, nil

}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) DeleteStudent(id int64) (int64, error) {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)

	if err != nil {
		return 0, err
	}
	rowsaffected, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsaffected, nil
}

func (s *Sqlite) UpdateStudent(id int64, name string, email string, age int) (int64, error) {
	_, err := s.GetStudentById(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("no student found with id %s", fmt.Sprint(id)) {
			return 0, fmt.Errorf("cannot update: %w", err) // student not found
		}
		return 0, err // other query error
	}
	
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return 0, err
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsaffected, nil
}

package main

import "fmt"

type Student struct {
	Name  string
	Age   int
	Grade string
}

func AddStudent(students []Student, student Student) []Student {
	return append(students, student)
}

func RemoveStudent(students []Student, index int) []Student {
	return append(students[:index], students[index+1:]...)
}

func FindStudent(students []Student, name string) (Student, bool) {
	for _, student := range students {
		if student.Name == name {
			return student, true
		}
	}
	return Student{}, false
}

func ListStudents(students []Student) {
	for _, student := range students {
		fmt.Printf("Name: %s, Age: %d, Grade: %s\n.", student.Name, student.Age, student.Grade)
	}
}

func main() {
	// Initialize an empty slice of students
	students := []Student{}

	// Add students
	students = AddStudent(students, Student{"Alice", 20, "A"})
	students = AddStudent(students, Student{"Bob", 22, "B"})
	students = AddStudent(students, Student{"Charlie", 21, "A"})

	// List all students
	fmt.Println("All students:")
	ListStudents(students)

	// Remove a student by index (removing Bob)
	students = RemoveStudent(students, 1)

	// List students after removal
	fmt.Println("\nStudents after removal:")
	ListStudents(students)

	// Find a student by name
	nameToFind := "Alice"
	student, found := FindStudent(students, nameToFind)
	if found {
		fmt.Printf("\nFound student: Name: %s, Age: %d, Grade: %s\n", student.Name, student.Age, student.Grade)
	} else {
		fmt.Printf("\nStudent %s not found\n", nameToFind)
	}
}

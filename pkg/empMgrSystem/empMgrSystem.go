package empMgrSystem

import "fmt"

type employee struct {
	id int
	name string
}

type ems struct {
	employee map[int]employee
}

func (e ems) showEmployee() {
	for _, v := range e.employee {
		fmt.Printf("ID: %d, Name: %s\n", v.id, v.name)
	}
}

func (e ems) addEmployee() {
	var (
		empID int
		empName string
	)
	fmt.Print("new employee ID is: ")
	fmt.Scanln(&empID)
	fmt.Print("new employee Name is: ")
	fmt.Scanln(&empName)
	newEmp := employee{
		id: empID,
		name: empName,
	}
	e.employee[newEmp.id] = newEmp
}

func (e ems) editEmployee() {
	var(
		empID int
		empName string
	)
	fmt.Print("which employee you want to edit? please provide the ID: ")
	fmt.Scanln(&empID)
	empObj, err := e.employee[empID]
	if !err {
		fmt.Println("there is no such employee")
		return
	}
	fmt.Printf("the employee you want to edit is - ID: %d, Name: %s\n", empObj.id, empObj.name)
	fmt.Print("please enter the new name: ")
	fmt.Scanln(&empName)
	empObj.name = empName
	e.employee[empID] = empObj
	fmt.Println("done")
}

func (e ems) deleteEmployee() {
	var empID int
	fmt.Print("which employee you want to delete? please provide the ID: ")
	fmt.Scanln(&empID)
	_, err := e.employee[empID]
	if !err {
		fmt.Println("there is no such employee")
		return
	}
	delete(e.employee, empID)
	fmt.Println("done")
}
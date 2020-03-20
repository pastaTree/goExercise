package main

import (
	"fmt"
	"os"
)

var (
	empMgrSys ems
)

func showMenu() {
	fmt.Println("welcome to Employee Management System")
	fmt.Println(`
		1. show all employees
		2. add employee
		3. edit employee
		4. delete employee
		5. exit
	`)
}

func main() {
	empMgrSys = ems{
		employee: make(map[int]employee, 100),
	}
	for {
		showMenu()
		fmt.Print("please make a choice: ")
		var choice int
		fmt.Scanln(&choice)
		fmt.Printf("your choice is: %d\n", choice)
		switch choice {
		case 1:
			empMgrSys.showEmployee()
		case 2:
			empMgrSys.addEmployee()
		case 3:
			empMgrSys.editEmployee()
		case 4:
			empMgrSys.deleteEmployee()
		case 5:
			os.Exit(1)
		default:
			fmt.Println("invalid choice, try again!")
		}
	}

}
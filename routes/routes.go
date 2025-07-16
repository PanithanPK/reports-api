package routes

import (
	"reports-api/handlers"

	"github.com/gorilla/mux"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(r *mux.Router) {

	// Problem routes
	r.HandleFunc("/problemEntry/reportProblem", handlers.GetTasksHandler).Methods("GET")
	r.HandleFunc("/problemEntry/reportProblem", handlers.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/problemEntry/reportProblem/{id}", handlers.UpdateTaskHandler).Methods("PUT")
	r.HandleFunc("/problemEntry/reportProblem/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

	// Phone routes
	r.HandleFunc("/phoneEntry/ipPhones", handlers.ListIPPhonesHandler).Methods("GET")
	r.HandleFunc("/phoneEntry/ipPhone", handlers.CreateIPPhoneHandler).Methods("POST")
	r.HandleFunc("/phoneEntry/ipPhone", handlers.UpdateIPPhoneHandler).Methods("PUT")
	r.HandleFunc("/phoneEntry/ipPhone/{id}", handlers.DeleteIPPhoneHandler).Methods("DELETE")

	// Program routes
	r.HandleFunc("/programEntry/programs", handlers.ListProgramsHandler).Methods("GET")
	r.HandleFunc("/programEntry/program", handlers.CreateProgramHandler).Methods("POST")
	r.HandleFunc("/programEntry/program/{id}", handlers.UpdateProgramHandler).Methods("PUT")
	r.HandleFunc("/programEntry/program/{id}", handlers.DeleteProgramHandler).Methods("DELETE")

	// Department routes
	r.HandleFunc("/departmentEntry/departments", handlers.ListDepartmentsHandler).Methods("GET")
	r.HandleFunc("/departmentEntry/department", handlers.CreateDepartmentHandler).Methods("POST")
	r.HandleFunc("/departmentEntry/department/{id}", handlers.UpdateDepartmentHandler).Methods("PUT")
	r.HandleFunc("/departmentEntry/department/{id}", handlers.DeleteDepartmentHandler).Methods("DELETE")

	// branch routes
	r.HandleFunc("/branchEntry/branches", handlers.ListBranchesHandler).Methods("GET")
	r.HandleFunc("/branchEntry/branch", handlers.CreateBranchHandler).Methods("POST")
	r.HandleFunc("/branchEntry/branch/{id}", handlers.UpdateBranchHandler).Methods("PUT")
	r.HandleFunc("/branchEntry/branch/{id}", handlers.DeleteBranchHandler).Methods("DELETE")

}

// RegisterAuthRoutes registers all authentication-related routes
func RegisterAuthRoutes(r *mux.Router) {
	// Authentication routes
	r.HandleFunc("/authEntry/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/authEntry/registerUser", handlers.RegisterHandler("user")).Methods("POST")
	r.HandleFunc("/authEntry/registerAdmin", handlers.RegisterHandler("admin")).Methods("POST")
	r.HandleFunc("/authEntry/updateUser", handlers.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/authEntry/deleteUser", handlers.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/authEntry/logout", handlers.LogoutHandler).Methods("POST")
}

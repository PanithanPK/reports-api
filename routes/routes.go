package routes

import (
	"reports-api/handlers"

	"github.com/gorilla/mux"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(r *mux.Router) {

	// Problem routes
	r.HandleFunc("/api/v1/problem/list", handlers.GetTasksHandler).Methods("GET")
	r.HandleFunc("/api/v1/problem/{id}", handlers.GetTaskDetailHandler).Methods("GET")
	r.HandleFunc("/api/v1/problem/create", handlers.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/api/v1/problem/update/{id}", handlers.UpdateTaskHandler).Methods("PUT")
	r.HandleFunc("/api/v1/problem/delete/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

	// Phone routes
	r.HandleFunc("/api/v1/ipphone/list", handlers.ListIPPhonesHandler).Methods("GET")
	r.HandleFunc("/api/v1/ipphone/create", handlers.CreateIPPhoneHandler).Methods("POST")
	r.HandleFunc("/api/v1/ipphone/update/{id}", handlers.UpdateIPPhoneHandler).Methods("PUT")
	r.HandleFunc("/api/v1/ipphone/delete/{id}", handlers.DeleteIPPhoneHandler).Methods("DELETE")

	// Program routes
	r.HandleFunc("/api/v1/program/list", handlers.ListProgramsHandler).Methods("GET")
	r.HandleFunc("/api/v1/program/create", handlers.CreateProgramHandler).Methods("POST")
	r.HandleFunc("/api/v1/program/update/{id}", handlers.UpdateProgramHandler).Methods("PUT")
	r.HandleFunc("/api/v1/program/delete/{id}", handlers.DeleteProgramHandler).Methods("DELETE")

	// Department routes
	r.HandleFunc("/api/v1/department/list", handlers.ListDepartmentsHandler).Methods("GET")
	r.HandleFunc("/api/v1/department/{id}", handlers.GetDepartmentDetailHandler).Methods("GET")
	r.HandleFunc("/api/v1/department/create", handlers.CreateDepartmentHandler).Methods("POST")
	r.HandleFunc("/api/v1/department/update/{id}", handlers.UpdateDepartmentHandler).Methods("PUT")
	r.HandleFunc("/api/v1/department/delete/{id}", handlers.DeleteDepartmentHandler).Methods("DELETE")

	// branch routes
	r.HandleFunc("/api/v1/branch/list", handlers.ListBranchesHandler).Methods("GET")
	r.HandleFunc("/api/v1/branch/{id}", handlers.GetBranchDetailHandler).Methods("GET")
	r.HandleFunc("/api/v1/branch/create", handlers.CreateBranchHandler).Methods("POST")
	r.HandleFunc("/api/v1/branch/update/{id}", handlers.UpdateBranchHandler).Methods("PUT")
	r.HandleFunc("/api/v1/branch/delete/{id}", handlers.DeleteBranchHandler).Methods("DELETE")

	//Dashboard routes
	r.HandleFunc("/api/v1/dashboard/data", handlers.GetDashboardDataHandler).Methods("GET")

	r.HandleFunc("/api/v1/telegramMessage", handlers.SendTelegramNotificationHandler).Methods("POST")

	r.HandleFunc("/api/v1/updateTaskStatus", handlers.UpdateTaskStatusHandler).Methods("PUT")
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

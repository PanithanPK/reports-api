package routes

import (
	"reports-api/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(r *fiber.App) {

	// Problem routes
	r.Get("/api/v1/problem/list", handlers.GetTasksHandler)
	r.Get("/api/v1/problem/list/:query", handlers.GetTasksWithQueryHandler)
	r.Get("/api/v1/problem/search/:query", handlers.SearchTasksHandler)
	r.Post("/api/v1/problem/create", handlers.CreateTaskHandler)
	r.Get("/api/v1/problem/:id", handlers.GetTaskDetailHandler)
	r.Put("/api/v1/problem/update/:id", handlers.UpdateTaskHandler)
	r.Delete("/api/v1/problem/delete/:id", handlers.DeleteTaskHandler)

	// Phone routes
	r.Get("/api/v1/ipphone/list", handlers.ListIPPhonesHandler)
	r.Get("/api/v1/ipphone/list/:query", handlers.ListIPPhonesQueryHandler)
	r.Get("/api/v1/ipphone/listall", handlers.AllIPPhonesHandler)
	r.Get("/api/v1/ipphone/search/:query", handlers.SearchIPPhonesHandler)
	r.Post("/api/v1/ipphone/create", handlers.CreateIPPhoneHandler)
	r.Put("/api/v1/ipphone/update/:id", handlers.UpdateIPPhoneHandler)
	r.Delete("/api/v1/ipphone/delete/:id", handlers.DeleteIPPhoneHandler)

	// Program routes
	r.Get("/api/v1/program/list", handlers.ListProgramsHandler)
	r.Get("/api/v1/program/list/:query", handlers.ListProgramsQueryHandler)
	r.Post("/api/v1/program/create", handlers.CreateProgramHandler)
	r.Put("/api/v1/program/update/:id", handlers.UpdateProgramHandler)
	r.Delete("/api/v1/program/delete/:id", handlers.DeleteProgramHandler)

	// Department routes
	r.Get("/api/v1/department/list", handlers.ListDepartmentsHandler)
	r.Get("/api/v1/department/list/:query", handlers.ListDepartmentsQueryHandler)
	r.Get("/api/v1/department/listall", handlers.AllDepartmentsHandler)
	r.Get("/api/v1/department/search/:query", handlers.SearchDepartmentsHandler)
	r.Post("/api/v1/department/create", handlers.CreateDepartmentHandler)
	r.Get("/api/v1/department/:id", handlers.GetDepartmentDetailHandler)
	r.Put("/api/v1/department/update/:id", handlers.UpdateDepartmentHandler)
	r.Delete("/api/v1/department/delete/:id", handlers.DeleteDepartmentHandler)

	// branch routes
	r.Get("/api/v1/branch/list", handlers.ListBranchesHandler)
	r.Get("/api/v1/branch/list/:query", handlers.ListBranchesQueryHandler)
	r.Post("/api/v1/branch/create", handlers.CreateBranchHandler)
	r.Get("/api/v1/branch/:id", handlers.GetBranchDetailHandler)
	r.Put("/api/v1/branch/update/:id", handlers.UpdateBranchHandler)
	r.Delete("/api/v1/branch/delete/:id", handlers.DeleteBranchHandler)

	//Dashboard routes
	r.Get("/api/v1/dashboard/data", handlers.GetDashboardDataHandler)

	r.Put("/api/v1/updateTaskStatus", handlers.UpdateTaskStatusHandler)

	r.Get("/api/v1/scores/list", handlers.ListScoresHandler)
	r.Get("/api/v1/scores/:id", handlers.GetScoreDetailHandler)
	r.Put("/api/v1/scores/update/:id", handlers.UpdateScoreHandler)
	r.Delete("/api/v1/scores/delete/:id", handlers.DeleteScoreHandler)

	r.Get("/api/v1/respons/list", handlers.GetresponsHandler)
	r.Post("/api/v1/respons/create", handlers.AddresponsHandler)
}

// RegisterAuthRoutes registers all authentication-related routes
func RegisterAuthRoutes(r *fiber.App) {
	// Authentication routes
	r.Post("/api/authEntry/login", handlers.LoginHandler)
	r.Post("/api/authEntry/registerUser", handlers.RegisterHandler("user"))
	r.Post("/api/authEntry/registerAdmin", handlers.RegisterHandler("admin"))
	r.Put("/api/authEntry/updateUser", handlers.UpdateUserHandler)
	r.Delete("/api/authEntry/deleteUser", handlers.DeleteUserHandler)
	r.Post("/api/authEntry/logout", handlers.LogoutHandler)
}

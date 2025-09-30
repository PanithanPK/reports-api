package routes

import (
	"reports-api/handlers"

	"github.com/gofiber/fiber/v2"
)

// MainRoutes registers all API routes
func MainRoutes(r *fiber.App) {
	//Dashboard routes
	r.Get("/api/v1/dashboard/data", handlers.GetDashboardDataHandler)

	//Data export routes
	r.Get("/api/v1/dashboard/data/phonecsv", handlers.IpphonesExportCsv)
	r.Get("/api/v1/dashboard/data/departmetcsv", handlers.DepartmentsExportCsv)
	r.Get("/api/v1/dashboard/data/branchcsv", handlers.BranchExportCsv)
	r.Get("/api/v1/dashboard/data/systemcsv", handlers.SystemExportCsv)
	r.Get("/api/v1/dashboard/data/taskscsv", handlers.TasksExportCsv)

	r.Get("/api/v1/scores/list", handlers.ListScoresHandler)
	r.Get("/api/v1/scores/:id", handlers.GetScoreDetailHandler)
	r.Put("/api/v1/scores/update/:id", handlers.UpdateScoreHandler)
	r.Delete("/api/v1/scores/delete/:id", handlers.DeleteScoreHandler)

	r.Get("/api/v1/respons/list", handlers.GetresponsHandler)
	r.Get("/api/v1/respons/:id", handlers.GetResponsDetailHandler)
	r.Post("/api/v1/respons/create", handlers.AddresponsHandler)
	r.Put("/api/v1/respons/update/:id", handlers.UpdateResponsHandler)
	r.Delete("/api/v1/respons/delete/:id", handlers.DeleteResponsHandler)
}

// problemRoutes registers all problem-related routes
func problemRoutes(r *fiber.App) {
	r.Get("/api/v1/problem/list", handlers.GetTasksHandler)
	r.Get("/api/v1/problem/list/:query", handlers.GetTasksWithQueryHandler)
	r.Get("/api/v1/problem/list/:column/:query", handlers.GetTasksWithColumnQueryHandler)
	r.Get("/api/v1/problem/list/sort/:column/:query", handlers.GetTaskSort)
	r.Post("/api/v1/problem/create", handlers.CreateTaskHandler)
	r.Get("/api/v1/problem/:id", handlers.GetTaskDetailHandler)
	r.Put("/api/v1/problem/update/:id", handlers.UpdateTaskHandler)
	r.Delete("/api/v1/problem/delete/:id", handlers.DeleteTaskHandler)
	r.Put("/api/v1/problem/update/assignto/:id", handlers.UpdateAssignedTo)
}

// resolutionRoutes registers all resolution-related routes
func resolutionRoutes(r *fiber.App) {
	r.Get("/api/v1/resolution/:id", handlers.GetResolutionHandler)
	r.Post("/api/v1/resolution/create/:id", handlers.CreateResolutionHandler)
	r.Put("/api/v1/resolution/update/:id", handlers.UpdateResolutionHandler)
	r.Delete("/api/v1/resolution/delete/:id", handlers.DeleteResolutionHandler)
}

// progressRoutes registers all progress-related routes
func progressRoutes(r *fiber.App) {
	r.Get("/api/v1/progress/:id", handlers.GetProgressHandler)
	r.Post("/api/v1/progress/create/:id", handlers.CreateProgressHandler)
	r.Put("/api/v1/progress/update/:id/:pgid", handlers.UpdateProgressHandler)
	r.Delete("/api/v1/progress/delete/:id/:pgid", handlers.DeleteProgressHandler)
}

// ipphoneRoutes registers all IP phone-related routes
func ipphoneRoutes(r *fiber.App) {
	r.Get("/api/v1/ipphone/list", handlers.ListIPPhonesHandler)
	r.Get("/api/v1/ipphone/list/:query", handlers.ListIPPhonesQueryHandler)
	r.Get("/api/v1/ipphone/:id", handlers.GetIPPhonesDetailHandler)
	r.Get("/api/v1/ipphone/listall", handlers.AllIPPhonesHandler)
	r.Post("/api/v1/ipphone/create", handlers.CreateIPPhoneHandler)
	r.Put("/api/v1/ipphone/update/:id", handlers.UpdateIPPhoneHandler)
	r.Delete("/api/v1/ipphone/delete/:id", handlers.DeleteIPPhoneHandler)
}

// programRoutes registers all program-related routes
func programRoutes(r *fiber.App) {
	r.Get("/api/v1/program/list", handlers.ListProgramsHandler)
	r.Get("/api/v1/program/list/:query", handlers.ListProgramsQueryHandler)
	r.Post("/api/v1/program/create", handlers.CreateProgramHandler)
	r.Get("/api/v1/program/type/list", handlers.GETTypeProgramHandler)
	r.Get("/api/v1/program/type/list/:query", handlers.GetTypeWithQueryHandler)
	r.Post("/api/v1/program/type/create", handlers.AddTypeProgramHandler)
	r.Post("/api/v1/program/type/update/:id", handlers.UpdateTypeProgramHandler)
	r.Delete("/api/v1/program/type/delete/:id", handlers.DeleteTypeHandler)
	r.Get("/api/v1/program/:id", handlers.GetProgramDetailHandler)
	r.Put("/api/v1/program/update/:id", handlers.UpdateProgramHandler)
	r.Delete("/api/v1/program/delete/:id", handlers.DeleteProgramHandler)
}

// departmentRoutes registers all department-related routes
func departmentRoutes(r *fiber.App) {
	r.Get("/api/v1/department/list", handlers.ListDepartmentsHandler)
	r.Get("/api/v1/department/list/:query", handlers.ListDepartmentsQueryHandler)
	r.Get("/api/v1/department/listall", handlers.AllDepartmentsHandler)
	r.Post("/api/v1/department/create", handlers.CreateDepartmentHandler)
	r.Get("/api/v1/department/:id", handlers.GetDepartmentDetailHandler)
	r.Put("/api/v1/department/update/:id", handlers.UpdateDepartmentHandler)
	r.Delete("/api/v1/department/delete/:id", handlers.DeleteDepartmentHandler)
}

// branchRoutes registers all branch-related routes
func branchRoutes(r *fiber.App) {
	r.Get("/api/v1/branch/list", handlers.ListBranchesHandler)
	r.Get("/api/v1/branch/list/:query", handlers.ListBranchesQueryHandler)
	r.Post("/api/v1/branch/create", handlers.CreateBranchHandler)
	r.Get("/api/v1/branch/:id", handlers.GetBranchDetailHandler)
	r.Put("/api/v1/branch/update/:id", handlers.UpdateBranchHandler)
	r.Delete("/api/v1/branch/delete/:id", handlers.DeleteBranchHandler)
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

	// User management routes
	r.Get("/api/authEntry/users", handlers.GetAllUsersHandler)
	r.Get("/api/authEntry/user/:id", handlers.GetUserDetailHandler)
}

// RegisterRoutes registers all routes
func RegisterRoutes(r *fiber.App) {
	RegisterAuthRoutes(r)
	MainRoutes(r)
	problemRoutes(r)
	resolutionRoutes(r)
	progressRoutes(r)
	ipphoneRoutes(r)
	programRoutes(r)
	departmentRoutes(r)
	branchRoutes(r)
}

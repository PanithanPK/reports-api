package routes

import (
	"reports-api/handlers"
	"reports-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// MainRoutes registers all API routes
func MainRoutes(r *fiber.App) {
	// Dashboard routes (with session validation)
	r.Get("/api/v1/dashboard/data", middleware.SessionMiddleware(), handlers.GetDashboardDataHandler)

	// Data export routes
	r.Get("/api/v1/dashboard/data/phonecsv", middleware.SessionMiddleware(), handlers.IpphonesExportCsv)
	r.Get("/api/v1/dashboard/data/departmetcsv", middleware.SessionMiddleware(), handlers.DepartmentsExportCsv)
	r.Get("/api/v1/dashboard/data/branchcsv", middleware.SessionMiddleware(), handlers.BranchExportCsv)
	r.Get("/api/v1/dashboard/data/systemcsv", middleware.SessionMiddleware(), handlers.SystemExportCsv)
	r.Get("/api/v1/dashboard/data/taskscsv", middleware.SessionMiddleware(), handlers.TasksExportCsv)

	r.Get("/api/v1/scores/list", handlers.ListScoresHandler)
	r.Get("/api/v1/scores/:id", handlers.GetScoreDetailHandler)
	r.Put("/api/v1/scores/update/:id", middleware.SessionMiddleware(), handlers.UpdateScoreHandler)
	r.Delete("/api/v1/scores/delete/:id", middleware.SessionMiddleware(), handlers.DeleteScoreHandler)

	r.Get("/api/v1/respons/list", handlers.GetResponsibilitiesHandler)
	r.Get("/api/v1/respons/:id", handlers.GetResponsibilityDetailHandler)
	r.Post("/api/v1/respons/create", middleware.SessionMiddleware(), handlers.AddResponsibilityHandler)
	r.Put("/api/v1/respons/update/:id", middleware.SessionMiddleware(), handlers.UpdateResponsibilityHandler)
	r.Delete("/api/v1/respons/delete/:id", middleware.SessionMiddleware(), handlers.DeleteResponsibilityHandler)
}

// problemRoutes registers all problem-related routes
func problemRoutes(r *fiber.App) {
	r.Get("/api/v1/problem/list", handlers.GetTasksHandler)
	r.Get("/api/v1/problem/list/:query", handlers.GetTasksWithQueryHandler)
	r.Get("/api/v1/problem/list/:column/:query", handlers.GetTasksWithColumnQueryHandler)
	r.Get("/api/v1/problem/list/sort/:column/:query", middleware.SessionMiddleware(), handlers.GetTaskSort)
	r.Post("/api/v1/problem/create", middleware.SessionMiddleware(), middleware.RateLimiter(), handlers.CreateTaskHandler)
	r.Get("/api/v1/problem/:id", handlers.GetTaskDetailHandler)
	r.Put("/api/v1/problem/update/:id", middleware.SessionMiddleware(), middleware.RateLimiter(), handlers.UpdateTaskHandler)
	r.Delete("/api/v1/problem/delete/:id", middleware.SessionMiddleware(), middleware.RateLimiter(), handlers.DeleteTaskHandler)
	r.Put("/api/v1/problem/update/assignto/:id", middleware.SessionMiddleware(), handlers.UpdateAssignedTo)
}

// resolutionRoutes registers all resolution-related routes
func resolutionRoutes(r *fiber.App) {
	r.Get("/api/v1/resolution/:id", handlers.GetResolutionHandler)
	r.Post("/api/v1/resolution/create/:id", middleware.SessionMiddleware(), handlers.CreateResolutionHandler)
	r.Put("/api/v1/resolution/update/:id", middleware.SessionMiddleware(), handlers.UpdateResolutionHandler)
	r.Delete("/api/v1/resolution/delete/:id", middleware.SessionMiddleware(), handlers.DeleteResolutionHandler)
}

// progressRoutes registers all progress-related routes
func progressRoutes(r *fiber.App) {
	r.Get("/api/v1/progress/:id", handlers.GetProgressHandler)
	r.Post("/api/v1/progress/create/:id", middleware.SessionMiddleware(), handlers.CreateProgressHandler)
	r.Put("/api/v1/progress/update/:id/:pgid", middleware.SessionMiddleware(), handlers.UpdateProgressHandler)
	r.Delete("/api/v1/progress/delete/:id/:pgid", middleware.SessionMiddleware(), handlers.DeleteProgressHandler)
}

// ipphoneRoutes registers all IP phone-related routes
func ipphoneRoutes(r *fiber.App) {
	r.Get("/api/v1/ipphone/list", handlers.ListIPPhonesHandler)
	r.Get("/api/v1/ipphone/list/:query", handlers.ListIPPhonesQueryHandler)
	r.Get("/api/v1/ipphone/:id", handlers.GetIPPhonesDetailHandler)
	r.Get("/api/v1/ipphone/listall", handlers.AllIPPhonesHandler)
	r.Post("/api/v1/ipphone/create", middleware.SessionMiddleware(), handlers.CreateIPPhoneHandler)
	r.Put("/api/v1/ipphone/update/:id", middleware.SessionMiddleware(), handlers.UpdateIPPhoneHandler)
	r.Delete("/api/v1/ipphone/delete/:id", middleware.SessionMiddleware(), handlers.DeleteIPPhoneHandler)
}

// programRoutes registers all program-related routes
func programRoutes(r *fiber.App) {
	r.Get("/api/v1/program/list", handlers.ListProgramsHandler)
	r.Get("/api/v1/program/list/:query", handlers.ListProgramsQueryHandler)
	r.Post("/api/v1/program/create", middleware.SessionMiddleware(), handlers.CreateProgramHandler)
	r.Get("/api/v1/program/type/list", handlers.GETTypeProgramHandler)
	r.Get("/api/v1/program/type/list/:query", handlers.GetTypeWithQueryHandler)
	r.Post("/api/v1/program/type/create", middleware.SessionMiddleware(), handlers.AddTypeProgramHandler)
	r.Put("/api/v1/program/type/update/:id", middleware.SessionMiddleware(), handlers.UpdateTypeProgramHandler)
	r.Delete("/api/v1/program/type/delete/:id", middleware.SessionMiddleware(), handlers.DeleteTypeHandler)
	r.Get("/api/v1/program/:id", handlers.GetProgramDetailHandler)
	r.Put("/api/v1/program/update/:id", middleware.SessionMiddleware(), handlers.UpdateProgramHandler)
	r.Delete("/api/v1/program/delete/:id", middleware.SessionMiddleware(), handlers.DeleteProgramHandler)
}

// departmentRoutes registers all department-related routes
func departmentRoutes(r *fiber.App) {
	r.Get("/api/v1/department/list", handlers.ListDepartmentsHandler)
	r.Get("/api/v1/department/list/:query", handlers.ListDepartmentsQueryHandler)
	r.Get("/api/v1/department/listall", handlers.AllDepartmentsHandler)
	r.Post("/api/v1/department/create", middleware.SessionMiddleware(), handlers.CreateDepartmentHandler)
	r.Get("/api/v1/department/:id", handlers.GetDepartmentDetailHandler)
	r.Put("/api/v1/department/update/:id", middleware.SessionMiddleware(), handlers.UpdateDepartmentHandler)
	r.Delete("/api/v1/department/delete/:id", middleware.SessionMiddleware(), handlers.DeleteDepartmentHandler)
}

// branchRoutes registers all branch-related routes
func branchRoutes(r *fiber.App) {
	r.Get("/api/v1/branch/list", handlers.ListBranchesHandler)
	r.Get("/api/v1/branch/list/:query", handlers.ListBranchesQueryHandler)
	r.Post("/api/v1/branch/create", middleware.SessionMiddleware(), handlers.CreateBranchHandler)
	r.Get("/api/v1/branch/:id", handlers.GetBranchDetailHandler)
	r.Put("/api/v1/branch/update/:id", middleware.SessionMiddleware(), handlers.UpdateBranchHandler)
	r.Delete("/api/v1/branch/delete/:id", middleware.SessionMiddleware(), handlers.DeleteBranchHandler)
}

// RegisterAuthRoutes registers all authentication-related routes
func RegisterAuthRoutes(r *fiber.App) {
	// Authentication routes with rate limiting
	r.Post("/api/authEntry/login", middleware.RateLimiter(), handlers.LoginHandler)
	r.Post("/api/authEntry/registerUser", middleware.SessionMiddleware(), middleware.RateLimiter(), handlers.RegisterHandler("user"))
	r.Post("/api/authEntry/registerAdmin", middleware.SessionMiddleware(), middleware.RateLimiter(), handlers.RegisterHandler("admin"))
	r.Put("/api/authEntry/updateUser", middleware.SessionMiddleware(), middleware.RateLimiter(), handlers.UpdateUserHandler)
	r.Delete("/api/authEntry/deleteUser", middleware.SessionMiddleware(), middleware.RateLimiter(), handlers.DeleteUserHandler)
	r.Post("/api/authEntry/logout", handlers.LogoutHandler)

	// User management routes
	r.Get("/api/authEntry/users", middleware.RateLimiter(), handlers.GetAllUsersHandler)
	r.Get("/api/authEntry/user/:id", middleware.RateLimiter(), handlers.GetUserDetailHandler)
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

package routes

import (
	"github.com/KalinduBihan/leave-management-api/config"
	"github.com/KalinduBihan/leave-management-api/internal/handler"
	"github.com/KalinduBihan/leave-management-api/internal/middleware"
	"github.com/KalinduBihan/leave-management-api/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all application routes
func SetupRoutes(
	router *gin.Engine,
	services *service.Services,
	cfg *config.Config,
) {
	// Apply global middleware
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// Health check (no auth required)
	healthHandler := handler.NewHealthHandler()
	router.GET("/health", healthHandler.HealthCheck)

	// Public routes
	authHandler := handler.NewAuthHandler(services.Auth)
	public := router.Group("/api/v1/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(services.Auth))
	{
		// Auth routes
		authRoutes := protected.Group("/auth")
		{
			authRoutes.POST("/refresh", authHandler.RefreshToken)
		}

		// Employee routes
		employeeHandler := handler.NewEmployeeHandler(services.Employee)
		employeeRoutes := protected.Group("/employees")
		{
			employeeRoutes.POST("", employeeHandler.CreateEmployee)
			employeeRoutes.GET("", employeeHandler.ListEmployees)
			employeeRoutes.GET("/:id", employeeHandler.GetEmployee)
			employeeRoutes.PUT("/:id", employeeHandler.UpdateEmployee)
			employeeRoutes.DELETE("/:id", employeeHandler.DeleteEmployee)
			employeeRoutes.GET("/:id/balance", employeeHandler.GetLeaveBalance)
		}

		// Department routes
		departmentHandler := handler.NewDepartmentHandler(services.Department)
		departmentRoutes := protected.Group("/departments")
		{
			departmentRoutes.POST("", departmentHandler.CreateDepartment)
			departmentRoutes.GET("", departmentHandler.ListDepartments)
			departmentRoutes.GET("/:id", departmentHandler.GetDepartment)
			departmentRoutes.PUT("/:id", departmentHandler.UpdateDepartment)
			departmentRoutes.DELETE("/:id", departmentHandler.DeleteDepartment)
		}

		// Leave type routes
		leaveTypeHandler := handler.NewLeaveTypeHandler(services.LeaveType)
		leaveTypeRoutes := protected.Group("/leave-types")
		{
			leaveTypeRoutes.POST("", leaveTypeHandler.CreateLeaveType)
			leaveTypeRoutes.GET("", leaveTypeHandler.ListLeaveTypes)
			leaveTypeRoutes.GET("/:id", leaveTypeHandler.GetLeaveType)
			leaveTypeRoutes.PUT("/:id", leaveTypeHandler.UpdateLeaveType)
			leaveTypeRoutes.DELETE("/:id", leaveTypeHandler.DeleteLeaveType)
		}

		// Leave request routes
		leaveRequestHandler := handler.NewLeaveRequestHandler(services.LeaveRequest)
		leaveRoutes := protected.Group("/leaves")
		{
			leaveRoutes.POST("", leaveRequestHandler.CreateLeaveRequest)
			leaveRoutes.GET("", leaveRequestHandler.ListLeaveRequests)
			leaveRoutes.GET("/:id", leaveRequestHandler.GetLeaveRequest)
			leaveRoutes.PUT("/:id", leaveRequestHandler.UpdateLeaveRequest)
			leaveRoutes.DELETE("/:id", leaveRequestHandler.DeleteLeaveRequest)
			leaveRoutes.POST("/:id/approve", leaveRequestHandler.ApproveLeaveRequest)
			leaveRoutes.POST("/:id/reject", leaveRequestHandler.RejectLeaveRequest)
		}

		// Employee leave routes
		leaveRoutes.GET("/employees/:employee_id/leaves", leaveRequestHandler.GetEmployeeLeaveRequests)

		// Manager approval routes
		managerRoutes := protected.Group("/managers")
		{
			managerRoutes.GET("/:manager_id/pending-approvals", leaveRequestHandler.GetPendingApprovals)
		}

		// Dashboard routes
		dashboardHandler := handler.NewDashboardHandler(services.Dashboard)
		dashboardRoutes := protected.Group("/dashboard")
		{
			dashboardRoutes.GET("/stats", dashboardHandler.GetDashboardStats)
			dashboardRoutes.GET("/employees/:employee_id/stats", dashboardHandler.GetEmployeeStats)
		}
	}
}
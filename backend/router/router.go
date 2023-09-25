package router

import (
	"backend/controllers"
	mw "backend/middleware"
	"backend/pkg/v1/mysql"
	"backend/utils/config"

	adminDelivv1 "backend/api/admin/v1/delivery"
	adminRepov1 "backend/api/admin/v1/repositories"
	adminUsecasev1 "backend/api/admin/v1/usecase"

	healthcheck "github.com/RaMin0/gin-health-check"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// CreateDefaultRouter is a function to create initiation router
func CreateRouter(isDev bool) *gin.Engine {
	// Swagger setup
	host := config.MyConfig.Host
	port := config.MyConfig.ServerPort
	urlSwagger := ginSwagger.URL(host + port + "/swagger/doc.json")
	// Create path url
	router := gin.New()
	// Use middleware
	router.Use(mw.Secure(isDev))
	router.Use(mw.CORSMiddleware())
	router.Use(requestid.New())
	router.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		controllers.HandlePanic(c, err)
	}))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, urlSwagger))
	router.Use(healthcheck.Default())
	router.Use(mw.RequestLoggerActivity())
	return router
}

func InitRouteV1_0_0(router *gin.Engine) {
	//v1Private := router.Group("/v1.0/private")
	v1Admin := router.Group("/v1.0/admin")
	v1Public := router.Group("/v1.0/public")

	db, err := mysql.GetConnectionItem()
	if err != nil {
		return
	}

	v1Admin.Use(controllers.MiddlewareFuncOverrideAdmin())

	// repositories
	ir := adminRepov1.NewTestRepoAdmin(db)

	// usecase
	iu := adminUsecasev1.NewAdminUsecase(ir)

	// handler
	adminDelivv1.NewAdminController(v1Public, v1Admin, iu)
}

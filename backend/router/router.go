package router

import (
	"backend/controllers"
	mw "backend/middleware"
	"backend/pkg/v1/mysql"
	"backend/utils/config"

	authDelivv1 "backend/api/manage/v1/auth/delivery"
	authRepov1 "backend/api/manage/v1/auth/repositories"
	authUsecasev1 "backend/api/manage/v1/auth/usecase"

	itemDelivv1 "backend/api/manage/v1/item/delivery"
	itemRepov1 "backend/api/manage/v1/item/repositories"
	itemUsecasev1 "backend/api/manage/v1/item/usecase"

	utilDeelivv1 "backend/api/manage/v1/utilities/delivery"
	utilRepov1 "backend/api/manage/v1/utilities/repositories"
	utilUsecasev1 "backend/api/manage/v1/utilities/usecase"

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
	v1Private := router.Group("/v1.0/private")
	v1Public := router.Group("/v1.0/public")

	db, err := mysql.GetConnectionItem()
	if err != nil {
		return
	}

	v1Private.Use(controllers.MiddlewareFuncOverrideUser())

	// repositories
	ar := authRepov1.NewTestRepoAuth(db)
	ir := itemRepov1.NewTestRepoItem(db)
	ur := utilRepov1.NewTestRepoUtilities(db)

	// usecase
	au := authUsecasev1.NewAuthUsecase(ar)
	iu := itemUsecasev1.NewItemUsecase(ir)
	uu := utilUsecasev1.NewUtilitiesUsecase(ur)

	// handler
	authDelivv1.NewAuthController(v1Public, v1Private, au)
	itemDelivv1.NewItemController(v1Private, iu)
	utilDeelivv1.NewUtilitiesController(v1Private, uu)
}

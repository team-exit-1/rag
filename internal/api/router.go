package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "refo-rag-server/docs"
	"refo-rag-server/internal/api/handler"
	"refo-rag-server/internal/service"
	"refo-rag-server/internal/storage"
)

// Router configures all API routes
func Router(conversationService *service.ConversationService, postgresStore storage.PostgresStoreInterface, qdrantStore storage.QdrantStoreInterface) *gin.Engine {
	router := gin.Default()

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// RAG API routes
	rag := router.Group("/api/rag")
	{
		// Health check endpoint
		healthHandler := handler.NewHealthCheckHandler(postgresStore, qdrantStore)
		rag.GET("/health", healthHandler.Handle)

		// Save conversation endpoint
		saveHandler := handler.NewSaveConversationHandler(conversationService)
		rag.POST("/conversation/store", saveHandler.Handle)

		// Search conversations endpoint
		searchHandler := handler.NewSearchConversationHandler(conversationService)
		rag.GET("/conversation/search", searchHandler.Handle)

		// Personal information endpoints
		personalInfoHandler := handler.NewPersonalInfoHandler(postgresStore)
		rag.POST("/personal-info", personalInfoHandler.CreatePersonalInfo)
		rag.GET("/personal-info/:info_id", personalInfoHandler.GetPersonalInfo)
		rag.GET("/personal-info/user/:user_id", personalInfoHandler.GetPersonalInfoByUser)
		rag.PUT("/personal-info/:info_id", personalInfoHandler.UpdatePersonalInfo)
		rag.DELETE("/personal-info/:info_id", personalInfoHandler.DeletePersonalInfo)
	}

	return router
}

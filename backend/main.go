package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"recruitment-platform/handlers"
	"recruitment-platform/store"
)

func main() {
	store.InitMockData()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		jobs := api.Group("/jobs")
		{
			jobs.GET("", handlers.ListJobs)
			jobs.GET("/:id", handlers.GetJob)
			jobs.POST("", handlers.CreateJob)
			jobs.PUT("/:id", handlers.UpdateJob)
		}

		applications := api.Group("/applications")
		{
			applications.GET("", handlers.ListApplications)
			applications.GET("/:id", handlers.GetApplication)
			applications.POST("", handlers.CreateApplication)
			applications.PUT("/:id/status", handlers.UpdateApplicationStatus)
		}

		messages := api.Group("/messages")
		{
			messages.GET("/application/:applicationId", handlers.ListMessages)
			messages.POST("", handlers.SendMessage)
		}

		stats := api.Group("/stats")
		{
			stats.GET("", handlers.GetStats)
		}

		requirements := api.Group("/requirements")
		{
			requirements.GET("", handlers.ListRequirements)
			requirements.GET("/stats", handlers.GetRequirementStats)
			requirements.GET("/:id", handlers.GetRequirement)
			requirements.POST("", handlers.CreateRequirement)
			requirements.PUT("/:id", handlers.UpdateRequirement)
			requirements.PUT("/:id/submit-approval", handlers.SubmitApproval)
			requirements.PUT("/:id/approve", handlers.ApproveRequirement)
			requirements.POST("/:id/convert-to-job", handlers.ConvertToJob)
		}

		candidates := api.Group("/candidates")
		{
			candidates.GET("", handlers.ListCandidates)
			candidates.GET("/detail", handlers.GetCandidate)
			candidates.PUT("/status", handlers.UpdateCandidateStatus)
		}

		notes := api.Group("/notes")
		{
			notes.GET("", handlers.ListNotes)
			notes.POST("", handlers.CreateNote)
		}

		interviews := api.Group("/interviews")
		{
			interviews.GET("", handlers.ListInterviews)
			interviews.GET("/:id", handlers.GetInterview)
			interviews.POST("", handlers.CreateInterview)
			interviews.PUT("/:id/reschedule", handlers.RescheduleInterview)
			interviews.PUT("/:id/cancel", handlers.CancelInterview)
			interviews.PUT("/:id/complete", handlers.CompleteInterview)
			interviews.PUT("/:id/note", handlers.AddInterviewNote)
		}

		offers := api.Group("/offers")
		{
			offers.GET("", handlers.ListOffers)
			offers.GET("/:id", handlers.GetOffer)
			offers.POST("", handlers.CreateOffer)
			offers.PUT("/:id/send", handlers.SendOffer)
			offers.PUT("/:id/reply", handlers.ReplyOffer)
			offers.PUT("/:id/withdraw", handlers.WithdrawOffer)
			offers.PUT("/:id/note", handlers.UpdateOfferNote)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

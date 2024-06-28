package main

import (
	"log"
	"net/http"
	"syndication-go-template/client/cloudesire"
	"time"

	"github.com/gin-gonic/gin"
)

type Event struct {
	Entity    string    `json:"entity" binding:"required"`
	ID        int       `json:"id" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	EntityURL string    `json:"entityUrl"`
	Date      time.Time `json:"date"`
}

func postEvent(r *gin.Engine) *gin.Engine {
	r.POST("event", func(c *gin.Context) {
		var event Event
		c.BindJSON(&event)
		log.Println("Received notification:", event)

		switch event.Entity {
		case "Subscription":
			subscription := cloudesire.GetSubscription(event.ID)

			switch event.Type {
			case "CREATED", "MODIFIED":
				deploy(subscription)
			case "DELETED":
				undeploy(subscription)
			}
		default:
			log.Printf("Skipping %s events", event.Entity)
		}

		c.Status(http.StatusNoContent)
	})

	return r
}

func deploy(subscription cloudesire.Subscription) {
	switch subscription.DeploymentStatus {
	case "PENDING":
		if subscription.Paid {
			log.Println("Provision tenant resources")
			cloudesire.UpdateSubscription(subscription.ID, "DEPLOYED")
		}
	case "STOPPED":
		log.Println("Temporarily suspend the subscription")
	case "DEPLOYED":
		log.Println("Check if tenant is OK")
	}
}

func undeploy(subscription cloudesire.Subscription) {
	log.Println("Unprovision tenant and release resources")
	cloudesire.UpdateSubscription(subscription.ID, "UNDEPLOYED")
}

func main() {
	r := gin.Default()
	r = postEvent(r)
	r.Run(":8080")
}

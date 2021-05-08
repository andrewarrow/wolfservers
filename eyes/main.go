package main

import (
	"os/exec"

	"github.com/gin-gonic/gin"
)

func StartEyes() {
	r := gin.Default()
	r.GET("/hi", func(c *gin.Context) {
		o, _ := exec.Command("ls", "-l", "/root/cardano-my-node/").Output()

		c.JSON(200, gin.H{
			"message": string(o),
		})
	})
	r.Run()
}

func main() {
	StartEyes()
}

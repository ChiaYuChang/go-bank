package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Empty struct{}

func (t Empty) Get(c *gin.Context)    { log.Println("function not yet implemented") }
func (t Empty) List(c *gin.Context)   { log.Println("function not yet implemented") }
func (t Empty) Create(c *gin.Context) { log.Println("function not yet implemented") }
func (t Empty) Update(c *gin.Context) { log.Println("function not yet implemented") }
func (t Empty) Delete(c *gin.Context) { log.Println("function not yet implemented") }

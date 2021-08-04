package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger(), gee.Recovery()) // global midlleware
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.GET("/renchen", func(c *gee.Context) {
		temp := "renchen"
		c.String(http.StatusOK, "hello world %s", temp[12]);
	})

	r.Run(":9999")
}
//
//func main() {
//	r := gee.New()
//	r.GET("/", func(c *gee.Context) {
//		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
//	})
//
//	r.GET("/hello", func(c *gee.Context) {
//		// expect /hello?name=geektutu
//		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
//	})
//
//	r.GET("/hello/:name", func(c *gee.Context) {
//		// expect /hello/geektutu
//		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
//	})
//
//	r.GET("/assets/*filepath", func(c *gee.Context) {
//		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
//	})
//
//	r.Run(":9999")
//}

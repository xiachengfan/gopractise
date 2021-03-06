package main

import (
	"fmt"
	"github.com/xiachengfan/gopractise/xhttp/framework"
	"net/http"
)

func main() {
	r := framework.NewCore()
	r.Get("/index", func(c *framework.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	r.Get("/hello", func(c *framework.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	r.Get("/", func(c *framework.Context) {
		fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
	})
	v1 := r.Group("/v1")
	{
		v1.Get("/", func(c *framework.Context) {
			c.HTML(http.StatusOK, "<h1>Hello framework</h1>")
		})

		v1.Get("/hello", func(c *framework.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.Get("/login", func(c *framework.Context) {
			c.JSON(http.StatusOK, framework.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}

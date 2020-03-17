# ginx 
a handler wrapper for gin

### quick start
```go
func main() {
	app := gin.New()
	router := ginx.XRouter(app, nil)
	router.Any("/greeting", func(who *struct {
		Name string `json:"name" form:"name" binding:"required"`
	}) string {
		return fmt.Sprintf("Hi, %s!", who.Name)
	})

	log.Fatal(app.Run(":8090"))
}
```
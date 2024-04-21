package main

func main() {
	app := App{}
	app.Initialise(DBUser, DBPassword, DBName)
	app.Run("127.0.0.1:8080")
}

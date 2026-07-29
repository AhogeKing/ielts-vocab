package main

import (
	"ielts-vocab/internal/app"
	"ielts-vocab/internal/config"
	"ielts-vocab/internal/database"
	"log"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectToPostgres(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	application := app.New(db, appConfig)

	err = application.Router.Run(appConfig.Port)
	if err != nil {
		log.Fatal(err)
	}
}

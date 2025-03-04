package main

import (
	"diploma/db_configs"
	"diploma/services/admin/pkg/app"
	// "fmt"
)

func main() {
    config := new(db_configs.DBConfigs)
    config.GetAllDBsConfigs()

    // fmt.Println(config.Admin, config.Client)

    app.LaunchAdminService(config.Admin)
}
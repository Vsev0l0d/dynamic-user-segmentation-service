package main

import "dynamic-user-segmentation-service/internal/app"

const configPath = "internal/config/config.yaml"

func main() {
	app.Run(configPath)
}

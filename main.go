package main

import (
	globalcomps "eleliafrika.com/backend/global_comps"
)

func main() {
	globalcomps.LoadEnv()
	globalcomps.LoadDatabase()
	globalcomps.ServeApplication()
}

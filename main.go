package main

import (
	"fmt"
)

func main() {

	nDao, err := ConstructNotionDaoFromEnv()
	if err != nil {
		panic(fmt.Errorf("configuration error: %w", err))
	}

	tasks := GetAllTasks()
	errs := make([]error, len(tasks))
	for i, t := range tasks {
		errs[i] = t.Run(nDao)
	}

	PanicOnErrors(errs)
}

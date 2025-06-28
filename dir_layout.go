package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func mustCreate(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
}

func writeMain(path string) {
	mainPath := filepath.Join(path, "cmd", "server", "main.go")
	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello from your Go backend!")
}
`
	err := os.WriteFile(mainPath, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	project := "myproject"
	if len(os.Args) > 1 {
		project = os.Args[1]
	}

	paths := []string{
		"cmd/server",
		"internal/api",
		"internal/model",
		"internal/store",
		"pkg/config",
		"static",
		"templates",
		"configs",
		"scripts",
		"docs",
	}

	for _, p := range paths {
		mustCreate(filepath.Join(project, p))
	}

	writeMain(project)

	fmt.Printf("✅ Project structure for '%s' created!\n", project)
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
  "errors"
)

type Task struct {
	Name      string
	Completed bool
}

type Project struct {
	Name  string
	Tasks []Task
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getContext(name string) (context string) {
	file, err := os.ReadFile(name)
	check(err)
	context = string(file)
	context = strings.ReplaceAll(context, "\n", "")
	return context
}

func parseContext(context string) (projects []Project) {
	projectsSlice := strings.Split(context, "---")
	for _, project := range projectsSlice {
		projects = append(projects, parseProject(project))
	}

	return projects
}

func parseProject(projectContext string) (project Project) {
	// first implementation : will be refactored later
	projectExpression := strings.Split(projectContext, ",")
	project.Name = strings.Split(projectExpression[0], ":")[1]
	tasksSection := strings.Split(projectExpression[1], ":")[1]
	project.Tasks = parseTasks(tasksSection)
	return project
}

func parseTasks(tasksContext string) (tasks []Task) {
	tasksSlice := strings.Split(tasksContext, "*")
	for _, task := range tasksSlice {
		task = strings.TrimSpace(task)
		if len(task) > 0 {
			tasks = append(tasks, parseTask(task))
		}
	}
	return tasks
}

func parseTask(taskContext string) (task Task) {

	taskSlice := strings.Split(taskContext, "]")
	task.Name = taskSlice[1]
	if strings.Contains(taskSlice[0], "x") || strings.Contains(taskSlice[0], "X") {
		task.Completed = true
	}
	return task
}

func convertToJson(parsed []Project) (jsonContext string) {
	content, err := json.Marshal(parsed)
	check(err)
	jsonContext = string(content)
	return jsonContext
}

func outputFile(name string, content string) {
  fileName := strings.Split(name, ".")[0] + ".json"
	file, err := os.Create(fileName)
	check(err)
	file.WriteString(content)
}

func getFileName() (name string, e error) {
	if len(os.Args) > 1 {
		fileName := os.Args[1]
		fmt.Println(fileName)
	}else {
    return "",errors.New("Enter File Name")
  }
  return os.Args[1] , nil
}

func main() {

  fileName , err := getFileName()
  check(err)
	context := getContext(fileName)
	parsed := parseContext(context)
	jsonContent := convertToJson(parsed)
	outputFile(fileName, jsonContent)

}

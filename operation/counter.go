package operation

var apiCalls int = 0

var projects int = 0

func IncrementApiCalls() {
	apiCalls++
}

func GetApiCalls() int {
	return apiCalls
}

func IncrementProjects() {
	projects++
}

func GetProjects() int {
	return projects
}
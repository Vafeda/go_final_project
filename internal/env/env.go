package env

import "os"

var (
	Dict  map[string]string
	names = []string{
		"TODO_PORT",
		"TODO_DBFILE",
		"TODO_DBTESTDATA",
		"TODO_PASSWORD",
	}
)

func Init() {
	Dict = make(map[string]string)

	for _, name := range names {
		if envParam, exist := os.LookupEnv(name); exist {
			Dict[name] = envParam
		}
	}
}

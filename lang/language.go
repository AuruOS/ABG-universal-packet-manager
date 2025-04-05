package lang

import (
	"fmt"
	"log"
	"os"
)

func GetText(lang string, name string) string {
	content, err := os.ReadFile(fmt.Sprintf("lang/%s/%s", lang, name))
	if err != nil {
		log.Println("FATAL! Can't get translation file!")
		log.Fatal(err)
	}
	return string(content)
}

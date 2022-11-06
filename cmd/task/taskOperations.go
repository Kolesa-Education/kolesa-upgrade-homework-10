package task

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"log"
	"strings"
)

func checkTask(taskArgs []string) bool {
	if len(taskArgs) == 0 {
		log.Println("Нет аргументов!")
		return false
	}

	if len(taskArgs) < 3 {
		log.Println("Слишком мало аргументов")
		return false
	}

	regexpDate := pcre.MustCompile(`^(?:(?:31(\/|-|\.)(?:0?[13578]|1[02]))\1|(?:(?:29|30)(\/|-|\.)(?:0?[13-9]|1[0-2])\2))(?:(?:1[6-9]|[2-9]\d)?\d{2})$|^(?:29(\/|-|\.)0?2\3(?:(?:(?:1[6-9]|[2-9]\d)?(?:0[48]|[2468][048]|[13579][26])|(?:(?:16|[2468][048]|[3579][26])00))))$|^(?:0?[1-9]|1\d|2[0-8])(\/|-|\.)(?:(?:0?[1-9])|(?:1[0-2]))\4(?:(?:1[6-9]|[2-9]\d)?\d{2})$`, 0)
	match := regexpDate.MatcherString(strings.TrimSpace(taskArgs[2]), 0).Matches()

	if !match {
		log.Println("Не прошла проверка по регулярке")
		return false
	}

	return true
}

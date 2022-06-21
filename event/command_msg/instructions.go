package command_msg

import "fmt"

func CreateSetupMsg(username string) string {
	text := fmt.Sprintf("Приветствую тебя, *%s*"+"\n"+"\n"+
		"Основные команды, которые Вам помогут: "+"\n"+
		"*/help* - _помощь_"+"\n"+
		"*Отправить локацию* - _внизу экрана_"+"\n", username)

	return text
}

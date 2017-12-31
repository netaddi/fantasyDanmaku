package requestProcessors

import (
	"danmakuBackend/danmakuLib"
)

func sendAdminCommandToFrontend(operation string, parameter string){
	message := &danmakuLib.FrontendAdminMessage{
		"adminMessage",
		operation,
		parameter}
	Frontend.SendMessage(danmakuLib.GetJSON(message))
}

func ProcessAdminCommand (commandTokens []string) {
	switch commandTokens[1] {
	case "ban":
		switch commandTokens[2] {
		case "user":
			danmakuLib.SetPermission( commandTokens[3], 0 )
			break
		case "keyword":
			sendAdminCommandToFrontend("banKeyword", commandTokens[3])
			break
		}
		break
	case "open":
		switch commandTokens[2] {
		case "lottery":
			sendAdminCommandToFrontend("openLottery", "")
			break
		case "log" :
			sendAdminCommandToFrontend("openLog", "")
			break
		}
		break
	case "close":
		switch commandTokens[2] {
		case "lottery":
			sendAdminCommandToFrontend("closeLottery", "")
			break
		case "log":
			sendAdminCommandToFrontend("closeLog", "")
			break
		}
		break
	case "question":
		switch commandTokens[2] {
		case "prepare":
			prepareAnswering()
			break
		case "start":
			startAnswering()
			break
		case "end":
			endAnswering()
			break
		case "ranking":
			sendRanking()
			break
		}
	}


}
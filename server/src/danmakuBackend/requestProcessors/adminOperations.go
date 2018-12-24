package requestProcessors

import (
	"danmakuBackend/danmakuLib"
)

func sendAdminCommandToFrontend(operation string, parameter string){
	message := &danmakuLib.FrontendAdminMessage{
		MessageType:        "adminMessage",
		AdminOperation:     operation,
		OperationParameter: parameter,
	}
	Frontend.SendMessage(danmakuLib.GetJSON(message))
}

func ProcessAdminCommand (commandTokens []string) {
	switch commandTokens[1] {
	case "display":
		if commandTokens[2] == "goto" {
			sendAdminCommandToFrontend("goto", commandTokens[3])
		} else {
			sendAdminCommandToFrontend(commandTokens[2], "")
		}
	case "ban":
		if commandTokens[2] == "user" {
			danmakuLib.SetPermission( commandTokens[3], 0 )
		}
		if commandTokens[2] == "keyword" {
			sendAdminCommandToFrontend("banKeyword", commandTokens[3])
		}
		break
	case "unban":
		if commandTokens[2] == "user" {
			danmakuLib.SetPermission( commandTokens[3], 1 )
		}
		break
	case "open":
		switch commandTokens[2] {
		case "lottery":
			sendAdminCommandToFrontend("openLottery", "")
			break
		case "log":
			sendAdminCommandToFrontend("openLog", "")
			break
		case "ranking":
			sendAdminCommandToFrontend("openCommentRanking", "")
			break
		}
		break
	case "close":
		sendAdminCommandToFrontend("close", "")
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
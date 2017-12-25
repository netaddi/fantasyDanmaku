package requestProcessors

import (
	"danmakuBackend/danmakuLib"
)
//
//type FrontendAdminMessage struct {
//	isAdminCommand bool
//	adminOperation string
//	operationParameter string
//}

func sendAdminCommandToFrontend(operation string, parameter string){
	message := &danmakuLib.FrontendAdminMessage{true, operation, parameter}
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
	}
}
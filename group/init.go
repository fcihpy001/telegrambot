package group

import "context"

func Init() {
	LoadChatRules()

	InitCallbackRouters()
	InitLuckyFilter(context.Background())
}

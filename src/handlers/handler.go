package handlers

import (
	"log"
	"server/src/components"
)

func RequestService(ctx *components.Context) error {
	request, err := components.NewRequest(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	request.Send(ctx.Body)
	res, err := request.AwaitResponse()
	if err != nil {
		log.Println(err)
		return err
	}

	new_ctx := components.NewContext(res)
	ctx.Header = new_ctx.Header
	ctx.Body = new_ctx.Body

	return nil
}

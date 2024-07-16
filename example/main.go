package main

import (
	"fmt"
	"github.com/nikumar1206/puff"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Response struct {
	Error string `json:"error"`
}

func main() {
	app := puff.DefaultApp()
	app.Get("/", puff.Field{}, func(c *puff.Context) {
		c.SendResponse(puff.GenericResponse{
			StatusCode: 200,
			Content:    fmt.Sprintf("Hello there!"),
		})
	})
	app.WebSocket("/getUser", puff.Field{}, func(c *puff.Context) {
		c.WebSocket.OnMessage = func(ws *puff.WebSocket, msg puff.WebSocketMessage) {
			user := new(User)
			if err := msg.To(user); err != nil {
				ws.SendJSON(Response{
					Error: "Bad data payload.",
				})
				return
			}
			ws.SendJSON(Response{})
		}
	})
	app.Get("/hi", puff.Field{}, func(c *puff.Context) {
		c.SendResponse(puff.RedirectResponse{
			To: "https://youtube.com",
		})
	})
	app.ListenAndServe()
}

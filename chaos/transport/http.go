package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

type Handler interface {
	Name() string
	Execute(ctx context.Context, request *http.Request) response.Response
}

func Run(ip, port string) {
	go func() {
		err := http.ListenAndServe(ip+":"+port, nil)
		if err != nil {
			log.Fatalf("start chaos server failed, %v", err)
		}
	}()
}

func RegisterHandler(handler Handler) {
	http.HandleFunc(handler.Name(), func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			fmt.Fprintf(writer, response.ReturnIllegalParameters(err.Error()).Print())
			return
		}
		fmt.Fprintf(writer, handler.Execute(context.Background(), request).Print())
	})
}

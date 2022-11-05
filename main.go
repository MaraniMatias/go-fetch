package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"syscall/js"
)

func main() {
	js.Global().Set("GoFetch", GoFetch())

	select {} // block
}

func GoFetch() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolveFunction := args[0]
			//rejectFunction := args[1]

			go func() {
				res, err := http.DefaultClient.Get("https://catfact.ninja/fact")
				if err != nil {
					log.Fatal(err)
				}
				defer res.Body.Close()

				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
				}
				// Resolve the Promise
				resolveFunction.Invoke(string(b))
			}()

			return nil
		})

		// Create promise
		promise := js.Global().Get("Promise")
		return promise.New(handler)
	})
}

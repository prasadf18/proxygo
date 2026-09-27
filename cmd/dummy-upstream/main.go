package main
 
import (
	"fmt"
	"log"
	"net/http"
)
 
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from dummy upstream, you requested %s\n", r.URL.Path)
	})
 
	log.Println("dummy upstream listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
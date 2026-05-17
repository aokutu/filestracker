package modules

import  ( 
"net/http"
"log"

  )
func Readlogs(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")
	log.Println(ua)
	w.Write([]byte("hello",))

}
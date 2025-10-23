package middleware

import (
	"log"
	"net/http"
)

func Useless(next http.Handler) http.Handler{
	return  http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		log.Printf("useless middleware")
		
		next.ServeHTTP(w,r)
	})

}
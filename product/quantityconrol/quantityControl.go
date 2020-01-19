package main
import (
	"log"
	"net/http"
	"sync"
)

var (
	//TODO: get productNum from mysql
	productNum int64 = 999
	sum int64 = 0
	count int64 =0
	countMutex sync.Mutex
)

func GetOneProduct() bool {
	countMutex.Lock()
	count ++
	countMutex.Unlock()

	if count % 1 == 0 {
		if sum < productNum {
			sum ++
			return true
		}
	}
	return false
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	if GetOneProduct() {
		_, _ = w.Write([]byte("true"))
		return
	}
	_, _ = w.Write([]byte("false"))
	return
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	log.Fatal(http.ListenAndServe(":8084",nil))
}
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"product/common"
	"product/encrypt"
	"strconv"
	"sync"
)

var (
	//ip address of validate cluster,
	hostArray = []string{"127.0.0.1","127.0.0.1"}
	localHost = "127.0.0.1"
	//port of validate server
	port	  = "8083"
	hashConsistent *common.Consistent
)

type AccessControl struct {
	sourceArray map[int]interface{}
	sync.RWMutex
}
var accessControl = &AccessControl{sourceArray:make(map[int]interface{})}

func (m *AccessControl) GetNewRecord(uid int ) interface{} {
	m.RWMutex.RLock()
	defer  m.RWMutex.RUnlock()
	data := m.sourceArray[uid]
	return data
}

func (m *AccessControl) SetNewRecord(uid int){

}

func (m *AccessControl) GetDistributeRight(req *http.Request) bool {
	uid , err := req.Cookie("uid")
	if err != nil {
		return false
	}
	targetHost, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	if targetHost == localHost {
		return m.GetDataFromMap(uid.Value)
	}else {
		return GetDataFromOtherMap(targetHost,req)
	}
}

func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)

	if data != nil {
		return true
	}
	return false
}

//proxy
func GetDataFromOtherMap(host string, r *http.Request) bool {
	uid, err := r.Cookie("uid")
	if err != nil {
		return false
	}
	uidSign, err := r.Cookie("sign")
	if err != nil {
		return false
	}

	client := &http.Client{}
	url := "http://"+host+":"+port+"/check"
	req, err := http.NewRequest("GET",url,nil)
	if err != nil {
		return false
	}

	cookieUid := &http.Cookie{Name:"uid",Value:uid.Value,Path:"/"}
	cookieSign := &http.Cookie{Name:"sign",Value:uidSign.Value,Path:"/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	response , err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		}else {
			return false
		}
	}

	return false
}

func Check(w http.ResponseWriter, r *http.Request){
	//TODO check r
	fmt.Println("Checking...")
}

func Auth(w http.ResponseWriter,r *http.Request) error {
	//TODO: cookie authority
	fmt.Println("Cookie Auth")
	err := CheckUserInfoByCookie(r)
	if err != nil {
		return err
	}
	return nil

}
func CheckUserInfoByCookie(r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		//DO Noting
	}

	signCookie, err := r.Cookie("sign")
	if err != nil {
		//Do Noting
	}

	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		//
	}

	if checkInfoByCompare(uidCookie.Value, string(signByte)){
		// Pass
	}
	return nil
}
func checkInfoByCompare(checkStr,signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func main() {
	hashConsistent = common.NewConsistent()
	for _ ,host := range hostArray {
		hashConsistent.Add(host)
	}

	filter := common.NewFilter()
	filter.RegisterFilterUri("/check",Auth)
	http.HandleFunc("/check",filter.Handle(Check))
	//log.Fatal(http.ListenAndServe(":8083",nil))
	http.ListenAndServe(":8083",nil)
}
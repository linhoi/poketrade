package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"product/common"
	"product/datamodels"
	"product/encrypt"
	"product/rabbitmq"
	"strconv"
	"sync"
)

var (
	//ip listed here should be changed with real network env
	//ip address of validate cluster,
	hostArray = []string{"127.0.0.1","127.0.0.1"}
	localHost = "127.0.0.1"
	//port of validate server
	port	  = "8083"
	hashConsistent *common.Consistent

	QuantityControlServerIp = "127.0.0.1"
	QuantityControlServerPort = "8084"

	rabbitMqValidate  *rabbitmq.RabbitMQ
)

//store control message
type AccessControl struct {
	sourceArray map[int]interface{}
	sync.RWMutex
}

var accessControl = &AccessControl{sourceArray:make(map[int]interface{})}


//get data for specify uid
func (m *AccessControl) GetNewRecord(uid int ) interface{} {
	m.RWMutex.RLock()
	defer  m.RWMutex.RUnlock()
	data := m.sourceArray[uid]
	return data
}

func (m *AccessControl) SetNewRecord(uid int){
	m.RWMutex.Lock()
	m.sourceArray[uid] = "control message"
	m.RWMutex.Unlock()
}

func (m *AccessControl) GetDistributeRight(req *http.Request) bool {
	// get user uid
	uid , err := req.Cookie("uid")
	if err != nil {
		return false
	}

	//validate uid from validate server cluster
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
func GetDataFromOtherMap_Discard(host string, r *http.Request) bool {
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
func GetDataFromOtherMap(validateHost string, request *http.Request) bool {
	validateUrl := "http://"+validateHost +":"+port+"/checkRight"
	response, err := GetResponseFromProxy(validateUrl,request)
	if err != nil {
		return false
	}
	if response.StatusCode == 200 {
		bodyDate, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false
		}
		if string(bodyDate) == "true" {
			return true
		}else {
			return false
		}
	}
	return false
}
func GetResponseFromProxy(validateUrl string, request *http.Request) (response *http.Response, err error) {
	uid, err := request.Cookie("uid")
	if err != nil {
		return
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	client := &http.Client{}
	req , err := http.NewRequest("GET", validateUrl,nil)
	if err != nil {
		return
	}

	uidInCookie := &http.Cookie{
		Name:       "uid",
		Value:      uid.Value,
		Path:       "/",
	}
	signInCookie := &http.Cookie{
		Name:       "sign",
		Value:      uidSign.Value,
		Path:       "/",
	}

	request.AddCookie(uidInCookie)
	request.AddCookie(signInCookie)

	response ,err = client.Do(req)
	defer response.Body.Close()
	if err !=nil {
		return
	}
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	hasRight := accessControl.GetDistributeRight(r)
	if !hasRight {
		_, _ = w.Write([]byte("false"))
		return
	}
	_, _ = w.Write([]byte("true"))
	return
}




func Check(w http.ResponseWriter, r *http.Request){
	fmt.Println("Checking...")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productIDString := queryForm["productID"][0]
	//fmt.Println(productString)

	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	//1. distribute validate
	hasRight := accessControl.GetDistributeRight(r)
	if !hasRight {
		w.Write([]byte("false"))
		return
	}
	//2. quantity control
	quantityControlUrl := "http://"+QuantityControlServerIp+":"+QuantityControlServerPort+"/getOne"
	response , err := GetResponseFromProxy(quantityControlUrl, r)
	if response.StatusCode == 200 {
		validateBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.Write([]byte("false"))
			return
		}
		//1
		if string(validateBody) == "true" {
			productID, err := strconv.ParseInt(productIDString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			//2
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			//3.
			message := datamodels.NewMessage(userID, productID)
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			//4
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
		}
	}
	w.Write([]byte("true"))
	return
}







//validate request through cookie
func Auth(w http.ResponseWriter,r *http.Request) error {
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
		return errors.New("get uid in cookie failed")
	}

	signCookie, err := r.Cookie("sign")
	if err != nil {
		return  errors.New("get sign in cookie failed")
	}

	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("validate uid with sign failed")
	}

	if checkInfoByCompare(uidCookie.Value, string(signByte)){
		return nil
	}
	return errors.New("validate user failed")
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

	localIp ,err := common.GetLocalIP()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localHost)

	rabbitMqValidate := rabbitmq.NewRabbitMQSimple("product")
	defer rabbitMqValidate.Destroy()


	// filter
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check",Auth)
	filter.RegisterFilterUri("/checkRight",Auth)

	http.HandleFunc("/check",filter.Handle(Check))
	http.HandleFunc("/checkRight",filter.Handle(CheckRight))
	log.Fatal(http.ListenAndServe(":8083",nil))
	//http.ListenAndServe(":8083",nil)
}
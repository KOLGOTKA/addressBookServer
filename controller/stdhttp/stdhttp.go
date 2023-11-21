package stdhttp

import (
	"encoding/json"
	"fmt"
	"httpserver/gates/psg"
	"httpserver/models/dto"
	"httpserver/pkg"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Controller struct {
	srv http.Server
	db  *psg.Psg
}

func NewController(addr string, postgres *psg.Psg) (hs *Controller) {
	hs = new(Controller)
	hs.srv = http.Server{}
	// mux := http.NewServeMux()
	//
	// mux.HandleFunc("/get", http.HandlerFunc(hs.RecordGetHandler))
	// mux.Handle("/update", http.HandlerFunc(hs.RecordUpdateHandler))
	// mux.Handle("/delete", http.HandlerFunc(hs.RecordDeleteByPhoneHandler))

	mux := mux.NewRouter()
	mux.HandleFunc("/create", hs.RecordCreateHandler).Methods(http.MethodPost)
	mux.HandleFunc("/get", hs.RecordGetHandler).Methods(http.MethodPost)
	mux.HandleFunc("/update", hs.RecordUpdateHandler).Methods(http.MethodPut)
	mux.HandleFunc("/delete", hs.RecordDeleteByPhoneHandler).Methods(http.MethodDelete)
	hs.srv.Handler = mux
	hs.srv.Addr = addr
	hs.db = postgres
	return hs
}

func (hs *Controller) Start() *pkg.Errorer {
	errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) Start()"}
	err := hs.srv.ListenAndServe()
	if err != nil {
		errorer.Add("hs.srv.ListenAndServe(): " + err.Error())
		return errorer
	}
	return nil
}

// RecordCreate обрабатывает HTTP запрос для добавления новой записи.
func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request) {
	errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request)"}
	record, erro := GetBody(req)
	if erro != nil {
		errorer.Add("GetBody(req): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
	record.Phone, erro = pkg.PhoneNormalize(record.Phone)
	if erro != nil {
		errorer.Add("pkg.PhoneNormalize(record.Phone): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
	erro = hs.db.RecordCreate(*record)
	if erro != nil {
		errorer.Add("hs.db.RecordCreate(record): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
}

// RecordsGet обрабатывает HTTP запрос для получения записей на основе предоставленных полей Record.
func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request) {
	errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request)"}
	record, erro := GetBody(req)
	if erro != nil {
		errorer.Add("GetBody(req): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
	if record.Phone != "" {
		record.Phone, erro = pkg.PhoneNormalize(record.Phone)
		if erro != nil {
			errorer.Add("pkg.PhoneNormalize(record.Phone): " + erro.GetError())
			log.Println(errorer.GetError())
			return
		}
	}

	result, erro := hs.db.RecordsGet(*record)
	if erro != nil {
		errorer.Add("hs.db.RecordsGet(*record): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
	fmt.Println(result) /// По хорошему вернуть бы эти данные
}

// RecordUpdate обрабатывает HTTP запрос для обновления записи.
func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request) {
	errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request)"}
	record, erro := GetBody(req)
	if erro != nil {
		errorer.Add("GetBody(req): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
	record.Phone, erro = pkg.PhoneNormalize(record.Phone)
	if erro != nil {
		errorer.Add("pkg.PhoneNormalize(record.Phone): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
	erro = hs.db.RecordUpdate(*record)
	if erro != nil {
		errorer.Add("hs.db.RecordUpdate(*record): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
}

// RecordDeleteByPhone обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request) {
	errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request)"}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	// log.Printf("Request Body: %s", string(byteReq))
	if err != nil {
		errorer.Add("io.ReadAll(req.Body): " + err.Error())
		log.Println(errorer.GetError())
		return
	}
	phone, erro := pkg.PhoneNormalize(string(byteReq))
	if erro != nil {
		errorer.Add("pkg.PhoneNormalize(string(byteReq)): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
	erro = hs.db.RecordDeleteByPhone(phone)
	if erro != nil {
		errorer.Add("hs.db.RecordDeleteByPhone(phone): " + erro.GetError())
		log.Println(errorer.GetError())
		return
	}
}

func GetBody(req *http.Request) (record *dto.Record, errorer *pkg.Errorer) {
	errorer = &pkg.Errorer{Where: "stdhttp: func GetBody(req *http.Request)"}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	// log.Printf("Request Body: %s", string(byteReq))
	if err != nil {
		errorer.Add("io.ReadAll(req.Body): " + err.Error())
		log.Println(errorer.GetError())
		return nil, errorer
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		errorer.Add("json.Unmarshal(byteReq, &record): " + err.Error())
		log.Println(errorer.GetError())
		return nil, errorer
	}
	return record, nil
}

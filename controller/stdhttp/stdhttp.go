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
	// "github.com/gorilla/mux"
)

type Controller struct {
	srv http.Server
	db  *psg.Psg
}

func NewController(addr string, postgres *psg.Psg) (hs *Controller) {
	hs = new(Controller)
	hs.srv = http.Server{}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed) /// наверно также как обычные ошибки нужно обрабатывать
			return
		}
		hs.RecordCreateHandler(w, r)
	})
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.RecordGetHandler(w, r)
	})
	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.RecordUpdateHandler(w, r)
	})
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.RecordDeleteByPhoneHandler(w, r)
	})

	// mux := mux.NewRouter()
	// mux.HandleFunc("/create", hs.RecordCreateHandler).Methods(http.MethodPost)
	// mux.HandleFunc("/get", hs.RecordGetHandler).Methods(http.MethodPost)
	// mux.HandleFunc("/update", hs.RecordUpdateHandler).Methods(http.MethodPut)
	// mux.HandleFunc("/delete", hs.RecordDeleteByPhoneHandler).Methods(http.MethodDelete)
	hs.srv.Handler = mux
	hs.srv.Addr = addr
	hs.db = postgres
	return hs
}

func (hs *Controller) Start() error {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) Start()")
	// errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) Start()"}
	err := hs.srv.ListenAndServe()
	if err != nil {
		return myerr.Wrap(err, "hs.srv.ListenAndServe()")
		// errorer.Add("hs.srv.ListenAndServe(): " + err.Error())
		// return errorer
	}
	return nil
}

// RecordCreate обрабатывает HTTP запрос для добавления новой записи.
func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request) {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request)")
	// errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request)"}
	record, err := GetBody(req)
	if err != nil {
		e := myerr.Wrap(err, "GetBody(req)")
		// errorer.Add("GetBody(req): " + erro.GetError())
		log.Println(e.Error())
		return /////////////////////////////////////////////////////////////////////////////////////////////////////
	}
	/// Проверяем нет ли пустых значений /// нужно проверить
	if record.Address == "" || record.LastName == "" || record.Name == "" || record.Phone == "" {
		e := myerr.Wrap(nil, "All fields except for the middle name must be filled in")
		log.Println(e.Error())
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		e := myerr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
		log.Println(e.Error())
		return
		// errorer.Add("pkg.PhoneNormalize(record.Phone): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
	err = hs.db.RecordCreate(*record)
	if err != nil {
		e := myerr.Wrap(err, "hs.db.RecordCreate(record)")
		log.Println(e.Error())
		return
		// errorer.Add("hs.db.RecordCreate(record): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
}

// RecordsGet обрабатывает HTTP запрос для получения записей на основе предоставленных полей Record.
func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request) {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request)")
	// errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request)"}
	record, err := GetBody(req)
	if err != nil {
		e := myerr.Wrap(err, "GetBody(req)")
		log.Println(e.Error())
		return
		// errorer.Add("GetBody(req): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
	if record.Phone != "" {
		record.Phone, err = pkg.PhoneNormalize(record.Phone)
		if err != nil {
			e := myerr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
			log.Println(e.Error())
			return
			// errorer.Add("pkg.PhoneNormalize(record.Phone): " + erro.GetError())
			// log.Println(errorer.GetError())
			// return
		}
	}

	result, err := hs.db.RecordsGet(*record)
	if err != nil {
		e := myerr.Wrap(err, "hs.db.RecordsGet(*record)")
		log.Println(e.Error())
		return
		// errorer.Add("hs.db.RecordsGet(*record): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
	fmt.Println(result) /// По хорошему вернуть бы эти данные
	/// Формируем ответ
	w.Header().Set("Content-Type", "text/plain")
	jsonData, err := json.Marshal(result) /// возможно стоит поменять формат ответа
	if err != nil {
		e := myerr.Wrap(err, "json.Marshal(records)")
		log.Println(e.Error())
		return
	}
	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}

// RecordUpdate обрабатывает HTTP запрос для обновления записи.
func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request) {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request)")
	// errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request)"}
	record, err := GetBody(req)
	if err != nil {
		e := myerr.Wrap(err, "GetBody(req)")
		log.Println(e.Error())
		return
		// errorer.Add("GetBody(req): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		e := myerr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
		log.Println(e.Error())
		return
		// errorer.Add("pkg.PhoneNormalize(record.Phone): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
	err = hs.db.RecordUpdate(*record)
	if err != nil {
		e := myerr.Wrap(err, "hs.db.RecordUpdate(*record)")
		log.Println(e.Error())
		return
		// errorer.Add("hs.db.RecordUpdate(*record): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
}

// RecordDeleteByPhone обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request) {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request)")
	// errorer := &pkg.Errorer{Where: "stdhttp: func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request)"}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	// log.Printf("Request Body: %s", string(byteReq))
	if err != nil {
		e := myerr.Wrap(err, "io.ReadAll(req.Body)")
		log.Println(e.Error())
		return
		// errorer.Add("io.ReadAll(req.Body): " + err.Error())
		// log.Println(errorer.GetError())
		// return
	}
	phone, err := pkg.PhoneNormalize(string(byteReq))
	if err != nil {
		e := myerr.Wrap(err, "pkg.PhoneNormalize(string(byteReq))")
		log.Println(e.Error())
		return
		// errorer.Add("pkg.PhoneNormalize(string(byteReq)): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
	err = hs.db.RecordDeleteByPhone(phone)
	if err != nil {
		e := myerr.Wrap(err, "hs.db.RecordDeleteByPhone(phone)")
		log.Println(e.Error())
		return
		// errorer.Add("hs.db.RecordDeleteByPhone(phone): " + erro.GetError())
		// log.Println(errorer.GetError())
		// return
	}
}

func GetBody(req *http.Request) (record *dto.Record, e error) {
	myerr := pkg.NewMyError("package stdhttp: func GetBody(req *http.Request)")
	// errorer = &pkg.Errorer{Where: "stdhttp: func GetBody(req *http.Request)"}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	// log.Printf("Request Body: %s", string(byteReq))
	if err != nil {
		e = myerr.Wrap(err, "io.ReadAll(req.Body)")
		// errorer.Add("io.ReadAll(req.Body): " + err.Error())
		log.Println(e.Error())
		return nil, e
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		e = myerr.Wrap(err, "json.Unmarshal(byteReq, &record)")
		log.Println(e.Error())
		return nil, e
		// errorer.Add("json.Unmarshal(byteReq, &record): " + err.Error())
		// log.Println(errorer.GetError())

	}
	return record, nil
}

// result ok/err
// body
// error строка с ошибкой
// 200 ok in any cases

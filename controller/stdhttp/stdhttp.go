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
	hs.srv.Handler = mux
	hs.srv.Addr = addr
	hs.db = postgres
	return hs
}

func (hs *Controller) Start() error {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) Start()")
	err := hs.srv.ListenAndServe()
	if err != nil {
		return myerr.Wrap(err, "hs.srv.ListenAndServe()") /////////////////////////// тоже до запуска сервера
	}
	return nil
}

// RecordCreate обрабатывает HTTP запрос для добавления новой записи.
func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request) {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}
	record, err := GetBody(req)
	if err != nil {
		e := myerr.Wrap(err, "GetBody(req)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error()) ///////////////////////
		return                 /////////////////////////////////////////////////////////////////////////////////////////////////////
	}
	/// Проверяем нет ли пустых значений
	if record.Address == "" || record.MiddleName == "" || record.Name == "" || record.Phone == "" {
		e := myerr.Wrap(nil, "All fields except for the lastname must be filled in")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		e := myerr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	err = hs.db.RecordCreate(*record)
	if err != nil {
		e := myerr.Wrap(err, "hs.db.RecordCreate(record)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	response.Result = "Success"
	w.Write(response.GetJson())
}

// RecordsGet обрабатывает HTTP запрос для получения записей на основе предоставленных полей Record.
func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request) {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}
	record, err := GetBody(req)
	if err != nil {
		e := myerr.Wrap(err, "GetBody(req)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	if record.Phone != "" {
		record.Phone, err = pkg.PhoneNormalize(record.Phone)
		if err != nil {
			e := myerr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
			response.Result = "Error"
			response.Error = e.Error()
			w.Write(response.GetJson())
			log.Println(e.Error())
			return
		}
	}

	result, err := hs.db.RecordsGet(*record)
	if err != nil {
		e := myerr.Wrap(err, "hs.db.RecordsGet(*record)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	fmt.Println(result) /// По хорошему вернуть бы эти данные
	/// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(result) /// возможно стоит поменять формат ответа
	if err != nil {
		e := myerr.Wrap(err, "json.Marshal(records)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	response.Result = "Succes"
	response.Data = jsonData
	w.Write(response.GetJson())
}

// RecordUpdate обрабатывает HTTP запрос для обновления записи.
func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request) {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}
	record, err := GetBody(req)
	if err != nil {
		e := myerr.Wrap(err, "GetBody(req)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		e := myerr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	err = hs.db.RecordUpdate(*record)
	if err != nil {
		e := myerr.Wrap(err, "hs.db.RecordUpdate(*record)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	response.Result = "Succes"
	w.Write(response.GetJson())
}

/// RecordDeleteByPhone обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request) {
	myerr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) /// Status 200 OK
	response := dto.Response{}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e := myerr.Wrap(err, "io.ReadAll(req.Body)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	phone, err := pkg.PhoneNormalize(string(byteReq))
	if err != nil {
		e := myerr.Wrap(err, "pkg.PhoneNormalize(string(byteReq))")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	err = hs.db.RecordDeleteByPhone(phone)
	if err != nil {
		e := myerr.Wrap(err, "hs.db.RecordDeleteByPhone(phone)")
		response.Result = "Error"
		response.Error = e.Error()
		w.Write(response.GetJson())
		log.Println(e.Error())
		return
	}
	response.Result = "Success"
	w.Write(response.GetJson())
}

/// Функция для преобразования responce в структуру Record
func GetBody(req *http.Request) (record *dto.Record, e error) {
	myerr := pkg.NewMyError("package stdhttp: func GetBody(req *http.Request)")
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e = myerr.Wrap(err, "io.ReadAll(req.Body)")
		log.Println(e.Error())
		return nil, e
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		e = myerr.Wrap(err, "json.Unmarshal(byteReq, &record)")
		log.Println(e.Error())
		return nil, e
	}
	return record, nil
}

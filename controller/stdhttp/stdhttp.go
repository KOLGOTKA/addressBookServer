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
	"github.com/pkg/errors"
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

func (hs *Controller) Start() (err error) {
	err = hs.srv.ListenAndServe()
	if err != nil {
		err = errors.Wrap(err, "hs.Start(): hs.srv.ListenAndServe()")
		return err
	}
	return err
}

// RecordCreate обрабатывает HTTP запрос для добавления новой записи.
func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request) {
	record := dto.Record{}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	log.Printf("Request Body: %s", string(byteReq))
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): io.ReadAll(req.Body)")
		log.Println(err)
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): json.Unmarshal(byteReq, &record)"+string(byteReq))
		log.Println(err)
		return
	}
	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): pkg.PhoneNormalize(record.Phone)")
		log.Println(err)
		return
	}
	err = hs.db.RecordCreate(record)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): hs.db.RecordCreate(record)")
		log.Println(err)
		return
	}
}

// RecordsGet обрабатывает HTTP запрос для получения записей на основе предоставленных полей Record.
func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request) {
	record := dto.Record{}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	log.Printf("Request Body: %s", string(byteReq))
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): io.ReadAll(req.Body)")
		log.Println(err)
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): json.Unmarshal(byteReq, &record)")
		log.Println(err)
		return
	}
	if record.Phone != "" {
		record.Phone, err = pkg.PhoneNormalize(record.Phone)
		if err != nil {
			err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): pkg.PhoneNormalize(record.Phone)")
			log.Println(err)
			return
		}
	}

	result, err := hs.db.RecordsGet(record)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): hs.db.RecordsGet(record)")
		log.Println(err)
		return
	}
	fmt.Println(result) /// По хорошему вернуть бы эти данные
}

// RecordUpdate обрабатывает HTTP запрос для обновления записи.
func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request) {
	record := dto.Record{}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	log.Printf("Request Body: %s", string(byteReq))
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): io.ReadAll(req.Body)")
		log.Println(err)
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): json.Unmarshal(byteReq, &record)")
		log.Println(err)
		return
	}
	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): pkg.PhoneNormalize(record.Phone)")
		log.Println(err)
		return
	}
	err = hs.db.RecordUpdate(record)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): hs.db.RecordUpdate(record)")
		log.Println(err)
		return
	}
}

// RecordDeleteByPhone обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request) {
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	log.Printf("Request Body: %s", string(byteReq))
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): io.ReadAll(req.Body)")
		log.Println(err)
		return
	}
	phone, err := pkg.PhoneNormalize(string(byteReq))
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): pkg.PhoneNormalize(record.Phone)")
		log.Println(err)
		return
	}
	err = hs.db.RecordDeleteByPhone(phone)
	if err != nil {
		err = errors.Wrap(err, "Controller: (hs *controller) RecordCreateHandler(): hs.db.RecordDeleteByPhone(record)")
		log.Println(err)
		return
	}
}

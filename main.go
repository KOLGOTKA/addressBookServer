package main

import (
	// "fmt"

	"httpserver/controller/stdhttp"
	"httpserver/gates/psg"
	// "httpserver/models/dto"
	// "httpserver/pkg"
)

func main() {
	psgr := psg.NewPsg("localhost")
	serv := stdhttp.NewController(":8080", psgr)

	serv.Start()
	// fmt.Println("All ride")

	// rec := dto.Record{
	// 	Name:       "Иван",
	// 	LastName:   "Иванов",
	// 	MiddleName: "Иванович",
	// 	Address:    "Улица Пушкина, дом 1",
	// 	Phone:      "+1234567890",
	// }
	// rec1 := dto.Record{
	// 	Name:       "Петр",
	// 	LastName:   "Сергеевич",
	// 	MiddleName: "Кашпо",
	// 	Address:    "Улица Смородиновая, дом 23",
	// 	Phone:      "+79133840018",
	// }
	// rec2 := dto.Record{
	// 	Name:       "Иван",
	// 	LastName:   "Романов",
	// 	MiddleName: "Георгиевич",
	// 	Address:    "Улица Пушкина, дом 1",
	// 	Phone:      "+1234567890",
	// }
	// search_rec := dto.Record{
	// 	Name:       "Иван",
	// 	LastName:   "",
	// 	MiddleName: "",
	// 	Address:    "",
	// 	Phone:      "",
	// }
	// _, err := psgr.RecordCreate(rec)
	// _, err = psgr.RecordCreate(rec1)
	// _, err = psgr.RecordCreate(rec2)
	// // fmt.Println("Created", err)
	// query, err := psgr.RecordsGet(search_rec)
	// fmt.Println("Searched", err, query)

	// // print(psgr)
	// // err := psgr.RecordUpdate(rec2)
	// // fmt.Println("Updated", err)

	// err = psgr.RecordDeleteByPhone("+79133840018")
	// fmt.Println("Deleted", err)

	// // psgr.conn.Close()
	// // p := psg.NewPsg()
	// // hs := httpserver.NewHttpServer(":8080")
	// // hs.Start()
	// // fmt.Println(pkg.PhoneNormalize("8(985)-46-52-19"))
}

package psg

import (
	"context"
	"fmt"
	"httpserver/models/dto"
	"httpserver/pkg"
	"log"
	"net/url"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Psg struct {
	conn *pgxpool.Pool
}

// NewPsg создает новый экземпляр Psg.
func NewPsg(psgAddr string, dbUser string, dbPassword string) *Psg { //////////////// Возвращать ошибку или нет
	psg := &Psg{}

	if dbPassword == "" || dbUser == "" || psgAddr == "" {
		fmt.Println("Error: please write correct arguments") ////////////////////////// 1
		return nil
	}
	db_url := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   psgAddr + ":5432",
		Path:   "records",
	}
	connPool, err := pgxpool.New(context.Background(), db_url.String())

	/// Так?
	if err != nil {
		log.Fatal(err, "Something going wrong with connecting to database") /////////////////////// 2
	}
	psg.conn = connPool
	return psg
}

// RecordCreate добавляет новую запись в базу данных.
func (p *Psg) RecordCreate(record dto.Record) error {
	myerr := pkg.NewMyError("package pkg: func (p *Psg) RecordCreate(record dto.Record) error")
	// err := pkg.NewMyError("package pkg: func (p *Psg) RecordCreate(record dto.Record) error")
	// err := &pkg.Errorer{Where: "psg: func (p *Psg) RecordCreate(record dto.Record)"}
	record_existence, err := p.CheckPhone(record.Phone)
	if err != nil {
		return err //////////////////////////////////////////////////////////////////
		// errorer.Add(erro.GetError())
		// return errorer
	}
	if record_existence {
		return myerr.Wrap(nil, "Phone already exist")
		// errorer.Add("Phone already exist")
		// return errorer
	}
	_, err = p.conn.Exec(context.Background(), "INSERT INTO records (name, last_name, middle_name, address, phone) VALUES ($1, $2, $3, $4, $5)",
		record.Name, record.LastName, record.MiddleName, record.Address, record.Phone)

	if err != nil {
		return myerr.Wrap(err, "Something going wrong with inserting into db")
		// errorer.Add("p.conn.Exec(context.Background(), 'INSERT INTO records (name, last_name, middle_name, address, phone) VALUES ($1, $2, $3, $4, $5)', record.Name, record.LastName, record.MiddleName, record.Address, record.Phone): " + "Something going wrong with inserting into table 'records': " + err.Error())
		// return errorer
	}

	return nil
}

// RecordsGet возвращает записи из базы данных на основе предоставленных полей Record.
func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, error) {
	myerr := pkg.NewMyError("package psg: func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, *pkg.Errorer)")
	// errorer := &pkg.Errorer{Where: "psg: func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, *pkg.Errorer)"}
	var result []dto.Record

	var values []interface{}
	query := "SELECT * FROM records WHERE TRUE"
	// Будем добавлять параметры в запрос только если они не пустые
	if record.Name != "" {
		query += " AND name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Name)
	}
	if record.LastName != "" {
		query += " AND last_name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.LastName)
	}
	if record.MiddleName != "" {
		query += " AND middle_name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.MiddleName)
	}
	if record.Address != "" {
		query += " AND address=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Address)
	}
	if record.Phone != "" {
		query += " AND phone=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Phone)
	}
	/// хорошим вариантом была бы text template lib

	// Выполняем запрос с параметрами
	rows, err := p.conn.Query(context.Background(), query, values...)

	fmt.Println(query, values) ///////////////////////////////////////

	if err != nil {
		return nil, myerr.Wrap(err, "p.conn.Query(context.Background(), query, values...): Something going wrong with query to database")
		// errorer.Add("p.conn.Query(context.Background(), query, values...): Something going wrong with query to database")
		// return nil, errorer
	}
	defer rows.Close()
	for rows.Next() {
		var rec dto.Record
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone); err != nil {
			return nil, myerr.Wrap(err, "rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone): Something going wrong scan")
			// errorer.Add("rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone): Something going wrong scan")
			// return nil, errorer
		}
		result = append(result, rec)
		fmt.Printf("ID: %d, Name: %s, Last Name: %s, Middle Name: %s, Address: %s, Phone: %s\n", rec.ID, rec.Name, rec.LastName, rec.MiddleName, rec.Address, rec.Phone)
	}
	return result, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(record dto.Record) error {
	myerr := pkg.NewMyError("package psg: func (p *Psg) RecordUpdate(record dto.Record)")
	// errorer := &pkg.Errorer{Where: "psg: func (p *Psg) RecordUpdate(record dto.Record)"}
	/// должны обновляться не все поля, а только те, которые передаём
	var values []interface{}
	query := "UPDATE records SET"
	// Будем добавлять параметры в запрос только если они не пустые
	if record.Name != "" {
		query += " name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Name)
	}
	if record.LastName != "" {
		if len(values) > 0 {
			query += ","
		}
		query += " last_name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.LastName)
	}
	if record.MiddleName != "" {
		if len(values) > 0 {
			query += ","
		}
		query += " middle_name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.MiddleName)
	}
	if record.Address != "" {
		if len(values) > 0 {
			query += ","
		}
		query += " address=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Address)
	}
	if record.Phone == "" {
		return myerr.Wrap(nil, "Empty phone number")
		// errorer.Add("Empty phone number")
		// return errorer
	}
	if len(values) == 0 {
		return myerr.Wrap(nil, "Nothing to update")
		// errorer.Add("Nothing to update")
		// return errorer
	}
	query += " WHERE phone=$" + strconv.Itoa(len(values)+1)
	values = append(values, record.Phone)
	// query := "UPDATE records SET name=$1, last_name=$2, middle_name=$3, address=$4 WHERE phone=$5"
	// _, err := p.conn.Exec(context.Background(), query, record.Name, record.LastName, record.MiddleName, record.Address, record.Phone)
	_, err := p.conn.Exec(context.Background(), query, values...)
	if err != nil {
		return myerr.Wrap(err, "p.conn.Exec(context.Background(), query, values...): Something going wrong with query to database")
		// errorer.Add("p.conn.Exec(context.Background(), query, values...): Something going wrong with query to database")
		// return errorer
	}
	return nil
}

// RecordDeleteByPhone удаляет запись из базы данных по номеру телефона.
func (p *Psg) RecordDeleteByPhone(phone string) error {
	myerr := pkg.NewMyError("package psg: func (p *Psg) RecordDeleteByPhone(phone string)")
	// errorer := &pkg.Errorer{Where: "psg: func (p *Psg) RecordDeleteByPhone(phone string)"}
	query := "DELETE FROM records WHERE phone=$1"
	_, err := p.conn.Exec(context.Background(), query, phone)
	if err != nil {
		return myerr.Wrap(err, "p.conn.Exec(context.Background(), query, phone): Something going wrong with delete from database")
		// errorer.Add("p.conn.Exec(context.Background(), query, phone): Something going wrong with delete from database")
		// return errorer
	}
	return nil
}

// Функция закрытия соединения с БД
func (p *Psg) Close() {
	p.conn.Close()
}

// функция для проверки наличия номера телефона в БД
func (p *Psg) CheckPhone(phone string) (bool, error) {
	myerr := pkg.NewMyError("func (p *Psg) CheckPhone(phone string) (bool, *pkg.Errorer)")
	// errorer := &pkg.Errorer{Where: "psg: func (p *Psg) CheckPhone(phone string) (bool, *pkg.Errorer)"}
	query := "SELECT EXISTS (SELECT * FROM records WHERE phone = $1)"
	var exists bool
	err := p.conn.QueryRow(context.Background(), query, phone).Scan(&exists)
	if err != nil {
		return false, myerr.Wrap(err, "p.conn.QueryRow(context.Background(), query, phone).Scan(&exists)")
		// errorer.Add("p.conn.QueryRow(context.Background(), query, phone).Scan(&exists): " + err.Error())
		// return false, errorer
	}
	return exists, nil
}

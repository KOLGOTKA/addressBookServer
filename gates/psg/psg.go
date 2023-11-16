package psg

import (
	"context"
	"fmt"
	"httpserver/models/dto"
	"log"
	"net/url"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Psg struct {
	conn *pgxpool.Pool
}

// NewPsg создает новый экземпляр Psg.
func NewPsg(psgAddr string, dbPassword string) *Psg {
	psg := &Psg{}

	if dbPassword == "" {
		fmt.Println("Error: no password")
		return nil
	}
	db_url := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("postgres", dbPassword),
		Host:   psgAddr + ":5432",
		Path:   "records",
	}
	connPool, err := pgxpool.New(context.Background(), db_url.String())

	/// Так?

	/// "postgresql://postgres:"+dbPassword+"@"+psgAddr+":5432/records")
	/// через библиотеку url
	if err != nil {
		log.Fatal(err, "Something going wrong with connecting to database")
	}
	psg.conn = connPool
	// createTableSQL := `
	// CREATE TABLE IF NOT EXISTS records (
	// 	id SERIAL PRIMARY KEY,
	// 	name VARCHAR(30),
	// 	last_name VARCHAR(30),
	// 	middle_name VARCHAR(30),
	// 	address VARCHAR(255),
	// 	phone VARCHAR(20)
	// )
	// ` ///такого создания таблиц не должно быть. Сделать через библиотеку
	// _, err = psg.conn.Exec(context.Background(), createTableSQL)
	/// А что не так? Как иначе. Все под капотом используют psg.conn.Exec.

	// if err != nil {
	// 	log.Fatal(err, "Something going wrong with table creating")
	// }
	// defer psg.conn.Close() Не забыть закрыть соединение
	return psg
}

// RecordCreate добавляет новую запись в базу данных.
func (p *Psg) RecordCreate(record dto.Record) error { /// убрать int64 в возвращанемом
	_, err := p.conn.Exec(context.Background(), "INSERT INTO records (name, last_name, middle_name, address, phone) VALUES ($1, $2, $3, $4, $5)",
		record.Name, record.LastName, record.MiddleName, record.Address, record.Phone)

	if err != nil {
		return err
	}

	return nil
}

// RecordsGet возвращает записи из базы данных на основе предоставленных полей Record.
func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, error) {
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
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	fmt.Println("there") /////////////////////////
	for rows.Next() {
		var rec dto.Record
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone); err != nil {
			log.Println(err)
			return nil, err
		}
		result = append(result, rec)
		fmt.Printf("ID: %d, Name: %s, Last Name: %s, Middle Name: %s, Address: %s, Phone: %s\n", rec.ID, rec.Name, rec.LastName, rec.MiddleName, rec.Address, rec.Phone)
	}
	return result, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(record dto.Record) error {
	/// должны обновляться не все поля, а только те, которые передаём
	query := "UPDATE records SET name=$1, last_name=$2, middle_name=$3, address=$4 WHERE phone=$5"
	_, err := p.conn.Exec(context.Background(), query, record.Name, record.LastName, record.MiddleName, record.Address, record.Phone)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// RecordDeleteByPhone удаляет запись из базы данных по номеру телефона.
func (p *Psg) RecordDeleteByPhone(phone string) error {
	query := "DELETE FROM records WHERE phone=$1"
	_, err := p.conn.Exec(context.Background(), query, phone)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// / Функция закрытия соединения с БД
func (p *Psg) Close() {
	p.conn.Close()
}

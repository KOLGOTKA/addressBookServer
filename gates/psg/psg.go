package psg

import (
	"context"
	"fmt"
	"httpserver/models/dto"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Psg struct {
	conn *pgxpool.Pool
}

// NewPsg создает новый экземпляр Psg.
func NewPsg(psgAddr string) *Psg {
	psg := &Psg{}
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPassword == "" {
		fmt.Println("Error: no password")
		return nil
	}
	connPool, err := pgxpool.New(context.Background(), "postgresql://postgres:"+dbPassword+"@"+psgAddr+":5432/records")
	if err != nil {
		log.Fatal(err, "Something going wrong with connecting to database")
	}
	psg.conn = connPool
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS records (
		id SERIAL PRIMARY KEY,
		name VARCHAR(30),
		last_name VARCHAR(30),
		middle_name VARCHAR(30),
		address VARCHAR(255),
		phone VARCHAR(20)
	)
	`
	_, err = psg.conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatal(err, "Something going wrong with table creating")
	}
	// defer psg.conn.Close() Не забыть закрыть соединение
	return psg
}

// RecordCreate добавляет новую запись в базу данных.
func (p *Psg) RecordCreate(record dto.Record) (int64, error) {
	_, err := p.conn.Exec(context.Background(), "INSERT INTO records (name, last_name, middle_name, address, phone) VALUES ($1, $2, $3, $4, $5)",
		record.Name, record.LastName, record.MiddleName, record.Address, record.Phone)

	if err != nil {
		return -1, err
	}

	return 0, nil
}

// RecordsGet возвращает записи из базы данных на основе предоставленных полей Record.
func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, error) {
	var result []dto.Record
	name := record.Name
	last_name := record.LastName
	middle_name := record.MiddleName
	address := record.Address
	phone := record.Phone
	query := "SELECT * FROM records"
	// Будем заполнять строку с фильтрами в зависимоти от непустых полей
	search_params := ""
	if name != "" {
		search_params = search_params + " AND name='" + name + "'"
	}
	if last_name != "" {
		search_params = search_params + " AND last_name='" + last_name + "'"
	}
	if middle_name != "" {
		search_params = search_params + " AND middle_name='" + middle_name + "'"
	}
	if address != "" {
		search_params = search_params + " AND address='" + address + "'"
	}
	if phone != "" {
		search_params = search_params + " AND phone='" + phone + "'"
	}
	// Теперь уберём первый союз AND и добави ключивое слово WHERE
	if search_params != "" {
		search_params = " WHERE" + search_params[4:]
	}
	// fmt.Println(query + search_params)
	rows, err := p.conn.Query(context.Background(), query+search_params)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec dto.Record
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone); err != nil {
			log.Fatal(err)
			return nil, err
		}
		result = append(result, rec)
		fmt.Printf("ID: %d, Name: %s, Last Name: %s, Middle Name: %s, Address: %s, Phone: %s\n", rec.ID, rec.Name, rec.LastName, rec.MiddleName, rec.Address, rec.Phone)
	}
	return result, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(record dto.Record) error {
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

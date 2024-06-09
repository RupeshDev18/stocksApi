package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rupeshdev18/stocksApi/models"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to postgres")
	return db
}
func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("unable to decode the request body. %v", err)

	}
	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "stock created successfully",
	}
	json.NewEncoder(w).Encode(res)
}
func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int")
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatal("unable to get stock. %v", err)
	}
	json.NewEncoder(w).Encode(stock)
}
func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()
	if err != nil {
		log.Fatalf("Unable to get all the stocks %v", err)
	}
	json.NewEncoder(w).Encode(stocks)

}
func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int")
	}
	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("unable to decode the req body. %v", err)
	}
	updatedStock := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Stock updated successfully,Total rows/records affected %v", updatedStock)

	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int")
	}
	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("Stocks deleted successfully. Total rows/records %v", deletedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO stocks(name,price,company) values ($1,$2.$3) RETURNING stockid`
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}
	fmt.Println("Inserted a single record %v", id)
	return id
}
func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()
	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No Rows were Returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatal("Unable to scan the rows. %v", err)
	}
	return stock, err
}

func getAllStocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()
	var stocks []models.Stock

	sqlStatement := `SELECT * FROM stocks`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatal("unable to execute the query. %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Unable to scan the row %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `Update stocks SET name=$2, price=$3, company=$4 where stockid=$1`
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Unable to exec the query %v", err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Println("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `DELETE from stocks Where stockid=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to exec the query %v", err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Println("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}

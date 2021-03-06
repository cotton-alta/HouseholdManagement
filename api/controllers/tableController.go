package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"gopkg.in/gorp.v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/joho/godotenv"
)

type FormData struct {
	Body   string `json:"title"` 
	Amount int64 `json:"amount"`
	Genre  int64 `json:"genre"`
}


type HouseHold struct {
	Id      int64      `db:"id, primarykey"`
	Created time.Time  `db:"created_at"`
	Updated time.Time  `db:"updated_at"`
	Amount  int64      `db:"amount_money"`
	Body    string     `db:"remark,size:1024"`
	Balance int64      `db:"balance"` 
	Genre   int64      `db:"genre"`  
}

type JoinTable struct {
	Id      int64      `db:"id, primarykey"`
	Created time.Time  `db:"created_at"`
	Updated time.Time  `db:"updated_at"`
	Amount  int64      `db:"amount_money"`
	Body    string     `db:"remark,size:1024"`
	Balance int64      `db:"balance"` 
	Genre   string      `db:"genre"` 
	GenreId int64      `db:"genre_id"` 
}

func GetTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		dbmap := initDb()
		defer dbmap.Db.Close()
		
		var management  []JoinTable
		//SQLのクエリを直接記入（今回はselect文）
		_, err := dbmap.Select(
			&management,
			"select * from main_list join genre_list on main_list.Genre = genre_list.genre_id order by Id desc")
		checkErr(err, "select failed")
		fmt.Printf("%#v", management)
		
		return c.JSON(http.StatusOK, management)
	}
}

func GetItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		dbmap := initDb()
		defer dbmap.Db.Close()
		
		pathName := c.Param("item")
		// fmt.Printf("%v", pathName)
		var item HouseHold
		err := dbmap.SelectOne(
			&item,
			"select * from main_list where id=? order by Id desc limit 1 ",
			pathName,
		)
		checkErr(err, "select failed")
		// fmt.Printf("%#v", item)
		
		return c.JSON(http.StatusOK, item)
	}
}

func UpdateItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		dbmap := initDb()
		defer dbmap.Db.Close()

		var data FormData
		err := c.Bind(&data)
		checkErr(err, "data catch error")
		fmt.Printf("%#v\n", data)
		
		pathName := c.Param("item")
		var item HouseHold
		err = dbmap.SelectOne(
			&item,
			"select * from main_list where id=? order by Id desc limit 1 ",
			pathName,
		)
		checkErr(err, "select failed")
		
		dateNow := time.Now()
		difference := data.Amount - item.Amount
		updateBalance := item.Balance - difference
		
		fmt.Printf("%v\n", updateBalance)

		updateData := &HouseHold{
			item.Id,
			item.Created,
			dateNow,
			data.Amount,
			data.Body,
			updateBalance,
			data.Genre}
		fmt.Printf("%#v\n", updateData)
		_, err = dbmap.Update(updateData)
		checkErr(err, "insert failed")

		_, err = dbmap.Exec(
			"update main_list set balance = balance - ? where id > ?",
			difference,
			pathName,
		)

		return c.JSON(http.StatusOK, "updeted item!")
	}
}

func DeleteItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		dbmap := initDb()
		defer dbmap.Db.Close()

		pathName := c.Param("item")
		fmt.Printf("%v\n", pathName)
		var item HouseHold
		err := dbmap.SelectOne(
			&item,
			"select * from main_list where id=? order by Id desc limit 1 ",
			pathName,
		)

		amount := item.Amount
		dbmap.Delete(&item)
		
		fmt.Printf("%#v\n", item)

		_, err = dbmap.Exec(
			"update main_list set balance = balance + ? where id > ?",
			amount,
			pathName,
		)

		_, err = dbmap.Exec("update main_list set id = id - 1 where id > ?", pathName)
		// (select id from main_list where id = ?)
		checkErr(err, "select failed")
		return c.String(http.StatusOK, "item deleted")
	}
}

func CreateItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		dbmap := initDb()
		defer dbmap.Db.Close()
		
		var data FormData
		err := c.Bind(&data)
		checkErr(err, "data catch error")
		fmt.Printf("%#v", data)

		var beforeItem HouseHold
		err = dbmap.SelectOne(&beforeItem,
			"select * from main_list order by Id desc limit 1")
		checkErr(err, "select failed")

		fmt.Printf("%#v\n", beforeItem)

		dateNow := time.Now()
		item := &HouseHold{
			beforeItem.Id + 1,
			dateNow,
			dateNow,
			data.Amount,
			data.Body,
			beforeItem.Balance - data.Amount,
			data.Genre}
		err = dbmap.Insert(item)
		checkErr(err, "insert failed")

		return c.String(http.StatusOK, "created item!")
	}
}

//dbに関する設定、今後modelsに移行予定
func initDb() *gorp.DbMap {
	err := godotenv.Load()
	checkErr(err, "env failed")
	db, err := sql.Open("mysql", os.Getenv("TABLENAME"))
	checkErr(err, "sql.Open failed")
	
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF-8"}}

	// err = dbmap.CreateTablesIfNotExists()
	// checkErr(err, "Create tables failed")
	
	table := dbmap.AddTableWithName(HouseHold{}, "main_list")
	table.ColMap("Id").Rename("id")
	table.ColMap("Created").Rename("created_at")
	table.ColMap("Updated").Rename("updated_at")
	table.ColMap("Amount").Rename("amount_money")
	table.ColMap("Body").Rename("remark")
	table.ColMap("Balance").Rename("balance")
	table.ColMap("Genre").Rename("genre")

	return dbmap
}

//エラーハンドリング用関数
func checkErr(err error, msg string) {
	if err != nil {
			log.Fatalln(msg, err)
	}
}
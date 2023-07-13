package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"net/http"
)
// 本の構造体
type Book struct {
	gorm.Model
	Title string `gorm:"type:varchar(100)" json:"title"`
}
// DBの変数
var DB *gorm.DB

func main() {
	// errは、エラーがあれば返す
	var err error
	// DB, err = gorm.Openを使って、DBに接続。ユーザー名/パスワード/DB名を指定
	DB, err = gorm.Open("mysql", "root:1234@/MyData?charset=utf8&parseTime=True&loc=Local")
	// エラーがあれば、panicで停止
	if err != nil {
		panic("failed to connect database")
	}
	// DBを閉じる
	defer DB.Close()
  // DBのログを表示
	DB.AutoMigrate(&Book{})
  // gin.Default()でginのルーターを作成
	r := gin.Default()
  // r.GETで/booksにアクセスしたときに、getBooks関数を実行
	r.GET("/books", getBooks)
	// r.POSTで/booksにアクセスしたときに、createBook関数を実行
	r.POST("/books", createBook)
  // r.Runでサーバーを起動
	r.Run()
}
// MySQLのbooksテーブルから、全ての本を取得して、JSONで返す
func getBooks(c *gin.Context) {
	var books []Book
	if err := DB.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting books"})
		return
	}
	c.JSON(http.StatusOK, books)
}
// リクエストのJSONをBook構造体に変換して、MySQLのbooksテーブルに保存
func createBook(c *gin.Context) {
	var book Book
	// BindJSONで、リクエストのJSONをBook構造体に変換
	if err := c.BindJSON(&book); err != nil {
		// エラーがあれば、エラーを返す
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// DBに保存
	if err := DB.Save(&book).Error; err != nil {
		// エラーがあれば、エラーを返す
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating book"})
		return
	}
	// 保存した本を返す
	c.JSON(http.StatusOK, book)
}

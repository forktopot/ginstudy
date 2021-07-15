package common

//数据库

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Collection

func InitDB() *mongo.Client {
	// driverName := viper.GetString("datasource.driverName")
	// host := viper.GetString("datasource.host")
	// port := viper.GetString("datasource.port")
	// database := viper.GetString("datasource.database")
	// username := viper.GetString("datasource.username")
	// password := viper.GetString("datasource.password")
	// charset := viper.GetString("datasource.charset")
	// args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", username, password, host, port, database, charset)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")) //连接mongodb,连接方式mongodb://username:password@192.168.1.237:27017/ichunt?authMechanism=SCRAM-SHA-1

	err = client.Ping(ctx, readpref.Primary()) //如果不为空，连接失败
	if err != nil {
		log.Fatal(err)
	}
	DB = client.Database("test").Collection("trainers")

	return client

}

func GetDB() *mongo.Collection {
	return DB
}

func CloseDB(db *mongo.Client) {
	err := db.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
}

// func FindOne(ctx context.Context, filter interface{}, user *model.User, opts ...*options.FindOneOptions) (interface{}, error) {
// 	DB := GetDB()
// 	err := DB.FindOne(context.Background(), filter).Decode(&user)
// 	return user, err
// }

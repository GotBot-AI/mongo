package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//HOW TO USE/////////////////////////////////////////////////

// //Create monogo client object
// client, err := mongogo.Create(mongogo.ClientOptions{URI: "mongodb://localhost:27017", DBName: "my_db"})
// if err != nil {
// 	log.Fatal(err)
// }

// //Save object
// res := client.Save("cars", Car{Name: "Porche"})
// fmt.Println(res)

// //Find One
// var car Car
// client.FindOne("cars", mongogo.BSOND{{"name", "Volkswagen"}}).Decode(&car)
// fmt.Println(car.Name)

// //Find many
// cars := client.Find("cars", mongogo.BSOND{{}})
// for _, car := range cars {
// 	var c Car
// 	err := car.Decode(&c)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// //Find one and update
// doc, _ := client.ToBsonDoc(Car{Name: "Toyota"})
// var c Car
// client.FindOneAndUpdate("cars", mongogo.BSOND{{"name", "Porche"}}, doc).Decode(&c)
// fmt.Println(c)
// var newCar Car
// client.FindOne("cars", mongogo.BSOND{{"_id", c.ID}}).Decode(&newCar)
// fmt.Println(newCar)

// //Delete Many or 1
// deleteTotal := client.DeleteMany("cars", mongogo.BSOND{{"name", "Porche"}})
// fmt.Println(deleteTotal)

//ObjectID - primitive reference
type ObjectID = primitive.ObjectID

//BSOND - bson.D reference
type BSOND = bson.D

//ClientOptions - connection settings and options for mongoDB
type ClientOptions struct {
	URI    string
	DBName string
}

//Client - client object
type Client struct {
	DB *mongo.Database
}

//Create - Create mongo client using the client options
func Create(co ClientOptions) (Client, error) {
	clientOptions := options.Client().ApplyURI(co.URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
		return Client{}, err
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return Client{}, err
	}
	fmt.Println("Connected to MongoDB!")
	db := client.Database(co.DBName)
	cl := Client{DB: db}
	return cl, nil
}

//Save - Save a document to mongoDB
func (clnt *Client) Save(coll string, model interface{}) *mongo.InsertOneResult {
	collection := clnt.DB.Collection(coll)
	insertResult, err := collection.InsertOne(context.TODO(), &model)
	if err != nil {
		log.Fatal(err)
	}
	return insertResult
}

//DeleteMany - Delete one or many objects
func (clnt *Client) DeleteMany(coll string, filter bson.D) int64 {
	collection := clnt.DB.Collection(coll)
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return deleteResult.DeletedCount
}

//FindOneAndUpdate - Update one document
func (clnt *Client) FindOneAndUpdate(coll string, filter bson.D, query *primitive.D) *mongo.SingleResult {
	collection := clnt.DB.Collection(coll)
	singleResult := collection.FindOneAndUpdate(context.TODO(), filter, bson.D{{"$set", query}})
	return singleResult
}

//FindOne - Find one document
func (clnt *Client) FindOne(coll string, query interface{}) *mongo.SingleResult {
	collection := clnt.DB.Collection(coll)
	singleResult := collection.FindOne(context.TODO(), query)
	return singleResult
}

//Find - find many documents
func (clnt *Client) Find(coll string, query interface{}) []*mongo.Cursor {
	collection := clnt.DB.Collection(coll)
	findOptions := options.Find()
	var results []*mongo.Cursor
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		results = append(results, cur)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return results
}

//ToBsonDoc - Helper fuction that converts struct to byson for filters and queries
func (clnt *Client) ToBsonDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type UserInfo struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Pass     []byte             `bson:"pass" json:"pass"`
	Words    []WordInfo         `bson:"words" json:"words"`
	UuidHash [16]byte           `bson:"uuidHash" json:"uuidHash"`
}

type UserInfoset struct {
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
}
type WordInfo struct {
	Word  string    `bson:"word" json:"word"`
	Tword string    `bson:"Tword" json:"Tword"`
	Info  string    `bson:"Info" json:"Info"`
	Ldate time.Time `bson:"Ldate" json:"Ldate"`
	Learn int       `bson:"learn" json:"learn"`
}
type WordInfotosend struct {
	Word  string `bson:"word" json:"word"`
	Tword string `bson:"Tword" json:"Tword"`
	Info  string `bson:"Info" json:"Info"`
	Learn int    `bson:"learn" json:"learn"`
	//	Ldate time.Time `bson:"Ldate" json:"Ldate"`
}
type Wordget struct {
	Words   []WordInfotosend `bson:"words" json:"words"`
	Getdate time.Time        `json:"Getdate"`
}
type timetype struct {
	LastModified time.Time `bson:"lastModified" json:"lastModified"`
}

func (u UserInfoset) insertUser() primitive.ObjectID {
	_, col := connectDB()
	res, err := col.InsertOne(context.TODO(), u)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}
func connectDB() (*mongo.Client, *mongo.Collection) {
	client, err := mongo.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection := client.Database("dbGo").Collection("test")
	if err != nil {
		log.Fatal(err)
	}
	return client, collection
}

func disconnectDb(client *mongo.Client) {
	err := client.Disconnect(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
func insertManyWord(words []WordInfo, id primitive.ObjectID) error {

	_, col := connectDB()
	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "words", Value: bson.D{{Key: "$each", Value: words}}},
		}},
		{Key: "$currentDate", Value: bson.D{
			{Key: "lastModified", Value: true}}}}
	_, err := col.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return err
}

func findByEmail(email string) (bool, UserInfo) {
	_, col := connectDB()
	//defer disconnectDb(client) // have error when used but still work
	findOneOptions := options.FindOne()
	findOneOptions.SetProjection(bson.D{{Key: "_id", Value: 1}, {Key: "pass", Value: 1}, {Key: "uuidHash", Value: 1}})
	filter := bson.D{{Key: "email", Value: email}}
	var result UserInfo
	err := col.FindOne(nil, filter, findOneOptions).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, result
		}
		log.Fatal(err)
	}
	fmt.Print(result.UuidHash, " ::::::::: PASSSSSSS::::::: ", result.Pass)

	return true, result
}
func IsEmailExist(email string) bool {
	_, col := connectDB()
	fmt.Println("after if call 1")
	//defer disconnectDb(client) // have error when used but still work
	filter := bson.D{{Key: "email", Value: email}}
	var result UserInfo
	fmt.Println("after if call 2")
	err := col.FindOne(nil, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		fmt.Println("after if call 4")
		log.Fatal(err)
	}
	return true
}
func addHashPass(hashPass []byte, uuidHash uuid.UUID, id primitive.ObjectID) {
	_, col := connectDB()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "pass", Value: hashPass},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "uuidHash", Value: uuidHash}}}}
	_, err := col.UpdateOne(nil, filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllWords(id primitive.ObjectID) Wordget {
	_, col := connectDB()
	//defer disconnectDb(client) // have error when used but still work
	findOneOptions := options.FindOne()
	findOneOptions.SetProjection(bson.D{{Key: "_id", Value: 0}, {Key: "words", Value: 1}})
	filter := bson.D{{Key: "_id", Value: id}}
	var result Wordget
	err := col.FindOne(nil, filter, findOneOptions).Decode(&result)
	result.Getdate = time.Now()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result
		}
		log.Fatal(err)
	}
	return result
}

func syncDate(id primitive.ObjectID, timeDate time.Time) Wordget {
	_, col := connectDB()
	//defer disconnectDb(client) // have error when used but still work
	findOneOptions := options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 0}, {Key: "lastModified", Value: 1}})
	filter := bson.D{{Key: "_id", Value: id}}
	var result timetype
	var results Wordget
	err := col.FindOne(nil, filter, findOneOptions).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return results
		}
		log.Fatal(err)
	}
	if result.LastModified.Equal(timeDate) {
		return results
	}
	fmt.Println(result.LastModified, timeDate, result.LastModified.Equal(timeDate))
	cur, err := col.Aggregate(nil, bson.A{bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "_id", Value: id}}}}, bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "words", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$words"}, {Key: "as", Value: "item"}, {Key: "cond", Value: bson.D{
						{Key: "$gt", Value: bson.A{"$$item.Ldate", timeDate}}}}}}}}, {Key: "_id", Value: 0}}}}})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		if err := cur.Decode(&results); err != nil {
			log.Fatal(err)
		}

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	results.Getdate = time.Now()
	fmt.Println(results, len(results.Words))
	return results
}

func learn(id primitive.ObjectID, data []dataReq) error {
	_, col := connectDB()
	fmt.Println("  Hello   ")
	for index, elem := range data {
		fmt.Println(index, "   ", elem)
		filter := bson.D{{Key: "_id", Value: id}, {Key: "words.word", Value: elem.Word}}
		update := bson.D{
			{Key: "$inc", Value: bson.D{
				{Key: "words.$.learn", Value: elem.Learn},
			}}}

		updateResult, err := col.UpdateOne(nil, filter, update)
		if nil != err {
			fmt.Println("   Goood")
			fmt.Println(err)
			return err
		}
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	return nil
}

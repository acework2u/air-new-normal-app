package main

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type AirIotService interface {
	GetAcAll() ([]*AirIot, error)
	GetAcById(deviceSN string) (*AirIot, error)
}

type acIotService struct {
	client *mongo.Client
	ctx    context.Context
}

func NewAcIotService(ctx context.Context, client *mongo.Client) AirIotService {
	return &acIotService{ctx: ctx, client: client}
}

func (s *acIotService) GetCollection(dbName string, collectionName string) *mongo.Collection {
	collection := s.client.Database(dbName).Collection(collectionName)
	return collection
}

func (s *acIotService) GetAcAll() ([]*AirIot, error) {
	query := bson.M{}
	dbName := "airs"
	collectionName := "airs_shadows"

	airShadownCollection := s.GetCollection(dbName, collectionName)
	cursor, err := airShadownCollection.Find(ctx, query)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.ctx)

	acDevices := []*AirIoTDB{}

	for cursor.Next(s.ctx) {
		ac := &AirIoTDB{}
		ok := cursor.Decode(ac)
		if ok != nil {
			return nil, ok
		}
		acDevices = append(acDevices, ac)
	}
	var acList []*AirIot
	if len(acDevices) != 0 {
		for _, val := range acDevices {
			ac := &AirIot{
				DeviceSn:  val.DeviceSn,
				Message:   val.Message,
				Timestamp: val.Timestamp,
			}
			acList = append(acList, ac)
		}
	}

	return acList, err

}
func (s *acIotService) GetAcById(deviceSN string) (*AirIot, error) {
	panic("No Action")
}

type AirIot struct {
	DeviceSn  string    `bson:"device_sn" json:"device_sn"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Message   string    `bson:"message" json:"message"`
}

type AirIoTDB struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeviceSn  string             `bson:"device_sn" json:"device_sn"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Message   string             `bson:"message" json:"message"`
}

var (
	ctx    context.Context
	DB_URL = os.Getenv("DB_URL")
	client *mongo.Client
)

func dbConnect() error {
	var err error

	DB_URL = os.Getenv("DB_URL")
	if len(DB_URL) == 0 {
		err = gotenv.Load()
		if err != nil {
			log.Panicln(err.Error())
			return err
		}
		DB_URL = os.Getenv("DB_URL")
	}

	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(DB_URL))
	if err != nil {
		return err
	}

	return nil
}
func main() {

	// 1. Load db init config
	err := dbConnect()

	if err != nil {
		log.Fatal(err.Error())
	}
	// 2. fins Ac all

	acIot := NewAcIotService(ctx, client)
	acList, err := acIot.GetAcAll()
	if err != nil {
		fmt.Println(err.Error())
	}

	resp := []AirIot{}
	for _, ac := range acList {
		acVal := AirIot{
			DeviceSn:  ac.DeviceSn,
			Message:   ac.Message,
			Timestamp: ac.Timestamp,
		}
		resp = append(resp, acVal)
	}
	//
	if len(resp) != 0 {
		for _, val := range resp {
			acData, er := GetClaimsFromToken(val.Message)
			if er != nil {
				log.Println(er)
			}

			acDecode := DecodeValAcShadow(acData)
			fmt.Println(acDecode.Reg2000)

		}

	}

	//fmt.Println(resp)

}

func GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {

	var secretKey = "SaijoDenkiSmartIOT"
	secKey := secretKey

	var secret = []byte(secKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return secret, nil
	})

	if err != nil {

		fmt.Println("Error :" + err.Error())
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// Helper
type AcStr struct {
	reg1000 []byte
	reg2000 []byte
	reg3000 []byte
	reg4000 []byte
}
type AcValReq struct {
	Reg1000 string
	Reg2000 string
	Reg3000 string
	Reg4000 string
}

func DecodeValAcShadow(valShadow jwt.MapClaims) *AcValReq {
	regis1000 := valShadow["data"].(map[string]interface{})["reg1000"].(string)
	regis2000 := valShadow["data"].(map[string]interface{})["reg2000"].(string)
	regis3000 := valShadow["data"].(map[string]interface{})["reg3000"].(string)
	regis4000 := valShadow["data"].(map[string]interface{})["reg4000"].(string)
	acValReq := &AcValReq{
		Reg1000: regis1000,
		Reg2000: regis2000,
		Reg3000: regis3000,
		Reg4000: regis4000,
	}
	return acValReq
}

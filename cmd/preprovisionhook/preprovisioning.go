package main

import (
	"context"
	"encoding/hex"
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

type OzoneWorkStage struct {
	DeviceSn   string     `json:"device_sn"`
	OzoneStage OzoneStage `json:"ozone_stage"`
}
type OzoneStage struct {
	Stage    string `json:"stage"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

func StageTitle(stage int32) string {
	switch stage {
	case 1, 2, 3, 4, 5:
		return fmt.Sprintf("ขั้นตอนที่ %v", stage)
	default:
		return "Ready"

	}
}

func DurationTitle(duration int32) string {
	return fmt.Sprintf("เวลาดำเนินการ %vนาที", duration)
}

func Title(stage int32) string {
	switch stage {
	case 1:
		return "ไล่ความชื้นในคอยล์"
	case 2:
		return "เครื่องสร้างโอโซนทำงาน"
	case 3:
		return "เครื่องคงปริมาณโอโซน"
	case 4:
		return "เครื่องทำการสลายโอโซน"
	case 5:
		return "เครื่องทำการสลายโอโซนเสร็จแล้ว"

	default:
		return "เครื่องสร้างโอโซนหยุดทำงาน สถานะปกติ"

	}
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
	for {

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
					return
				}
				acDecode := DecodeValAcShadow(acData)
				ac2000 := NewGetAcValue(acDecode)
				ozone := ac2000.Ac2000()
				ozoneEvent := &OzoneWorkStage{
					DeviceSn:   val.DeviceSn,
					OzoneStage: OzoneStage{Stage: StageTitle(ozone.Stage), Title: Title(ozone.Stage), Duration: DurationTitle(ozone.Duration)},
				}
				//fmt.Println(val.DeviceSn, " time= ", ozone.Duration, " stage =", ozone.Stage, "active=", time.Now())
				fmt.Printf("device =%v stage=%v title=%v duration=%v \n", ozoneEvent.DeviceSn, ozoneEvent.OzoneStage.Stage, ozoneEvent.OzoneStage.Title, ozoneEvent.OzoneStage.Duration)
			}

		}
		//time.Sleep(8 * time.Second)
		time.Sleep(20 * time.Second)

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
type AcValue interface {
	Ac2000() *OzoneEvent
}

type OzoneEvent struct {
	Stage    int32
	Duration int32
}

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

func NewGetAcValue(reg *AcValReq) AcValue {
	data1000, _ := hex.DecodeString(reg.Reg1000)
	data2000, _ := hex.DecodeString(reg.Reg2000)
	data3000, _ := hex.DecodeString(reg.Reg3000)
	data4000, _ := hex.DecodeString(reg.Reg4000)

	return &AcStr{
		reg1000: data1000,
		reg2000: data2000,
		reg3000: data3000,
		reg4000: data4000,
	}
}
func (h *AcStr) Ac2000() *OzoneEvent {

	data := &OzoneEvent{}
	//fmt.Println(h.reg2000[18:])
	if len(h.reg2000) == 20 {
		ac := h.reg2000[18:]

		data = &OzoneEvent{
			Stage:    int32(ac[0]),
			Duration: int32(ac[1]),
		}

	}

	return data

	//dataOzone := OzoneEvent{
	//	Duration: int(duration),
	//}
	//
	//panic("No Action")
}

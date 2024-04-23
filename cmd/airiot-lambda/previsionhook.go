package main

import (
	"context"
	"fmt"
	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var (
	Ctx        context.Context
	DBUrl      string
	DBName     string
	Collection string
	client     *mongo.Client
)

type AirIotAuthentication interface {
	AcList() ([]*AcInfoDB, error)
	GetAirIoTCert()
	GetCollection(dbName string, collectionName string)
}

type airAwsIot struct {
	ctx                  context.Context
	client               *mongo.Client
	airProductCollection *mongo.Collection
}

type AcInfo struct {
	DeviceSn string `json:"device_sn"`
}
type AcInfoDB struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeviceSn string             `json:"serial" bson:"serial"`
}

func NewAirAWSIoT(ctx context.Context, client *mongo.Client) AirIotAuthentication {
	return &airAwsIot{ctx: ctx, client: client}
}

func (s *airAwsIot) GetCollection(dbName string, collectionName string) {
	s.airProductCollection = s.client.Database(dbName).Collection(collectionName)

}
func (s *airAwsIot) AcList() ([]*AcInfoDB, error) {

	//query := bson.M{"status": true}
	query := bson.M{}

	cursor, err := s.airProductCollection.Find(s.ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(s.ctx)

	acList := []*AcInfoDB{}
	for cursor.Next(s.ctx) {
		acInfo := &AcInfoDB{}
		ok := cursor.Decode(acInfo)
		if ok != nil {
			log.Println(ok)
			return nil, ok
		}

		acList = append(acList, acInfo)
	}

	return acList, err

}
func (s *airAwsIot) GetAirIoTCert() {
	panic("Noo Action")
}

type AirProduction struct {
	DeviceSn   string    `json:"device_sn"`
	Cert       string    `json:"cert"`
	Active     bool      `json:"active"`
	Status     bool      `json:"status"`
	RegisterAt time.Time `json:"registerAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
type PreProvisionHookInput struct {
	ClaimCertificateId string            `json:"claimCertificateId"`
	CertificateId      string            `json:"certificateId"`
	CertificatePem     string            `json:"certificatePem"`
	TemplateArn        string            `json:"templateArn"`
	ClientId           string            `json:"clientId"`
	Parameters         map[string]string `json:"parameters"`
}

type PreProvisionHookOutput struct {
	AllowProvisioning  bool
	ParameterOverrides map[string]string
}

func dbConnect() error {
	var err error
	DBUrl = os.Getenv("DB_URL")
	if len(DBUrl) == 0 {
		err = gotenv.Load()
		if err != nil {
			log.Fatal(err)
			return err
		}
		DBUrl = os.Getenv("DB_URL")
	}

	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(DBUrl))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func init() {
	err := dbConnect()
	if err != nil {
		log.Fatal(err)
	}

	Collection = "products"
	// Set collection
	DBName = "airs"

}

func HandlerRequest(ctx context.Context, event PreProvisionHookInput) {
	defer client.Disconnect(ctx)
	cl := NewAirAWSIoT(ctx, client)
	cl.GetCollection(DBName, Collection)
	//acList,err := cl.AcList()
	//if er

} // end of HandlerRequest

func main() {

	// Cet all data in Collection
	cl := NewAirAWSIoT(Ctx, client)
	cl.GetCollection(DBName, Collection)
	acList, err := cl.AcList()
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range acList {

		fmt.Println(val.DeviceSn)

	}

	fmt.Println(len(acList))

}

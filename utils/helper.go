package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"strings"
	"time"
)

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

type powerFunc func(val int) string
type modeFunc func(val int) string
type tempFunc func(val int) string
type roomTempFunc func(val int) string
type setRhFunc func(val int) string
type roomRhFunc func(val int) string
type fanSpeedFunc func(val int) string
type louverFunc func(val int) string

type AC1000 struct {
	Power    powerFunc    `json:"power"`
	Mode     modeFunc     `json:"mode"`
	Temp     tempFunc     `json:"temp"`
	RoomTemp roomTempFunc `json:"roomTemp"`
	SetRh    setRhFunc    `json:"setRh"`
	RoomRh   roomRhFunc   `json:"roomRh"`
	FanSpeed fanSpeedFunc `json:"fanSpeed"`
	Louver   louverFunc   `json:"louver"`
}

// Indoor value register2000
type midCoilFunc func(val int) int
type outletFunc func(val int) int
type exvPositionFunc func(val int) int
type demandFunc func(val int) int
type indCapacityFunc func(val int) int
type statusFunc func(val int) int
type protectionFunc func(val int) int
type uvTimeRunningFunc func(val int) int
type co2Func func(val int) int
type dischargeFunc func(val int) int
type ambientFunc func(val int) int
type compActualFunc func(val int) int
type compCurrentFunc func(val []byte) int
type comStatus func(val []byte) int
type oduErrorFunc func(val []byte) int
type oduSuction func(val []byte) int
type oduFunc[T int | []byte] func(val T) int

type AC2000 struct {
	MidCoilTemp   midCoilFunc       `json:"midCoilTemp"`
	OutletTemp    outletFunc        `json:"outletTemp"`
	ExvPosition   exvPositionFunc   `json:"exvPosition"`
	Demand        demandFunc        `json:"demand"`
	IndCapacity   indCapacityFunc   `json:"indCapacity"`
	LedStatus     statusFunc        `json:"ledStatus"`
	Protection    protectionFunc    `json:"protection"`
	UvTimeRunning uvTimeRunningFunc `json:"uvTimeRunning"`
	Co2           co2Func           `json:"co2"`
}
type AC3000 struct {
	MidCoilTemp midCoilFunc     `json:"midCoilTemp"`
	OutletTemp  outletFunc      `json:"outletTemp"`
	Discharge   dischargeFunc   `json:"discharge"`
	Ambient     ambientFunc     `json:"ambient"`
	Suction     oduFunc[[]byte] `json:"suction"`
	CompActual  compActualFunc  `json:"compActual"`
	CompCurrent compCurrentFunc `json:"compCurrent"`
	Demand      demandFunc      `json:"demand"`
	StatusComp  comStatus       `json:"statusComp"`
	OudErrors   oduErrorFunc    `json:"oudErrors"`
}

type IndoorInfo struct {
	Power    string `json:"power"`
	Mode     string `json:"mode"`
	Temp     string `json:"temp"`
	RoomTemp string `json:"roomTemp"`
	RhSet    string `json:"rhSet"`
	RhRoom   string `json:"RhRoom"`
	FanSpeed string `json:"fanSpeed"`
	Louver   string `json:"louver"`
}

type Ind2000 struct {
	MidcoilTemp int `json:"midcoilTemp"`
	OutletTemp  int `json:"outletTemp"`
	ExvPosition int `json:"exvPosition"`
	Demand      int `json:"demand"`
	Capacity    int `json:"capacity"`
	LedStatus   int `json:"ledStatus"`
	Protection  int `json:"protection"`
}
type Out3000 struct {
	MidcoilTemp int `json:"midcoilTemp"`
	OutletTemp  int `json:"outletTemp"`
	Discharge   int `json:"discharge"`
	Ambient     int `json:"ambient"`
	Suction     int `json:"suction"`
	CompActual  int `json:"compActual"`
	CompCurrent int `json:"compCurrent"`
	Demand      int `json:"demand"`
	StatusComp  int `json:"statusComp"`
	LedStatus   int `json:"ledStatus"`
}

type AcValue interface {
	Ac1000() *IndoorInfo
	Ac2000() *Ind2000
	Ac3000() *Out3000
}

type AcStr struct {
	reg1000 []byte
	reg2000 []byte
	reg3000 []byte
	reg4000 []byte
}

func NewGetAcVal(reg1000 string) AcValue {
	data, err := hex.DecodeString(reg1000)
	if err != nil {
		panic(err)
	}

	return &AcStr{reg1000: data}
}
func NewAcVal(reg string, payload string) AcValue {

	data, err := hex.DecodeString(payload)
	if err != nil {
		panic(err)
	}
	//log.Println(payload)
	//log.Println("data")
	//log.Println(payload)
	//log.Println(len(payload))
	//log.Println(data)

	switch reg {
	case "1":
		return &AcStr{reg1000: data}
	case "2":
		return &AcStr{reg2000: data}
	case "3":
		return &AcStr{reg3000: data}
	case "4":
		return &AcStr{reg4000: data}

	default:
		return &AcStr{}
	}
	return &AcStr{}
}
func (ut *AcStr) Ac1000() *IndoorInfo {
	ac := &AC1000{
		Power:    power,
		Mode:     mode,
		Temp:     temp,
		RoomTemp: roomTemp,
		SetRh:    rh,
		RoomRh:   rh,
		FanSpeed: fanSpeed,
		Louver:   louver,
	}
	rs := &IndoorInfo{
		Power:    ac.Power(int(ut.reg1000[1])),
		Mode:     ac.Mode(int(ut.reg1000[3])),
		Temp:     ac.Temp(int(ut.reg1000[5])),
		RoomTemp: ac.RoomTemp(int(ut.reg1000[7])),
		RhSet:    ac.SetRh(int(ut.reg1000[9])),
		RhRoom:   ac.RoomRh(int(ut.reg1000[11])),
		FanSpeed: ac.FanSpeed(int(ut.reg1000[13])),
		Louver:   ac.Louver(int(ut.reg1000[15])),
	}

	return rs
}
func (ut *AcStr) Ac2000() *Ind2000 {

	//testReg := ut.reg2000[12:14]

	//log.Println("Reg ")
	//log.Println(len(ut.reg2000))
	//log.Println("Mid-Coil =", ut.reg2000[0:2])
	//log.Println("Out let =", ut.reg2000[2:4])
	//log.Println("ExV =", ut.reg2000[4:6])
	//log.Println("Demand =", ut.reg2000[6:8])
	//log.Println("Capacity =", ut.reg2000[8:10])
	//log.Println("Led Status =", ut.reg2000[10:12])
	//log.Println("Protection =", ut.reg2000[12:14])
	//log.Println("Pm2.5 =", ut.reg2000[14:16])
	//log.Println("PM2.5 =", ut.reg2000[16:18])
	//log.Println("CO2 =", ut.reg2000[18:])

	//data := binary.BigEndian.Uint64(testReg)
	//fmt.Println(data)

	ac := &AC2000{
		MidCoilTemp:   midCoilTemp,
		OutletTemp:    outletTemp,
		ExvPosition:   exvPosition,
		Demand:        demand,
		LedStatus:     ledStatus,
		IndCapacity:   indCapacity,
		Protection:    protection,
		UvTimeRunning: uvTimeRunning,
		Co2:           indCo2,
	}

	rs := &Ind2000{
		MidcoilTemp: ac.MidCoilTemp(int(ut.reg2000[1])),
		OutletTemp:  ac.OutletTemp(int(ut.reg2000[3])),
		ExvPosition: ac.ExvPosition(int(ut.reg2000[5])),
		Demand:      ac.Demand(int(ut.reg2000[7])),
		Capacity:    ac.IndCapacity(int(ut.reg2000[9])),
		LedStatus:   ac.LedStatus(int(ut.reg2000[11])),
		Protection:  ac.Protection(int(ut.reg2000[13])),
	}

	return rs
}
func (ut *AcStr) Ac3000() *Out3000 {
	ac := &AC3000{
		MidCoilTemp: midCoilTemp,
		OutletTemp:  outletTemp,
		Discharge:   oduDischarge,
		Ambient:     oudAmbient,
		Suction:     oduSectionFunc,
		CompActual:  oduComActual,
		CompCurrent: oduComCurrent,
		Demand:      oduDemand,
		StatusComp:  oduCompStatus,
		OudErrors:   oduErrorLed,
	}

	if len(ut.reg3000) == 80 {

		//log.Println("REg 3000")
		//log.Println("Len == >", len(ut.reg3000))
		//log.Println(ut.reg3000)
		//log.Println("Mid-Coil =", ut.reg3000[0:2])
		//log.Println("Co2 Low =", ut.reg3000[16:18])
		//log.Println("Co2 Hi =", ut.reg3000[18:20])
		//log.Println("Outlet =", ut.reg3000[2:4])
		//log.Println("Discharge =", ut.reg3000[4:6])
		//log.Println("Suction =", ut.reg3000[11])
		//log.Println("Suction =", ut.reg3000[10:12])
		//log.Println("Ambient =", ut.reg3000[8:10])
		//log.Println("Comp Actual =", ut.reg3000[20:22])
		//log.Println("Comp Current =", ut.reg3000[24:28])
		//log.Println("Demand  =", ut.reg3000[40:42])
		//log.Println("Demand  =", ut.reg3000[41])

		rs := &Out3000{
			MidcoilTemp: ac.MidCoilTemp(int(ut.reg3000[1])),
			OutletTemp:  ac.OutletTemp(int(ut.reg3000[3])),
			Discharge:   ac.Discharge(int(ut.reg3000[5])),
			Ambient:     ac.Ambient(int(ut.reg3000[9])),
			Suction:     ac.Suction(ut.reg3000[10:12]),
			CompActual:  ac.CompActual(int(ut.reg3000[19])),
			CompCurrent: ac.CompCurrent(ut.reg3000[24:28]),
			Demand:      ac.Demand(int(ut.reg3000[41])),
			StatusComp:  ac.StatusComp(ut.reg3000[44:46]),
			LedStatus:   ac.OudErrors(ut.reg3000[72:76]),
		}
		return rs
	}
	return &Out3000{}

}

func power(val int) string {
	powerTxt := ""
	switch val {
	case 0:
		powerTxt = "off"
	case 1:
		powerTxt = "on"
	default:
		powerTxt = "err"
	}
	return powerTxt
}
func mode(val int) string {
	displayTxt := ""
	switch val {
	case 0:
		displayTxt = "cool"
	case 1:
		displayTxt = "dry"
	case 2:
		displayTxt = "auto"
	case 3:
		displayTxt = "heat"
	case 4:
		displayTxt = "fan"
	default:
		displayTxt = "err"
	}
	return displayTxt
}
func temp(val int) string {
	displayTxt := ""

	if val < 0 || val > 60 {
		displayTxt = "err"
		return displayTxt
	}
	val2 := float64(val)
	tempVal := val2 / 2
	s := fmt.Sprintf("%3.1f", tempVal)
	displayTxt = s

	return displayTxt

}
func roomTemp(val int) string {
	displayTxt := ""

	if val < 0 || val > 255 {
		displayTxt = "err"
		return displayTxt
	}

	val2 := float64(val)
	tempVal := val2 / 4
	s := fmt.Sprintf("%3.2f", tempVal)
	displayTxt = s

	return displayTxt

}
func rh(val int) string {
	displayTxt := ""

	if val < 0 || val > 100 {
		displayTxt = "err"
		return displayTxt
	}

	displayTxt = fmt.Sprintf("%v", val)

	return displayTxt

}
func fanSpeed(val int) string {
	displayTxt := ""
	//Value 0 : Fan Auto
	//Value 1 : Fan Low
	//Value 2 : Fan Med
	//Value 3 : Fan High
	//Value 4 : Fan Hi Hi
	//Value 5 : Fan Turbo
	switch val {
	case 0:
		displayTxt = "auto"
	case 1:
		displayTxt = "low"
	case 2:
		displayTxt = "med"
	case 3:
		displayTxt = "high"
	case 4:
		displayTxt = "high+"
	default:
		displayTxt = "turbo"
	}
	return displayTxt
}
func louver(val int) string {
	displayTxt := ""
	//Value 0 :  Auto (Swing)
	//Value 1 :  Level 1
	//Value 2 :  Level 2
	//Value 3 :  Level 3
	//Value 4 :  Level 4
	//Value 5 :  Level 5

	switch val {
	case 0:
		displayTxt = "auto"
	case 1:
		displayTxt = "level 1"
	case 2:
		displayTxt = "level 2"
	case 3:
		displayTxt = "level 3"
	case 4:
		displayTxt = "level 4"
	case 5:
		displayTxt = "level 5"
	default:
		displayTxt = "err"
	}
	return displayTxt
}
func midCoilTemp(val int) int {
	midTemp := 0

	if val == 255 {
		return val
	}

	if val >= 0 {
		midTemp = val - 40
	}

	return midTemp
}
func outletTemp(val int) int {
	oultTemp := 0
	if val == 255 {
		return val
	}
	if val >= 0 {
		oultTemp = val - 40
	}
	return oultTemp
}
func exvPosition(val int) int {
	exv := 0
	exv = val
	return exv
}
func demand(val int) int {
	demand := 0
	if fmt.Sprintf("%T", val) == "int" {
		demand = val
	}
	return demand

}
func indCapacity(val int) int {
	capacity := 0

	capacity = val * 1000

	return capacity
}
func ledStatus(val int) int {
	return val
}
func protection(val int) int {
	return val
}
func uvTimeRunning(val int) int {
	return val
}
func indCo2(val int) int {
	return val
}
func oduDischarge(val int) int {
	discharge := 0
	if val == 255 {
		return val
	}

	if val > 0 && val != 255 {
		discharge = val - 40
	}
	return discharge
}
func oudAmbient(val int) int {
	ambient := 0
	if val == 255 {
		return val
	}
	if val >= 0 {
		ambient = val - 40
	}
	return ambient
}
func oduConvTemp(val int) int {
	temp := 0
	if val == 255 {
		return val
	}
	temp = val - 40

	return temp
}
func oduComActual(val int) int {
	return val
}
func oduDemand(val int) int {
	return val
}
func oduCompStatus(val []byte) int {

	status := 0
	if len(val) == 2 {

		status = int(val[1])
	}

	return status
}
func oduComCurrent(val []byte) int {

	comCurrent := 0
	if len(val) == 4 {
		z0 := int(val[0]) * 1
		z1 := int(val[1]) * 2
		z2 := int(val[2]) * 4
		z3 := int(val[3]) * 8

		comCurrent = (z0 + z1 + z2 + z3) / 10

	} else {
		comCurrent = -1
	}
	return comCurrent
}

func oduErrorLed(val []byte) int {
	errCode := 0
	if len(val) == 4 {
		for _, v := range val {
			errCode = errCode + int(v)
		}
	}

	return errCode
}

func oduSectionFunc(val []byte) int {
	sectionTemp := 0
	if len(val) == 2 {
		pipTemp := int(val[1])
		if pipTemp == 255 {
			return pipTemp
		}
		sectionTemp = int(val[1]) - 40
	}

	return sectionTemp
}

//Custom Datetime

type ConvDateDB time.Time

func (v ConvDateDB) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(v))
}

func (v *ConvDateDB) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	rv := bson.RawValue{
		Type:  t,
		Value: b,
	}
	res := time.Time{}
	if err := rv.Unmarshal(&res); err != nil {
		return err
	}
	*v = ConvDateDB(res)
	return nil

}
func (v ConvDateDB) String() string {
	return time.Time(v).String()
}

type CustomTime struct {
	time.Time
}
type TestModel struct {
	Date CustomTime `json:"date"`
}

func (t CustomTime) MarshalJSON() ([]byte, error) {
	date := t.Time.Format("2006-01-02")
	date = fmt.Sprintf(`"%s"`, date)
	return []byte(date), nil
}
func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func getInt1(s []byte) int {
	var b [8]byte
	copy(b[8-len(s):], s)
	return int(binary.BigEndian.Uint64(b[:]))
}

func getInt2(s []byte) int {
	var res int
	for _, v := range s {
		res <<= 8
		res |= int(v)
	}
	return res
}

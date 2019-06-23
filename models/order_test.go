package models

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"order-service/pkgs/e"
	"order-service/services/distance"
	"testing"
)

var ErrBadDriver = errors.New("driver: bad connection")

func SetupTests() {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	db, _ = gorm.Open(mocket.DriverName, "connection_string")
}

func TestCreateOrder(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	var (
		expectDistance = rand.Intn(100)
		expectedId     = rand.Int63n(100)
	)
	// init the mock calculator
	distance.InitMockCalculator(expectDistance, nil)

	// mock the query that create order
	mocket.Catcher.Reset().NewMock().WithQuery(`INSERT  INTO "orders" ("distance","status") VALUES (?,?)`).WithID(expectedId)

	var (
		src = []string{"1", "2"}
		des = []string{"1.5", "1.6"}
	)

	o , err := CreateOrder(src, des)

	// check if the order return without error
	a.Nil(err, "order should be created without err")
	a.NotNil(o, "order should be created")
	a.Equal(expectedId, o.ID, "id should be as expected")
	a.Equal(expectDistance, o.Distance, "distance should be as expected")
	a.Equal(StatusUnassigned, o.Status, "status should be UNASSIGNED")
}

func TestCreateOrder_Query_Exception(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	// init the mock calculator
	var expectDistance = 100
	distance.InitMockCalculator(expectDistance, nil)

	// mock the query that create order with exception
	mocket.Catcher.Reset().NewMock().WithQuery(`INSERT  INTO "orders" ("distance","status") VALUES (?,?)`).WithExecException()

	var (
		src = []string{"1", "2"}
		des = []string{"1.5", "1.6"}
	)

	o , err := CreateOrder(src, des)

	a.NotNil(err, "error should occur based on the query")
	a.Equal(ErrBadDriver, err, "error should the expected error")
	a.Nil(o, "order should not be created")
}

func TestCreateOrder_Service_Exception(t *testing.T)  {
	a := assert.New(t)

	SetupTests()

	// init the mock calculator with error
	var expectedErr = errors.New("test for service exception")
	distance.InitMockCalculator(0, expectedErr)

	var (
		src = []string{"1", "2"}
		des = []string{"1.5", "1.6"}
	)

	o , err := CreateOrder(src, des)

	a.Equal(expectedErr ,err, "error should the expected error")
	a.Nil(o, "order should not be created")
}

func TestGetOrders(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	// set up expected result
	order1 := Order{ID: 1, Status:StatusUnassigned, Distance:rand.Intn(5000)}
	order2 := Order{ID: 2, Status:StatusTaken, Distance:rand.Intn(5000)}
	orders := []Order{order1, order2}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that query the orders
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "orders"   LIMIT 1 OFFSET 1`).WithReply(expectMap)

	results, err := GetOrders(1, 1)

	a.Nil(err, "order should be created without err")
	a.NotNil(results, "result should not be nil")
	a.Equal(len(orders), len(results), "length of the result should match expected")
	a.Equal(order1, *results[0], "Order 1 should match exactly")
	a.Equal(order2, *results[1], "Order 2 should match exactly")
}

func TestGetOrders_Query_Exception(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	// mock the query that query the orders with exception
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "orders"   LIMIT 1 OFFSET 1`).WithQueryException()

	os , err := GetOrders(1, 1)
	a.NotNil(err, "error should occur based on the query")
	a.Equal(ErrBadDriver, err, "error should the expected error")
	a.Nil(os, "order should not be returned")
}

func TestTakeOrder(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	// set up expected result
	const orderId = 1
	order1 := Order{ID: orderId, Status:StatusUnassigned, Distance:rand.Intn(5000)}
	orders := []Order{order1}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that get the order by id
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`).WithReply(expectMap)

	err := TakeOrder(orderId)
	a.Nil(err, "error should be nil")
}

func TestTakeOrder_Query_Exception_On_Select(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	// mock the query that get the order by id with exception
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`).WithQueryException()

	err := TakeOrder(rand.Int63n(100))
	a.NotNil(err, "error should be returned")
	a.Equal(ErrBadDriver, err, "error should the expected error")
}

func TestTakeOrder_Query_Exception_On_Update(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	// set up expected result
	const orderId = 1
	order1 := Order{ID: orderId, Status:StatusUnassigned, Distance:rand.Intn(5000)}
	orders := []Order{order1}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that get the order by id
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`).WithReply(expectMap)

	// mock the query which update the order
	mocket.Catcher.NewMock().WithQuery(`UPDATE "orders" SET "status" = ?  WHERE "orders"."id" = ?`).WithExecException()

	err := TakeOrder(orderId)
	a.NotNil(err, "error should be returned")
	a.Equal(ErrBadDriver, err, "error should the expected error")
}

func TestTakeOrder_Already_Taken(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	// set up expected result
	const orderId = 1
	order1 := Order{ID: orderId, Status:StatusTaken, Distance:rand.Intn(5000)}
	orders := []Order{order1}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that get the order by id
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`).WithReply(expectMap)

	err := TakeOrder(orderId)
	a.NotNil(err, "error should be returned")
	a.Equal(e.ErrOrderAlreadyTaken, err, "error should not found from gorm")
}

func TestTakeOrder_Not_Exist(t *testing.T) {
	a := assert.New(t)

	SetupTests()

	// mock the query that get the order by id
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`)

	err := TakeOrder(rand.Int63n(100))
	a.NotNil(err, "error should be returned")
	a.Equal(gorm.ErrRecordNotFound, err, "error should not found from gorm")
}
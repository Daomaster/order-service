package api

import (
	"bytes"
	"encoding/json"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"order-service/api/order"
	"order-service/api/requests"
	"order-service/models"
	"order-service/pkgs/e"
	"order-service/services/distance"
	"testing"
)

// helper function to parse json
func parseJson(b io.Reader, i interface{}) error {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		return err
	}

	return nil
}

// helper function to create json io
func createJson(i interface{}) (io.Reader, error) {
	body, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(body), err
}

func TestCreateOrder(t *testing.T)  {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// init the mock calculator
	var expectDistance = rand.Intn(5000)
	distance.InitMockCalculator(expectDistance, nil)

	// get the router
	r := InitRouter()

	// create request body
	var createOrder requests.CreateOrderRequest
	createOrder.Origin = []string{"35.9984617", "-115.1432558"}
	createOrder.Destination = []string{"36.0222811", "-115.0980736"}
	reqBody, err := createJson(createOrder)

	a.Nil(err, "should not have problem with create json")

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPost, "/orders", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusOK, w.Code, "server should return back 200 OK")

	// parse the response
	var orderResponse models.Order
	err = parseJson(w.Body, &orderResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(orderResponse, "server should have return the order")
	a.Equal(expectDistance, orderResponse.Distance, "server should return the correct distance")
	a.Equal(models.StatusUnassigned, orderResponse.Status, "server should create order with default UNASSIGNED")
}

func TestCreateOrder_Bad_Request(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// init the mock calculator
	var expectDistance = rand.Intn(5000)
	distance.InitMockCalculator(expectDistance, nil)

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPost, "/orders", nil)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusBadRequest, w.Code, "server should return back 400 Bad Request")

	// parsing the error response
	var errorResponse e.ResponseError
	err := parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrOrderRequestInvalid.Error(), errorResponse.Error, "error response should match the error content")
}

func TestCreateOrder_Internal_Sever_Error(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// init the mock calculator
	var expectDistance = rand.Intn(5000)
	distance.InitMockCalculator(expectDistance, nil)

	// mock the query that create order with exception
	mocket.Catcher.NewMock().WithQuery(`INSERT  INTO "orders" ("distance","status") VALUES (?,?)`).WithExecException()
	defer mocket.Catcher.Reset()

	// create request body
	var createOrder requests.CreateOrderRequest
	createOrder.Origin = []string{"35.9984617", "-115.1432558"}
	createOrder.Destination = []string{"36.0222811", "-115.0980736"}
	reqBody, err := createJson(createOrder)

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPost, "/orders", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusInternalServerError, w.Code, "server should return back 500 Internal Sever Error")

	// parsing the error response
	var errorResponse e.ResponseError
	err = parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrInternalError.Error(), errorResponse.Error, "error response should match the error content")
}

func TestGetOrders(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// set up expected result
	order1 := models.Order{ID: 1, Status:models.StatusUnassigned, Distance:rand.Intn(5000)}
	order2 := models.Order{ID: 2, Status:models.StatusTaken, Distance:rand.Intn(5000)}
	orders := []models.Order{order1, order2}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that query the orders
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "orders"   LIMIT 2 OFFSET 0`).WithReply(expectMap)
	defer mocket.Catcher.Reset()

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodGet, "/orders", nil)

	// create the query string
	q := req.URL.Query()
	q.Add("page", "0")
	q.Add("limit", "2")
	req.URL.RawQuery = q.Encode()
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusOK, w.Code, "server should return back 200 OK")

	// parsing the response
	var ordersResponse []models.Order
	err := parseJson(w.Body, &ordersResponse)
	a.Nil(err, "should not error out upon parsing response")
	a.NotNil(ordersResponse, "server should have a response")
	a.Equal(len(orders), len(ordersResponse), "response should return 2 orders")
	a.Equal(order1, ordersResponse[0], "order 1 should match exactly")
	a.Equal(order2, ordersResponse[1], "order 2 should match exactly")
}

func TestGetOrders_No_Query(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodGet, "/orders", nil)
	r.ServeHTTP(w, req)

	// it should return 200 since default value would 0, 0
	// check response code
	a.Equal(http.StatusOK, w.Code, "server should return back 200 OK")

	// parsing the response
	var ordersResponse []models.Order
	err := parseJson(w.Body, &ordersResponse)
	a.Nil(err, "should not error out upon parsing response")
	a.NotNil(ordersResponse, "server should have a response")
	a.Equal(0, len(ordersResponse), "response should be an empty array")
}

func TestGetOrders_Invalid_Query_Type(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodGet, "/orders", nil)

	// create the query string
	q := req.URL.Query()
	q.Add("page", "test")
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusBadRequest, w.Code, "server should return back 400 Bad Request")

	// parsing the error response
	var errorResponse e.ResponseError
	err := parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrQueryStringInvalid.Error(), errorResponse.Error, "error response should match the error content")
}

func TestGetOrders_Invalid_Query_Value(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodGet, "/orders", nil)

	// create the query string
	q := req.URL.Query()
	q.Add("page", "-1")
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusBadRequest, w.Code, "server should return back 400 Bad Request")

	// parsing the error response
	var errorResponse e.ResponseError
	err := parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrQueryStringInvalid.Error(), errorResponse.Error, "error response should match the error content")
}

func TestGetOrders_Internal_Server_Error(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// mock the query that query the orders
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "orders"   LIMIT 2 OFFSET 0`).WithQueryException()
	defer mocket.Catcher.Reset()

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodGet, "/orders", nil)

	// create the query string
	q := req.URL.Query()
	q.Add("page", "0")
	q.Add("limit", "2")
	req.URL.RawQuery = q.Encode()
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusInternalServerError, w.Code, "server should return back 500 Internal Sever Error")

	// parsing the error response
	var errorResponse e.ResponseError
	err := parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrInternalError.Error(), errorResponse.Error, "error response should match the error content")
}

func TestUpdateOrder(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// set up expected result
	const orderId = 1
	order1 := models.Order{ID: orderId, Status:models.StatusUnassigned, Distance:rand.Intn(5000)}
	orders := []models.Order{order1}

	// get the router
	r := InitRouter()

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that get the order by id
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`).WithReply(expectMap)
	defer mocket.Catcher.Reset()

	// create request body
	var takeOrderRequest requests.TakeOrderRequest
	takeOrderRequest.Status = models.StatusTaken
	reqBody, err := createJson(takeOrderRequest)

	a.Nil(err, "should not have problem with create json")

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPatch, "/orders/1", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusOK, w.Code, "server should return back 200 OK")

	// parsing response
	var takeOrderResponse order.TakeOrderResponse
	err = parseJson(w.Body, &takeOrderResponse)
	a.Nil(err, "should not error out upon parsing response")
	a.NotNil(takeOrderRequest, "server should have a response")
	a.Equal(models.StatusSuccess, takeOrderResponse.Status, "response should contain SUCCESS")
}

func TestUpdateOrder_Already_Taken(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// set up expected result
	const orderId = 1
	order1 := models.Order{ID: orderId, Status:models.StatusTaken, Distance:rand.Intn(5000)}
	orders := []models.Order{order1}

	// get the router
	r := InitRouter()

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that get the order by id
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`).WithReply(expectMap)
	defer mocket.Catcher.Reset()

	// create request body
	var takeOrderRequest requests.TakeOrderRequest
	takeOrderRequest.Status = models.StatusTaken
	reqBody, err := createJson(takeOrderRequest)

	a.Nil(err, "should not have problem with create json")

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPatch, "/orders/1", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusConflict, w.Code, "server should return back 409 Conflict")

	// parsing the error response
	var errorResponse e.ResponseError
	err = parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrOrderAlreadyTaken.Error(), errorResponse.Error, "error response should match the error content")
}

func TestUpdateOrder_Not_Found(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// set up expected result
	const orderId = 1
	order1 := models.Order{ID: orderId, Status:models.StatusTaken, Distance:rand.Intn(5000)}
	orders := []models.Order{order1}

	// get the router
	r := InitRouter()

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that get the order by id
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`)
	defer mocket.Catcher.Reset()

	// create request body
	var takeOrderRequest requests.TakeOrderRequest
	takeOrderRequest.Status = models.StatusTaken
	reqBody, err := createJson(takeOrderRequest)

	a.Nil(err, "should not have problem with create json")

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPatch, "/orders/1", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusNotFound, w.Code, "server should return back 404 Not Found")

	// parsing the error response
	var errorResponse e.ResponseError
	err = parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrOrderNotExist.Error(), errorResponse.Error, "error response should match the error content")
}

func TestUpdateOrder_No_ID(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPatch, "/orders", nil)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusNotFound, w.Code, "server should not even match the route")
}

func TestUpdateOrder_Invalid_ID(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// get the router
	r := InitRouter()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPatch, "/orders/testtest", nil)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusBadRequest, w.Code, "server should return back 400 Bad Request")

	// parsing the error response
	var errorResponse e.ResponseError
	err := parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrOrderRequestInvalid.Error(), errorResponse.Error, "error response should match the error content")
}

func TestUpdateOrder_Invalid_Request(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// get the router
	r := InitRouter()

	// create request body
	var takeOrderRequest requests.TakeOrderRequest
	takeOrderRequest.Status = "test"
	reqBody, err := createJson(takeOrderRequest)

	a.Nil(err, "should not have problem with create json")

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPatch, "/orders/1234", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusBadRequest, w.Code, "server should return back 400 Bad Request")

	// parsing the error response
	var errorResponse e.ResponseError
	err = parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrOrderRequestInvalid.Error(), errorResponse.Error, "error response should match the error content")
}

func TestUpdateOrder_Internal_Server_Error(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// set up expected result
	const orderId = 1
	order1 := models.Order{ID: orderId, Status:models.StatusTaken, Distance:rand.Intn(5000)}
	orders := []models.Order{order1}

	// get the router
	r := InitRouter()

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(orders)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that get the order by id
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "orders"  WHERE`).WithQueryException()
	defer mocket.Catcher.Reset()

	// create request body
	var takeOrderRequest requests.TakeOrderRequest
	takeOrderRequest.Status = models.StatusTaken
	reqBody, err := createJson(takeOrderRequest)

	a.Nil(err, "should not have problem with create json")

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ :=http.NewRequest(http.MethodPatch, "/orders/1", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusInternalServerError, w.Code, "server should return back 500 Internal Sever Error")

	// parsing the error response
	var errorResponse e.ResponseError
	err = parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrInternalError.Error(), errorResponse.Error, "error response should match the error content")
}
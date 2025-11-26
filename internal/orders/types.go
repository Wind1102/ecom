package orders

type orderItem struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type createOrderParams struct {
	CustomerId string      `json:"customerId"`
	Items      []orderItem `json:"items"`
}

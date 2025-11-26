package orders

type orderItem struct {
	ProductId int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerId int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}

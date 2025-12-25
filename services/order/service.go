package order

type Service interface {
	PostOrder()
	GetOrdersByAccountID()
}
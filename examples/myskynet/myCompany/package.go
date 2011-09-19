package myCompany

type SubscriptionRequest struct {
	Name string
	EmailAddress string
}

type SubscriptionResponse struct {
	Success bool
	Errors               []string
}

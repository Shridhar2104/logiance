package main
// resolver.go
type Resolver interface {
    Mutation() MutationResolver
    Query() QueryResolver
    Account() AccountResolver
    Order() OrderResolver
    Shipping() ShippingResolver  // Add this line
}
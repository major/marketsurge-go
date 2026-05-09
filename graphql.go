package marketsurge

// GraphQLRequest models a GraphQL request envelope.
type GraphQLRequest[V any] struct {
	OperationName string `json:"operationName,omitempty"`
	Variables     V      `json:"variables"`
	Query         string `json:"query"`
}

// GraphQLResponse models a GraphQL response envelope.
type GraphQLResponse[T any] struct {
	Data   T                   `json:"data"`
	Errors []GraphQLFieldError `json:"errors,omitempty"`
}

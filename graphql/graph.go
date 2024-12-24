// graphql/graph.go
package main

import (
    "github.com/99designs/gqlgen/graphql"
    "github.com/Shridhar2104/logilo/account"
    "github.com/Shridhar2104/logilo/shopify"
)

type Server struct {
    accountClient *account.Client
    shopifyClient *shopify.Client
}

func NewGraphQLServer(accountUrl, shopifyUrl string) (*Server, error) {
    accountClient, err := account.NewClient(accountUrl)
    if err != nil {
        return nil, err
    }

    shopifyClient, err := shopify.NewClient(shopifyUrl)    
    if err != nil {
        accountClient.Close()
        return nil, err
    }

    return &Server{accountClient, shopifyClient}, nil
}

func (s *Server) Mutation() MutationResolver {
    return &mutationResolver{s}
}

func (s *Server) Query() QueryResolver {
    return &queryResolver{s}
}

func (s *Server) Account() AccountResolver {
    return &accountResolver{s}
}

func (s *Server) Order() OrderResolver {
    return &orderResolver{s}
}

func (s *Server) ToNewExecutableSchema() graphql.ExecutableSchema {
    return NewExecutableSchema(Config{
        Resolvers: s,
    })
}
package models

import "github.com/graphql-go/graphql"

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"user": UserQuery,
	},
})

//var rootMutation = graphql.NewObject(graphql.ObjectConfig{
//	Name: "RootMutation",
//	Fields: graphql.Fields{
//		/*
//			curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:"My+new+todo"){id,text,done}}'
//		*/
//		"createUser": &graphql.Field{
//			Type:        userType, // the return type for this field
//			Description: "Create new user",
//			Args: graphql.FieldConfigArgument{
//				"name": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.String),
//				},
//			},
//			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//
//				// marshall and cast the argument value
//				text, _ := params.Args["name"].(string)
//
//				user := User{
//					ID:   1,
//					Name: text,
//				}
//
//				return user, nil
//			},
//		},
//	},
//})
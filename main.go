package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	Userc "xira/controllers/User" // Import the user controller
	"xira/db"
	"xira/dbinit"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to run the server: %v", err)
	}
}

func run() error {
	dbinit.DbInit()
	client := dbinit.Client
	ctx := context.Background()

	// Create a user using the user controller
	createUser, err := Userc.CreateUser(ctx, "yodevmw@gmail.com", "yossaf", "Av0129453c", []string{})
	if err != nil {
		fmt.Printf("failed to create user: %v\n", err)
	} else {
		result, _ := json.MarshalIndent(createUser, "", "  ")
		fmt.Printf("created user: %s\n", result)

		client.Task.CreateOne(
			db.Task.Title.Set("content"),
			db.Task.User.Link(
				db.User.ID.Equals(createUser.ID),
			),
		).Exec(ctx)

	}

	UserType, TaskType, ShopItemType := graphqlObjectCreator()
	MutationType := graphqlMutationCreator()

	var schemaConfig = graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type:        UserType,
					Description: "Get user by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["id"].(string)
						if ok {
							// Find user
							user, err := client.User.FindUnique(
								db.User.ID.Equals(id),
							).With(
								db.User.Tasks.Fetch(),
								db.User.Shopitems.Fetch(),
							).Exec(ctx)
							if err != nil {
								return nil, err
							}
							mp := bruteForceStrucToMap(user, "User")
							println("user GOT BY graphql: ", user.Name, user.Coins, user.Xp)
							return mp, err
						}
						return nil, nil
					},
				},
				"task": &graphql.Field{
					Type:        TaskType,
					Description: "Get task by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["id"].(string)
						if ok {
							// Find task
							task, err := client.Task.FindUnique(
								db.Task.ID.Equals(id),
							).Exec(ctx)
							if err != nil {
								return nil, err
							}
							mp := bruteForceStrucToMap(task, "Task")
							println("task GOT BY graphql: ", task.Title)
							return mp, err
						}
						return nil, nil
					},
				},
				"shopitem": &graphql.Field{
					Type:        ShopItemType,
					Description: "Get shopitem by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["id"].(string)
						if ok {
							// Find shopitem
							shopitem, err := client.ShopItem.FindUnique(
								db.ShopItem.ID.Equals(id),
							).Exec(ctx)
							if err != nil {
								return nil, err
							}
							mp := bruteForceStrucToMap(shopitem, "Shopitem")
							println("shopitem GOT BY graphql: ", shopitem.Name, shopitem.Price)
							return mp, err
						}
						return nil, nil
					},
				},
			},
		}),
		Mutation: MutationType,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	http.Handle("/graphql", h)
	println("Server started at http://localhost:8686/graphql")
	log.Fatal(http.ListenAndServe(":8686", nil))
	return nil
}

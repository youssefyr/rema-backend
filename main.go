package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"rema/db"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

// run is a function that initializes the database client, creates a user, and starts the GraphQL server.
// It returns an error if any operation fails.
// funcction  run () error

func run() error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	// create a user
	createUser, err := client.User.CreateOne(
		db.User.Name.Set("yossaf"),
		db.User.Coins.Set(0),
		db.User.Xp.Set(0),
	).Exec(ctx)
	if err != nil {
		return err
	}

	result, _ := json.MarshalIndent(createUser, "", "  ")
	fmt.Printf("created user: %s\n", result)

	client.Task.CreateOne(
		db.Task.Title.Set("content"),
		db.Task.User.Link(
			db.User.ID.Equals(createUser.ID),
		),
	).Exec(ctx)

	UserType, TaskType, ShopItemType, _ := graphqlObjectCreator()

	var querytype = graphql.NewObject(graphql.ObjectConfig{
		Name: "query",
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
	})

	schemaConfig := graphql.SchemaConfig{Query: querytype}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	//var query = `
	//	{
	//		hello
	//	}
	//`
	///params := graphql.Params{Schema: schema, RequestString: query}
	//r := graphql.Do(params)
	//if len(r.Errors) > 0 {
	//	log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	//}
	//rJSON, _ := json.Marshal(r)
	//fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}

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

//tasks, _ := client.Task.FindMany( // Get the user's tasks
//	db.Task.User.Where(
//		db.User.ID.Equals(user.ID),
//	),
//).Exec(ctx)

//	if tasks != nil {
//		tasksStr, _ := json.MarshalIndent(user, "", "  ")
//		fmt.Printf("The user's tasks are: %s\n", tasksStr)
//	} else {
//
// fmt.Println("The user's tasks are null")
// }

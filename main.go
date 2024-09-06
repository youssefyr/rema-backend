package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	// adapt "demo" to your module name if it differs
	"rema/db"

	"net/http"

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

func bruteForceStrucToMap(s interface{}, inputType string) map[string]interface{} {
	switch inputType {

	case "User":
		user := s.(*db.UserModel)
		return map[string]interface{}{
			"id":        user.ID,
			"name":      user.Name,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
			"coins":     user.Coins,
			"xp":        user.Xp,
			"tasks":     user.Tasks,
			"shopItems": user.Shopitems,
		}
	case "Task":
		task := s.(*db.TaskModel)
		return map[string]interface{}{
			"id":             task.ID,
			"name":           task.Title,
			"createdAt":      task.CreatedAt,
			"updatedAt":      task.UpdatedAt,
			"title":          task.Title,
			"description":    task.Description,
			"completed":      task.Completed,
			"user":           task.User,
			"repeatable":     task.Repeatable,
			"repeatduration": task.Repeatduration,
		}
	case "Shopitem":
		shopItem := s.(*db.ShopItemModel)
		return map[string]interface{}{
			"id":        shopItem.ID,
			"name":      shopItem.Name,
			"createdAt": shopItem.CreatedAt,
			"updatedAt": shopItem.UpdatedAt,
			"price":     shopItem.Price,
			"quantity":  shopItem.Quantity,
			"user":      shopItem.User,
		}
	default:
		return nil
	}
}

func graphqlObjectCreator() (userType, taskType, shopItemType, baseUserType *graphql.Object) {

	baseUserType = graphql.NewObject(graphql.ObjectConfig{
		Name: "BaseUser",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"coins": &graphql.Field{
				Type: graphql.Int,
			},
			"xp": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	taskType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Task",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"completed": &graphql.Field{
				Type: graphql.Boolean,
			},
			"user": &graphql.Field{
				Type: baseUserType,
			},
			"repeatable": &graphql.Field{
				Type: graphql.Boolean,
			},
			"repeatduration": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	shopItemType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Shopitem",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
			"quantity": &graphql.Field{
				Type: graphql.Int,
			},
			"user": &graphql.Field{
				Type: baseUserType,
			},
		},
	})

	userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"coins": &graphql.Field{
				Type: graphql.Int,
			},
			"xp": &graphql.Field{
				Type: graphql.Int,
			},
			"tasks": &graphql.Field{
				Type: taskType,
			},
			"shopItems": &graphql.Field{
				Type: shopItemType,
			},
		},
	})

	return userType, taskType, shopItemType, baseUserType
}

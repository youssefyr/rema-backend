package main

import (
	"context"
	"fmt"
	Userc "xira/controllers/User"
	"xira/db"
	"xira/dbinit"

	"github.com/graphql-go/graphql"
)

// Define the baseUser type
var BaseUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "BaseUser",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"coins": &graphql.Field{
			Type: graphql.Int,
		},
		"xp": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

// Define other types
var TaskType = graphql.NewObject(graphql.ObjectConfig{
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
			Type: BaseUserType,
		},
		"repeatable": &graphql.Field{
			Type: graphql.Boolean,
		},
		"repeatduration": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var ShopItemType = graphql.NewObject(graphql.ObjectConfig{
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
			Type: BaseUserType,
		},
	},
})

func resolveUserTasks(p graphql.ResolveParams) (interface{}, error) {
	userMap, ok := p.Source.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map[string]interface{} but got %T", p.Source)
	}

	userID, ok := userMap["id"].(string)
	if !ok {
		return nil, fmt.Errorf("expected user ID to be a string but got %T", userMap["id"])
	}

	client := dbinit.Client
	ctx := context.Background()

	tasks, err := client.Task.FindMany(
		db.Task.UserID.Equals(userID),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func resolveUserShopItems(p graphql.ResolveParams) (interface{}, error) {
	userMap, ok := p.Source.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map[string]interface{} but got %T", p.Source)
	}

	userID, ok := userMap["id"].(string)
	if !ok {
		return nil, fmt.Errorf("expected user ID to be a string but got %T", userMap["id"])
	}

	client := dbinit.Client
	ctx := context.Background()

	shopItems, err := client.ShopItem.FindMany(
		db.ShopItem.UserID.Equals(userID),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return shopItems, nil
}

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
		"sessionID": &graphql.Field{
			Type: graphql.String,
		},
		"rememberMe": &graphql.Field{
			Type: graphql.NewList(graphql.String),
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
			Type:    graphql.NewList(TaskType),
			Resolve: resolveUserTasks,
		},
		"shopItems": &graphql.Field{
			Type:    graphql.NewList(ShopItemType),
			Resolve: resolveUserShopItems,
		},
	},
})

// graphqlObjectCreator is a function that creates the GraphQL objects for the user, task, and shopitem types.
func graphqlObjectCreator() (*graphql.Object, *graphql.Object, *graphql.Object) {
	return UserType, TaskType, ShopItemType
}

// Add the mutation type to your GraphQL schema
func graphqlMutationCreator() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "CreateUserPayload",
					Fields: graphql.Fields{
						"user": &graphql.Field{
							Type: UserType,
						},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"rememberMe": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email := p.Args["email"].(string)
					name := p.Args["name"].(string)
					password := p.Args["password"].(string)
					rememberMeInterface := p.Args["rememberMe"].([]interface{})

					// Convert []interface{} to []string
					rememberMe := make([]string, len(rememberMeInterface))
					for i, v := range rememberMeInterface {
						rememberMe[i] = v.(string)
					}

					ctx := context.Background()
					user, err := Userc.CreateUser(ctx, email, name, password, rememberMe)
					if err != nil {
						return nil, err
					}

					fmt.Println("Created user:", user)

					userMap := bruteForceStrucToMap(user, "User")
					fmt.Println("User map:", userMap)

					return map[string]interface{}{
						"user": userMap,
					}, nil
				},
			},
		},
	})
}

// bruteForceStrucToMap is a function that converts a struct to a map. It is used to convert the user, task, and shopitem structs to maps.
func bruteForceStrucToMap(s interface{}, inputType string) map[string]interface{} {
	switch inputType {

	case "User":
		user := s.(*db.UserModel)
		return map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"password":   user.Password,
			"sessionID":  user.SessionID,
			"rememberMe": user.RememberMe,
			"createdAt":  user.CreatedAt,
			"updatedAt":  user.UpdatedAt,
			"coins":      user.Coins,
			"xp":         user.Xp,
			"tasks":      user.Tasks,
			"shopItems":  user.Shopitems,
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

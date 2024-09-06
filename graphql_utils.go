package main

import (
	"rema/db"

	"github.com/graphql-go/graphql"
)

// graphqlObjectCreator is a function that creates the GraphQL objects for the user, task, and shopitem types.
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

// bruteForceStrucToMap is a function that converts a struct to a map. It is used to convert the user, task, and shopitem structs to maps.
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

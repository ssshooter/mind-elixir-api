package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SSShooter/mind-elixir-backend-go/models"
)

// @Summary getAllPrivateMaps
// @Schemes
// @Description getAllPrivateMaps
// @Tags map
// @Router /api/map [get]
func getAllPrivateMaps(mapColl *mongo.Collection) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		loginId := c.MustGet("loginId").(int)
		cursor, err := mapColl.Find(
			context.TODO(),
			bson.D{{"author", loginId}},
		)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, gin.H{"data": results})
	}
}

// @Summary createPrivateMap
// @Schemes
// @Description createPrivateMap
// @Tags map
// @Router /api/map [post]
func createPrivateMap(mapColl *mongo.Collection) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var mapData *models.Map
		c.ShouldBind(&mapData)
		loginId := c.MustGet("loginId").(int)
		mapData.Author = loginId
		res, err := mapColl.InsertOne(context.TODO(), mapData)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": gin.H{"_id": res.InsertedID}})
	}
}

// @Summary getPrivateMap
// @Schemes
// @Description getPrivateMap
// @Tags map
// @Param id path string true "Map ID"
// @Router /api/map/{id} [get]
func getPrivateMap(mapColl *mongo.Collection) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		loginId := c.MustGet("loginId").(int)
		var result bson.M
		err := mapColl.FindOne(
			context.TODO(),
			bson.D{{"_id", id}, {"author", loginId}},
		).Decode(&result)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": result})
	}
}

// @Summary updatePrivateMap
// @Schemes
// @Description updatePrivateMap
// @Tags map
// @Param id path string true "Map ID"
// @Router /api/map/{id} [patch]
func updatePrivateMap(mapColl *mongo.Collection) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		var data map[string]interface{}
		c.ShouldBind(&data)
		loginId := c.MustGet("loginId").(int)
		var result bson.M
		update := bson.D{{"$set", data}}
		err := mapColl.FindOneAndUpdate(
			context.TODO(),
			bson.D{{"_id", id}, {"author", loginId}},
			update,
		).Decode(&result)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": result})
	}
}

// @Summary deletePrivateMap
// @Schemes
// @Description deletePrivateMap
// @Tags map
// @Param id path string true "Map ID"
// @Router /api/map [delete]
func deletePrivateMap(mapColl *mongo.Collection) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		loginId := c.MustGet("loginId").(int)
		var result bson.M
		err := mapColl.FindOneAndDelete(
			context.TODO(),
			bson.D{{"_id", id}, {"author", loginId}},
		).Decode(&result)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": result})
	}
}

func AddMapRoutes(rg *gin.RouterGroup, mapColl *mongo.Collection) {
	rg.GET("", getAllPrivateMaps(mapColl))
	rg.POST("", createPrivateMap(mapColl))
	rg.GET("/:id", getPrivateMap(mapColl))
	rg.PATCH("/:id", updatePrivateMap(mapColl))
	rg.DELETE("/:id", deletePrivateMap(mapColl))
}

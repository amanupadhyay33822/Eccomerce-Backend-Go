package utils

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUserID extracts user_id from context as ObjectID
func GetUserID(c *gin.Context) primitive.ObjectID {
	return c.MustGet("user_id").(primitive.ObjectID)
}

// GetParamObjectID converts URL parameter to ObjectID
func GetParamObjectID(c *gin.Context, param string) (primitive.ObjectID, error) {
	idStr := c.Param(param)
	return primitive.ObjectIDFromHex(idStr)
}

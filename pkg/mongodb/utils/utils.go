package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

// ConvertStringToObjectID convert string to objectID.
func ConvertStringToObjectID(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)

	return objID
}

// ConvertStringsToObjectID convert list strings to objectIDs.
func ConvertStringsToObjectID(ids ...string) []primitive.ObjectID {
	objectIDs := make([]primitive.ObjectID, 0)
	for _, id := range ids {
		objID, _ := primitive.ObjectIDFromHex(id)
		objectIDs = append(objectIDs, objID)
	}

	return objectIDs
}

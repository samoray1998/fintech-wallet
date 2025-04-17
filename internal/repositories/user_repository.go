package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/samoray1998/fintech-wallet/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database, collectionName string) *UserRepository {
	return &UserRepository{
		collection: db.Collection(collectionName),
	}
}

// / CreateUser with new hash password
func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	/// Hash password before storing

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.ID = primitive.NewObjectID()
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.KYCStatus = "unverified"

	_, err = r.collection.InsertOne(context.Background(), user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

///findBy id

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email (for authentication)

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

///  UpdateKYCStatus func

func (r *UserRepository) UpdateKYCStatus(id string, status string) (*models.User, error) {

	objctId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	var user *models.User

	err = r.collection.FindOne(context.Background(), bson.M{"_id": objctId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	validStatuses := map[string]bool{
		"unverified": true,
		"pending":    true,
		"verified":   true,
		"rejected":   true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid KYC status")
	}
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objctId}, bson.M{"$set": bson.M{"kyc_status": status, "updated_at": time.Now().Unix()}})
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objctId}).Decode(&user)

	return user, err

}

/// update password

func (r *UserRepository) UpdateUserPassword(userId string, newPass string) error {
	objectId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return errors.New("invalid user ID")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": bson.M{"password": hashedPassword, "updatedAt": time.Now()}})

	if err != nil {
		return err
	}

	return nil

}

//ListUsersWithKYCStatus

func (r *UserRepository) ListUsersWithKYCStatus(status string, page int, limit int) ([]models.User, error) {

	var users []models.User

	filter := bson.M{}

	if status != "" {
		filter["kyc_status"] = status
	}

	opts := options.Find().SetSkip(int64(page - 1)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

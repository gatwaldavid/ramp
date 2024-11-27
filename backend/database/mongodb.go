package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client
var db *mongo.Database

type Patient struct {
    ID        string    `bson:"_id,omitempty"`
    FirstName string    `bson:"firstName"`
    LastName  string    `bson:"lastName"`
    DOB       string    `bson:"dob"`
    Gender    string    `bson:"gender"`
    CreatedAt time.Time `bson:"createdAt"`
    UpdatedAt time.Time `bson:"updatedAt"`
}

type User struct {
    ID        string    `bson:"_id,omitempty"`
    Username  string    `bson:"username"`
    Password  string    `bson:"password"`
    Role      string    `bson:"role"`
    CreatedAt time.Time `bson:"createdAt"`
    UpdatedAt time.Time `bson:"updatedAt"`
}

func InitDatabase() error {
    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    // Connect to MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        return err
    }

    // Check the connection
    err = client.Ping(ctx, nil)
    if err != nil {
        return err
    }

    db = client.Database("hospital")

    // Initialize collections and indexes
    err = initializeCollections()
    if err != nil {
        return err
    }

    log.Println("Connected to MongoDB!")
    return nil
}

func initializeCollections() error {
    ctx := context.Background()

    // Create unique index for username in users collection
    usersCollection := db.Collection("users")
    indexModel := mongo.IndexModel{
        Keys:    bson.D{{Key: "username", Value: 1}},
        Options: options.Index().SetUnique(true),
    }
    _, err := usersCollection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        return err
    }

    // Insert admin user if not exists
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    adminUser := User{
        Username:  "admin",
        Password:  string(hashedPassword),
        Role:      "admin",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Try to insert admin user if not exists
    _, err = usersCollection.UpdateOne(
        ctx,
        bson.M{"username": "admin"},
        bson.M{"$setOnInsert": adminUser},
        options.Update().SetUpsert(true),
    )
    if err != nil {
        return err
    }

    // Insert sample patients if collection is empty
    patientsCollection := db.Collection("patients")
    count, err := patientsCollection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return err
    }

    if count == 0 {
        samplePatients := []interface{}{
            Patient{
                FirstName: "John",
                LastName:  "Doe",
                DOB:       "1990-01-01",
                Gender:    "Male",
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            },
            Patient{
                FirstName: "Jane",
                LastName:  "Doe",
                DOB:       "1995-02-01",
                Gender:    "Female",
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            },
        }

        _, err = patientsCollection.InsertMany(ctx, samplePatients)
        if err != nil {
            return err
        }
    }

    return nil
}

// Database operations for Users
func CreateUser(user User) error {
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    _, err := db.Collection("users").InsertOne(context.Background(), user)
    return err
}

func GetUserByUsername(username string) (User, error) {
    var user User
    err := db.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
    return user, err
}

// Database operations for Patients
func CreatePatient(patient Patient) error {
    patient.CreatedAt = time.Now()
    patient.UpdatedAt = time.Now()

    _, err := db.Collection("patients").InsertOne(context.Background(), patient)
    return err
}

func GetAllPatients() ([]Patient, error) {
    var patients []Patient
    cursor, err := db.Collection("patients").Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    err = cursor.All(context.Background(), &patients)
    return patients, err
}

func GetPatientByID(id string) (Patient, error) {
    var patient Patient
    err := db.Collection("patients").FindOne(context.Background(), bson.M{"_id": id}).Decode(&patient)
    return patient, err
}

func UpdatePatient(id string, patient Patient) error {
    patient.UpdatedAt = time.Now()

    _, err := db.Collection("patients").UpdateOne(
        context.Background(),
        bson.M{"_id": id},
        bson.M{"$set": patient},
    )
    return err
}

func DeletePatient(id string) error {
    _, err := db.Collection("patients").DeleteOne(context.Background(), bson.M{"_id": id})
    return err
}

// Helper functions
func GetDB() *mongo.Database {
    return db
}

func CloseDB() error {
    return client.Disconnect(context.Background())
}

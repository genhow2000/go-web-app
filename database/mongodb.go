package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client
var MongoDBDatabase *mongo.Database

// InitMongoDB 初始化MongoDB连接
func InitMongoDB() error {
	// 从环境变量获取MongoDB连接信息
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		// 如果没有设置MongoDB URI，跳过MongoDB初始化
		log.Println("MONGODB_URI not set, skipping MongoDB initialization")
		return nil
	}

	log.Printf("Attempting to connect to MongoDB with URI: %s", mongoURI[:50]+"...")

	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "chatbot" // 默认数据库名
	}

	// 设置连接超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	log.Printf("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("MongoDB connection failed: %v", err)
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	log.Printf("Testing MongoDB connection...")
	if err = client.Ping(ctx, nil); err != nil {
		log.Printf("MongoDB ping failed: %v", err)
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoDBClient = client
	MongoDBDatabase = client.Database(databaseName)

	log.Println("MongoDB connected successfully!")
	
	// 创建索引
	if err := createMongoIndexes(); err != nil {
		log.Printf("Warning: failed to create MongoDB indexes: %v", err)
	}

	return nil
}

// createMongoIndexes 创建MongoDB索引
func createMongoIndexes() error {
	ctx := context.Background()
	
	// 为conversations集合创建索引
	conversationsCollection := MongoDBDatabase.Collection("conversations")
	
	// 用户ID索引
	_, err := conversationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"user_id": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create user_id index: %w", err)
	}

	// 创建时间索引
	_, err = conversationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"created_at": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create created_at index: %w", err)
	}

	// 更新时间索引
	_, err = conversationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"updated_at": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create updated_at index: %w", err)
	}

	// 活跃状态索引
	_, err = conversationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"is_active": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create is_active index: %w", err)
	}

	// 复合索引：用户ID + 活跃状态
	_, err = conversationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"user_id": 1,
			"is_active": 1,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create user_id+is_active index: %w", err)
	}

	log.Println("MongoDB indexes created successfully!")
	return nil
}

// CloseMongoDB 关闭MongoDB连接
func CloseMongoDB() error {
	if MongoDBClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return MongoDBClient.Disconnect(ctx)
	}
	return nil
}

// GetMongoDBCollection 获取MongoDB集合
func GetMongoDBCollection(name string) *mongo.Collection {
	if MongoDBDatabase == nil {
		return nil
	}
	return MongoDBDatabase.Collection(name)
}

// IsMongoDBConnected 检查MongoDB是否已连接
func IsMongoDBConnected() bool {
	if MongoDBClient == nil {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	err := MongoDBClient.Ping(ctx, nil)
	return err == nil
}

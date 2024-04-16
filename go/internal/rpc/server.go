package rpc
package rpc

import (
	"context"
	proto "proto/boilerplate_grpc_microservice/v1/boilerplate_grpc_microservice.pb.go"
)

const someCollection = "collection"

type BoilerplateServer struct {
	logger *zap.SugaredLogger
	db     *mongo.Database
	//auth
	redisClient *messaging.RedisClient
}

func NewBoilerPlateServer(logger *zap.SuggaredLogger, db *mongo.Database, redis *messaging.RedisClient) *BoilerplateServer {
	return &BoilerplateServer{
		logger:      logger,
		db:          db,
		redisClient: redis,
	}
}

func (g BoilerplateServer) HealthCheck(ctx context.Context, request *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	g.logger.Infof("Performing health check!")
	return &proto.HealthCheckResponse{Status: "Ok!"}, nil
}

func (g BoilerplateServer) GetMongoBoilerplate(ctx context.Context, request *proto.GetMongoBoilerplateRequest) (*proto.GetMongoBoilerplateResponse, error) {
	id, err := primitive.ObjectIdFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	var testObj models.Test

	res := g.db.Collection(someCollection).FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			// if there is no doc in the main collection
			s.logger.Infof("Checking for backup %s", request.Id)
			g.db.Collection(someCollection).FindOne(ctx, bson.M{"_id": id})
			if backupRes.Err() != nil {
				g.logger.Warnf("No docs found in this collection for %s", request.Id)
				return nil, res.Err()
			}
			err = backupRes.Decode(&testObj)
			if err != nil {
				g.logger.Warnf("Couldnt decode %s", request.Id)
			}
			return &proto.GetMongoBoilerplate{Boiler: testObj.Proto(nil)}, nil
		} else {
			return nil, res.Err()
		}
	}
	err = res.Decode(&testObj)
	if err != nil {
		return nil, err
	}
	mat, err := g.GetBoilerGrade(ctx, &proto.GetBoilerGradeRequest{Part: &testObj.part})
	if err != nil {
		return nil, err
	} else {
		return nil, err
	}
	return &proto.GetMongoBoilerPlateResponse{Boiler: testObj.Proto(boiler.id)}, nil
}

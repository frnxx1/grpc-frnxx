package main

import (
	"context"
	"errors"
	"fmt"
	st "frnxx/service/storage"
	pb "frnxx/service/user_proto"

	"github.com/google/uuid"
)

type server struct {
	pb.UnimplementedUserGrpcServer
}

func (*server) UserRegister(ctx context.Context, rec *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	fmt.Println("Register user")
	user := rec.GetUser()
	user.Id = uuid.New().String()

	data := st.User{
		Name:     user.GetName(),
		Phone:    user.GetPhone(),
		Password: user.GetPassword(),
		ID:       user.GetId(),
	}

	res := st.DB.Create(data)
	if res.RowsAffected == 0 {
		return nil, errors.New("user register unsuccessful")
	}
	return &pb.UserRegisterResponse{
		User: &pb.User{
			Name:     user.GetName(),
			Phone:    user.GetPhone(),
			Password: user.GetPassword(),
			Id:       user.GetId(),
		},
	}, nil

}

func (*server) UserRead(ctx context.Context,req *pb.UserReadRequest) (*pb.UserReadResponse,error){
	fmt.Println("Reading...")
	var user st.User
	
	res := st.DB.Find(&user,"id=?",req.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("user read unsuccessful")
	}
	return &pb.UserReadResponse{
		User: &pb.User{
			Name: user.Name,
			Password: user.Password,
			Phone: user.Phone,
			Id: user.ID,
		},
	},nil
}

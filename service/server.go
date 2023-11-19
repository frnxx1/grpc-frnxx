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
			Id:       user.GetId(),
			Name:     user.GetName(),
			Phone:    user.GetPhone(),
			Password: user.GetPassword(),
		},
	}, nil

}

func (*server) UserRead(ctx context.Context, req *pb.UserReadRequest) (*pb.UserReadResponse, error) {
	fmt.Println("Reading...")
	var user st.User

	res := st.DB.Find(&user, "id=?", req.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("user read unsuccessful")
	}
	return &pb.UserReadResponse{
		User: &pb.User{
			Id:       user.ID,
			Name:     user.Name,
			Phone:    user.Phone,
			Password: user.Password,
		},
	}, nil
}

func (*server) UserUpdate(ctx context.Context, req *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	fmt.Println("Updating phone...")
	var user st.User

	res := st.DB.Model(&user).Where("id = ?", req.GetUser().Id).Updates(st.User{Phone: req.User.Phone, Password: req.User.Password,Name: req.User.Name})
	if res.RowsAffected == 0 {
		return nil, errors.New("user update unsuccessful")
	}

	return &pb.UserUpdateResponse{
		User: &pb.User{
			Id:       user.ID,
			Name:     user.Name,
			Phone:    user.Phone,
			Password: user.Password,
		},
	}, nil
}




package main

import (
	"context"
	"flag"
	"fmt"
	pb "frnxx/client/user_proto"
	st "frnxx/service/storage"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewUserGrpcClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	r := GetRouter()

	r.POST("/UserCreate", func(ct *gin.Context) {
		var user st.User
		err := ct.ShouldBind(&user)
		if err != nil {
			ct.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		userData := &pb.User{
			Name:     user.Name,
			Phone:    user.Phone,
			Password: user.Password,
		}
		res, err := c.UserRegister(ctx, &pb.UserRegisterRequest{User: userData})

		if err != nil {
			ct.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		ct.JSON(http.StatusOK, gin.H{
			"User": res.User,
		})

	})

	r.GET("/User/:id", func(ct *gin.Context) {
		id := ct.Param("id")

		res, err := c.UserRead(ctx, &pb.UserReadRequest{Id: id})
		if err != nil {
			ct.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ct.JSON(http.StatusOK, gin.H{
			"User": res.User,
		})
	})

	r.PUT("/User/:id", func(ctx *gin.Context) {
		var user st.User
		err := ctx.ShouldBind(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error1": err,
			})
			return
		}

		res, err := c.UserUpdate(ctx, &pb.UserUpdateRequest{
			User: &pb.User{
				Id:       user.ID,
				Name:     user.Name,
				Phone:    user.Phone,
				Password: user.Password,
			},
		})
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error2": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Update": res.User,
		})

	})

	r.Run("localhost:8080")

}

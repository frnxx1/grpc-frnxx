package main


import (
	"context"
	pb "frnxx/client/user_proto"
	st "frnxx/service/storage"
	"flag"
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
			ct.JSON(http.StatusBadRequest,gin.H{
				"error": err,
			})
			return
		}

		userData := &pb.User{
			Name: user.Name,
			Password: user.Password,
			Phone: user.Phone,
		}
		res,err := c.UserRegister(ctx,&pb.UserRegisterRequest{User: userData})

		if err != nil {
			ct.JSON(http.StatusBadRequest,gin.H{
				"error": err,
			})
			return
		}

		ct.JSON(http.StatusOK,gin.H{
			"User":res.User,
		})

	})

	r.GET("/User/:id",func(ct *gin.Context) {
		id := ct.Param("id")

		res,err := c.UserRead(ctx,&pb.UserReadRequest{Id: id })
		if err != nil {
			ct.JSON(http.StatusBadRequest,gin.H{
				"error": err,
			})
			return
		}
		ct.JSON(http.StatusOK,gin.H{
			"User":res.User,
		})
	})

	r.Run("localhost:8080")

}
package main

import (
	"context"
	"errors"
	"time"

	// "net/http"

	"connectrpc.com/connect"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	calculatorv1 "github.com/yinhuaz/ex_project/gen"
	"github.com/yinhuaz/ex_project/gen/calculatorv1connect"
)

type CalculatorServer struct{}

func (s *CalculatorServer) Add(
	ctx context.Context,
	req *connect.Request[calculatorv1.AddRequest],
) (*connect.Response[calculatorv1.AddResponse], error) {
	result := req.Msg.A + req.Msg.B
	return connect.NewResponse(&calculatorv1.AddResponse{Result: result}), nil
}

func (s *CalculatorServer) Subtract(
	ctx context.Context,
	req *connect.Request[calculatorv1.SubtractRequest],
) (*connect.Response[calculatorv1.SubtractResponse], error) {
	result := req.Msg.A - req.Msg.B
	return connect.NewResponse(&calculatorv1.SubtractResponse{Result: result}), nil
}

func (s *CalculatorServer) Multiply(
	ctx context.Context,
	req *connect.Request[calculatorv1.MultiplyRequest],
) (*connect.Response[calculatorv1.MultiplyResponse], error) {
	result := req.Msg.A * req.Msg.B
	return connect.NewResponse(&calculatorv1.MultiplyResponse{Result: result}), nil
}

func (s *CalculatorServer) Divide(
	ctx context.Context,
	req *connect.Request[calculatorv1.DivideRequest],
) (*connect.Response[calculatorv1.DivideResponse], error) {
	if req.Msg.B == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("除数不能为零"))
	}
	result := req.Msg.A / req.Msg.B
	return connect.NewResponse(&calculatorv1.DivideResponse{Result: result}), nil
}

func main() {
	router := gin.Default()

	// 添加CORS中间件配置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	_, calculatorHandler := calculatorv1connect.NewCalculatorServiceHandler(&CalculatorServer{})
	router.POST(calculatorv1connect.CalculatorServiceAddProcedure, gin.WrapH(calculatorHandler))
	router.POST(calculatorv1connect.CalculatorServiceSubtractProcedure, gin.WrapH(calculatorHandler))
	router.POST(calculatorv1connect.CalculatorServiceMultiplyProcedure, gin.WrapH(calculatorHandler))
	router.POST(calculatorv1connect.CalculatorServiceDivideProcedure, gin.WrapH(calculatorHandler))

	router.Run("localhost:8080")
}

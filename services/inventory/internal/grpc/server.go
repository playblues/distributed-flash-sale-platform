package grpc

import (
	"context"
	"fmt"

	"inventory-service/internal/repository"
	inventoryv1 "inventory-service/pkg/api/inventory/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryGrpcServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	repo *repository.InventoryRepository
}

func NewInventoryGrpcServer(repo *repository.InventoryRepository) *InventoryGrpcServer {
	return &InventoryGrpcServer{repo: repo}
}

func (s *InventoryGrpcServer) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	if req.EventId <= 0 || req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "event_id and user_id must be greater than zero")
	}

	eventIDStr := fmt.Sprintf("%d", req.EventId)
	userIDStr := fmt.Sprintf("%d", req.UserId)

	result, err := s.repo.ReserveStock(ctx, eventIDStr, userIDStr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	switch result {
	case -1:
		return &inventoryv1.ReserveStockResponse{
			Status:  inventoryv1.ReservationStatus_RESERVATION_STATUS_ALREADY_RESERVED,
			Message: "user already has an active reservation for this event",
		}, nil
	case 0:
		return &inventoryv1.ReserveStockResponse{
			Status:  inventoryv1.ReservationStatus_RESERVATION_STATUS_SOLD_OUT,
			Message: "tickets are sold out",
		}, nil
	default:
		return &inventoryv1.ReserveStockResponse{
			Status:  inventoryv1.ReservationStatus_RESERVATION_STATUS_RESERVED,
			Message: fmt.Sprintf("ticket successfully reserved for user %s", userIDStr),
		}, nil
	}
}

func (s *InventoryGrpcServer) RollbackStock(ctx context.Context, req *inventoryv1.RollbackStockRequest) (*inventoryv1.RollbackStockResponse, error) {
	if req.EventId <= 0 || req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "event_id and user_id must be greater than zero")
	}

	eventIDStr := fmt.Sprintf("%d", req.EventId)
	userIDStr := fmt.Sprintf("%d", req.UserId)

	if err := s.repo.RollBack(ctx, eventIDStr, userIDStr); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to rollback reservation: %v", err)
	}

	return &inventoryv1.RollbackStockResponse{
		Success: true,
		Message: fmt.Sprintf("reservation successfully rolled back for user %s", userIDStr),
	}, nil
}

func (s *InventoryGrpcServer) GetStock(ctx context.Context, req *inventoryv1.GetStockRequest) (*inventoryv1.GetStockResponse, error) {
	if req.EventId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "event_id must be greater than zero")
	}

	eventIDStr := fmt.Sprintf("%d", req.EventId)
	stock, err := s.repo.GetStock(ctx, eventIDStr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch stock: %v", err)
	}

	return &inventoryv1.GetStockResponse{
		EventId: req.EventId,
		Stock:   int64(stock),
	}, nil
}

func (s *InventoryGrpcServer) InitStock(ctx context.Context, req *inventoryv1.InitStockRequest) (*inventoryv1.InitStockResponse, error) {
	if req.EventId <= 0 || req.Amount < 0 {
		return nil, status.Error(codes.InvalidArgument, "event_id and amount must be valid")
	}

	eventIDStr := fmt.Sprintf("%d", req.EventId)
	if err := s.repo.InitializeStock(ctx, eventIDStr, int(req.Amount)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to initialize stock: %v", err)
	}

	return &inventoryv1.InitStockResponse{
		Success: true,
		Message: "stock successfully initialized",
	}, nil
}

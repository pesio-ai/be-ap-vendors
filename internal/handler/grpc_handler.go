package handler

import (
	"context"

	"github.com/pesio-ai/be-lib-common/auth"
	"github.com/pesio-ai/be-lib-common/logger"
	commonpb "github.com/pesio-ai/be-lib-proto/gen/go/common"
	pb "github.com/pesio-ai/be-lib-proto/gen/go/ap"
	"github.com/pesio-ai/be-vendors-service/internal/repository"
	"github.com/pesio-ai/be-vendors-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler handles gRPC requests for vendors service
type GRPCHandler struct {
	pb.UnimplementedVendorsServiceServer
	vendorService *service.VendorService
	log           *logger.Logger
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(vendorService *service.VendorService, log *logger.Logger) *GRPCHandler {
	return &GRPCHandler{
		vendorService: vendorService,
		log:           log,
	}
}

// CreateVendor creates a new vendor
func (h *GRPCHandler) CreateVendor(ctx context.Context, req *pb.CreateVendorRequest) (*pb.Vendor, error) {
	// Extract user context from authenticated request
	userCtx, err := auth.GetUserContext(ctx)
	if err != nil {
		h.log.Warn().Err(err).Msg("User context not found")
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	h.log.Info().
		Str("entity_id", req.EntityId).
		Str("vendor_code", req.VendorCode).
		Str("vendor_name", req.VendorName).
		Str("user_id", userCtx.UserID).
		Msg("gRPC CreateVendor request")

	// Verify entity_id matches authenticated user's entity
	if req.EntityId != userCtx.EntityID {
		h.log.Warn().
			Str("req_entity_id", req.EntityId).
			Str("user_entity_id", userCtx.EntityID).
			Msg("Entity ID mismatch")
		return nil, status.Error(codes.PermissionDenied, "access denied: entity mismatch")
	}

	svcReq := &service.CreateVendorRequest{
		EntityID:          req.EntityId,
		VendorCode:        req.VendorCode,
		VendorName:        req.VendorName,
		LegalName:         stringPtr(req.LegalName),
		VendorType:        req.VendorType,
		TaxID:             stringPtr(req.TaxId),
		IsTaxExempt:       req.IsTaxExempt,
		Is1099Vendor:      req.Is_1099Vendor,
		Email:             stringPtr(req.Email),
		Phone:             stringPtr(req.Phone),
		Fax:               stringPtr(req.Fax),
		Website:           stringPtr(req.Website),
		AddressLine1:      stringPtr(req.AddressLine1),
		AddressLine2:      stringPtr(req.AddressLine2),
		City:              stringPtr(req.City),
		StateProvince:     stringPtr(req.StateProvince),
		PostalCode:        stringPtr(req.PostalCode),
		Country:           req.Country,
		PaymentTerms:      req.PaymentTerms,
		PaymentMethod:     stringPtr(req.PaymentMethod),
		Currency:          req.Currency,
		CreditLimit:       int64Ptr(req.CreditLimit),
		BankName:          stringPtr(req.BankName),
		BankAccountNumber: stringPtr(req.BankAccountNumber),
		BankRoutingNumber: stringPtr(req.BankRoutingNumber),
		SwiftCode:         stringPtr(req.SwiftCode),
		IBAN:              stringPtr(req.Iban),
		Notes:             stringPtr(req.Notes),
		Tags:              req.Tags,
		CreatedBy:         userCtx.UserID, // Use authenticated user ID
	}

	vendor, err := h.vendorService.CreateVendor(ctx, svcReq)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to create vendor")
		return nil, toGRPCError(err)
	}

	return vendorToProto(vendor), nil
}

// GetVendor retrieves a vendor by ID
func (h *GRPCHandler) GetVendor(ctx context.Context, req *pb.GetVendorRequest) (*pb.Vendor, error) {
	h.log.Info().
		Str("id", req.Id).
		Str("entity_id", req.EntityId).
		Msg("gRPC GetVendor request")

	vendor, err := h.vendorService.GetVendor(ctx, req.Id, req.EntityId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get vendor")
		return nil, toGRPCError(err)
	}

	return vendorToProto(vendor), nil
}

// UpdateVendor updates a vendor
func (h *GRPCHandler) UpdateVendor(ctx context.Context, req *pb.UpdateVendorRequest) (*pb.Vendor, error) {
	// Extract user context from authenticated request
	userCtx, err := auth.GetUserContext(ctx)
	if err != nil {
		h.log.Warn().Err(err).Msg("User context not found")
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	h.log.Info().
		Str("id", req.Id).
		Str("entity_id", req.EntityId).
		Str("user_id", userCtx.UserID).
		Msg("gRPC UpdateVendor request")

	// Verify entity_id matches authenticated user's entity
	if req.EntityId != userCtx.EntityID {
		h.log.Warn().
			Str("req_entity_id", req.EntityId).
			Str("user_entity_id", userCtx.EntityID).
			Msg("Entity ID mismatch")
		return nil, status.Error(codes.PermissionDenied, "access denied: entity mismatch")
	}

	svcReq := &service.UpdateVendorRequest{
		ID:                req.Id,
		EntityID:          req.EntityId,
		VendorCode:        req.VendorCode,
		VendorName:        req.VendorName,
		LegalName:         stringPtr(req.LegalName),
		VendorType:        req.VendorType,
		Status:            req.Status,
		TaxID:             stringPtr(req.TaxId),
		IsTaxExempt:       req.IsTaxExempt,
		Is1099Vendor:      req.Is_1099Vendor,
		Email:             stringPtr(req.Email),
		Phone:             stringPtr(req.Phone),
		Fax:               stringPtr(req.Fax),
		Website:           stringPtr(req.Website),
		AddressLine1:      stringPtr(req.AddressLine1),
		AddressLine2:      stringPtr(req.AddressLine2),
		City:              stringPtr(req.City),
		StateProvince:     stringPtr(req.StateProvince),
		PostalCode:        stringPtr(req.PostalCode),
		Country:           req.Country,
		PaymentTerms:      req.PaymentTerms,
		PaymentMethod:     stringPtr(req.PaymentMethod),
		Currency:          req.Currency,
		CreditLimit:       int64Ptr(req.CreditLimit),
		BankName:          stringPtr(req.BankName),
		BankAccountNumber: stringPtr(req.BankAccountNumber),
		BankRoutingNumber: stringPtr(req.BankRoutingNumber),
		SwiftCode:         stringPtr(req.SwiftCode),
		IBAN:              stringPtr(req.Iban),
		Notes:             stringPtr(req.Notes),
		Tags:              req.Tags,
		UpdatedBy:         userCtx.UserID, // Use authenticated user ID
	}

	vendor, err := h.vendorService.UpdateVendor(ctx, svcReq)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to update vendor")
		return nil, toGRPCError(err)
	}

	return vendorToProto(vendor), nil
}

// DeleteVendor deletes a vendor
func (h *GRPCHandler) DeleteVendor(ctx context.Context, req *pb.DeleteVendorRequest) (*commonpb.Response, error) {
	h.log.Info().
		Str("id", req.Id).
		Str("entity_id", req.EntityId).
		Msg("gRPC DeleteVendor request")

	err := h.vendorService.DeleteVendor(ctx, req.Id, req.EntityId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to delete vendor")
		return nil, toGRPCError(err)
	}

	return &commonpb.Response{
		Success: true,
		Message: "Vendor deleted successfully",
	}, nil
}

// ListVendors lists vendors with filtering and pagination
func (h *GRPCHandler) ListVendors(ctx context.Context, req *pb.ListVendorsRequest) (*pb.ListVendorsResponse, error) {
	h.log.Info().
		Str("entity_id", req.EntityId).
		Int32("page", req.Page).
		Int32("page_size", req.PageSize).
		Msg("gRPC ListVendors request")

	var status *string
	if req.Status != "" {
		status = &req.Status
	}

	var vendorType *string
	if req.VendorType != "" {
		vendorType = &req.VendorType
	}

	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	vendors, total, err := h.vendorService.ListVendors(ctx, req.EntityId, status, vendorType, req.ActiveOnly, page, pageSize)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to list vendors")
		return nil, toGRPCError(err)
	}

	pbVendors := make([]*pb.Vendor, len(vendors))
	for i, vendor := range vendors {
		pbVendors[i] = vendorToProto(vendor)
	}

	return &pb.ListVendorsResponse{
		Vendors:  pbVendors,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// ActivateVendor activates a vendor
func (h *GRPCHandler) ActivateVendor(ctx context.Context, req *pb.ActivateVendorRequest) (*commonpb.Response, error) {
	// Extract user context from authenticated request
	userCtx, err := auth.GetUserContext(ctx)
	if err != nil {
		h.log.Warn().Err(err).Msg("User context not found")
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	h.log.Info().
		Str("id", req.Id).
		Str("entity_id", req.EntityId).
		Str("user_id", userCtx.UserID).
		Msg("gRPC ActivateVendor request")

	// Verify entity_id matches authenticated user's entity
	if req.EntityId != userCtx.EntityID {
		h.log.Warn().
			Str("req_entity_id", req.EntityId).
			Str("user_entity_id", userCtx.EntityID).
			Msg("Entity ID mismatch")
		return nil, status.Error(codes.PermissionDenied, "access denied: entity mismatch")
	}

	err = h.vendorService.ActivateVendor(ctx, req.Id, req.EntityId, userCtx.UserID)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to activate vendor")
		return nil, toGRPCError(err)
	}

	return &commonpb.Response{
		Success: true,
		Message: "Vendor activated successfully",
	}, nil
}

// DeactivateVendor deactivates a vendor
func (h *GRPCHandler) DeactivateVendor(ctx context.Context, req *pb.DeactivateVendorRequest) (*commonpb.Response, error) {
	// Extract user context from authenticated request
	userCtx, err := auth.GetUserContext(ctx)
	if err != nil {
		h.log.Warn().Err(err).Msg("User context not found")
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	h.log.Info().
		Str("id", req.Id).
		Str("entity_id", req.EntityId).
		Str("user_id", userCtx.UserID).
		Msg("gRPC DeactivateVendor request")

	// Verify entity_id matches authenticated user's entity
	if req.EntityId != userCtx.EntityID {
		h.log.Warn().
			Str("req_entity_id", req.EntityId).
			Str("user_entity_id", userCtx.EntityID).
			Msg("Entity ID mismatch")
		return nil, status.Error(codes.PermissionDenied, "access denied: entity mismatch")
	}

	err = h.vendorService.DeactivateVendor(ctx, req.Id, req.EntityId, userCtx.UserID)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to deactivate vendor")
		return nil, toGRPCError(err)
	}

	return &commonpb.Response{
		Success: true,
		Message: "Vendor deactivated successfully",
	}, nil
}

// ValidateVendor validates a vendor
func (h *GRPCHandler) ValidateVendor(ctx context.Context, req *pb.ValidateVendorRequest) (*pb.ValidateVendorResponse, error) {
	h.log.Info().
		Str("id", req.Id).
		Str("entity_id", req.EntityId).
		Msg("gRPC ValidateVendor request")

	valid, message, err := h.vendorService.ValidateVendor(ctx, req.Id, req.EntityId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to validate vendor")
		return nil, toGRPCError(err)
	}

	return &pb.ValidateVendorResponse{
		Valid:   valid,
		Message: message,
	}, nil
}

// UpdateBalance updates the vendor's current balance
func (h *GRPCHandler) UpdateBalance(ctx context.Context, req *pb.UpdateBalanceRequest) (*commonpb.Response, error) {
	h.log.Info().
		Str("id", req.Id).
		Str("entity_id", req.EntityId).
		Int64("amount", req.Amount).
		Msg("gRPC UpdateBalance request")

	err := h.vendorService.UpdateBalance(ctx, req.Id, req.EntityId, req.Amount)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to update vendor balance")
		return nil, toGRPCError(err)
	}

	return &commonpb.Response{
		Success: true,
		Message: "Vendor balance updated successfully",
	}, nil
}

// Helper functions

func vendorToProto(vendor *repository.Vendor) *pb.Vendor {
	return &pb.Vendor{
		Id:                vendor.ID,
		EntityId:          vendor.EntityID,
		VendorCode:        vendor.VendorCode,
		VendorName:        vendor.VendorName,
		LegalName:         stringToProto(vendor.LegalName),
		VendorType:        vendor.VendorType,
		Status:            vendor.Status,
		TaxId:             stringToProto(vendor.TaxID),
		IsTaxExempt:       vendor.IsTaxExempt,
		Is_1099Vendor:     vendor.Is1099Vendor,
		Email:             stringToProto(vendor.Email),
		Phone:             stringToProto(vendor.Phone),
		Fax:               stringToProto(vendor.Fax),
		Website:           stringToProto(vendor.Website),
		AddressLine1:      stringToProto(vendor.AddressLine1),
		AddressLine2:      stringToProto(vendor.AddressLine2),
		City:              stringToProto(vendor.City),
		StateProvince:     stringToProto(vendor.StateProvince),
		PostalCode:        stringToProto(vendor.PostalCode),
		Country:           vendor.Country,
		PaymentTerms:      vendor.PaymentTerms,
		PaymentMethod:     stringToProto(vendor.PaymentMethod),
		Currency:          vendor.Currency,
		CreditLimit:       int64ToProto(vendor.CreditLimit),
		CurrentBalance:    vendor.CurrentBalance,
		BankName:          stringToProto(vendor.BankName),
		BankAccountNumber: stringToProto(vendor.BankAccountNumber),
		BankRoutingNumber: stringToProto(vendor.BankRoutingNumber),
		SwiftCode:         stringToProto(vendor.SwiftCode),
		Iban:              stringToProto(vendor.IBAN),
		Notes:             stringToProto(vendor.Notes),
		Tags:              vendor.Tags,
		CreatedAt:         timestamppb.New(vendor.CreatedAt),
		UpdatedAt:         timestamppb.New(vendor.UpdatedAt),
	}
}

func stringToProto(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int64ToProto(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

func toGRPCError(err error) error {
	// TODO: Map common errors to gRPC status codes
	return status.Error(codes.Internal, err.Error())
}

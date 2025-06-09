package office

import (
	"context"

	pb "legoas/legoas/proto"
)

type OfficeServiceServer struct {
	pb.UnimplementedOfficeServiceServer
	Repo *OfficeRepository
}

func NewOfficeServiceServer(repo *OfficeRepository) *OfficeServiceServer {
	return &OfficeServiceServer{Repo: repo}
}

func (s *OfficeServiceServer) SearchOffices(ctx context.Context, req *pb.SearchOfficesRequest) (*pb.SearchOfficesResponse, error) {
	offices, err := s.Repo.SearchOffices(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}

	var result []*pb.Office
	for _, o := range offices {
		result = append(result, &pb.Office{
			OfficeCode: o.OfficeCode,
			OfficeName: o.OfficeName,
		})
	}

	return &pb.SearchOfficesResponse{Offices: result}, nil
}

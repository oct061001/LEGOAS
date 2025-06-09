package menu

import (
	"context"

	pb "legoas/legoas/proto"
)

type MenuServiceServer struct {
	pb.UnimplementedMenuServiceServer
	Repo *MenuRepository
}

func NewMenuServiceServer(repo *MenuRepository) *MenuServiceServer {
	return &MenuServiceServer{Repo: repo}
}

func (s *MenuServiceServer) SearchMenus(ctx context.Context, req *pb.SearchMenusRequest) (*pb.SearchMenusResponse, error) {
	menus, err := s.Repo.SearchMenus(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}

	var result []*pb.Menu
	for _, m := range menus {
		result = append(result, &pb.Menu{
			MenuCode: m.MenuCode,
			MenuName: m.MenuName,
		})
	}
	return &pb.SearchMenusResponse{Menus: result}, nil

}

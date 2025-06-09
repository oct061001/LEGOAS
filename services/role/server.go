package role

import (
	"context"
	pb "legoas/legoas/proto"
)

type RoleServiceServer struct {
	pb.UnimplementedRoleServiceServer
	Repo *RoleRepository
}

func NewRoleServiceServer(repo *RoleRepository) *RoleServiceServer {
	return &RoleServiceServer{Repo: repo}
}

func (s *RoleServiceServer) SearchRoles(ctx context.Context, req *pb.SearchRolesRequest) (*pb.SearchRolesResponse, error) {
	roles, err := s.Repo.SearchRoles(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}

	var result []*pb.Role
	for _, r := range roles {
		result = append(result, &pb.Role{
			RoleCode: r.RoleCode,
			RoleName: r.RoleName,
		})
	}

	return &pb.SearchRolesResponse{Roles: result}, nil
}

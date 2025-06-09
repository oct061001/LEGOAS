package account

import (
	pb "legoas/legoas/proto"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	Name       string `bson:"name"`
	Address    string `bson:"address"`
	PostalCode string `bson:"postal_code"`
	Province   string `bson:"province"`
	OfficeCode string `bson:"office_code"`
}

type AccessRight struct {
	MenuCode  string `bson:"menu_code"`
	MenuName  string `bson:"menu_name"`
	CanCreate bool   `bson:"can_create"`
	CanRead   bool   `bson:"can_read"`
	CanUpdate bool   `bson:"can_update"`
	CanDelete bool   `bson:"can_delete"`
}

type AccountDoc struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	AccountName  string             `bson:"account_name"`
	Password     string             `bson:"password"`
	UserInfo     UserInfo           `bson:"user_info"`
	RoleCodes    []string           `bson:"role_codes"`
	OfficeCodes  []string           `bson:"office_codes"`
	AccessRights []AccessRight      `bson:"access_rights"`
	CreatedAt    time.Time          `bson:"created_at"`
}

func pbUserInfoToBSON(u *pb.UserInfo) UserInfo {
	if u == nil {
		return UserInfo{}
	}
	return UserInfo{
		Name:       u.GetName(),
		Address:    u.GetAddress(),
		PostalCode: u.GetPostalCode(),
		Province:   u.GetProvince(),
		OfficeCode: u.GetOfficeCode(),
	}
}

func bsonUserInfoToPB(u UserInfo) *pb.UserInfo {
	return &pb.UserInfo{
		Name:       u.Name,
		Address:    u.Address,
		PostalCode: u.PostalCode,
		Province:   u.Province,
		OfficeCode: u.OfficeCode,
	}
}

func pbAccessRightsToBSON(ars []*pb.AccessRight) []AccessRight {
	var res []AccessRight
	for _, ar := range ars {
		res = append(res, AccessRight{
			MenuCode:  ar.GetMenuCode(),
			CanCreate: ar.GetCanCreate(),
			CanRead:   ar.GetCanRead(),
			CanUpdate: ar.GetCanUpdate(),
			CanDelete: ar.GetCanDelete(),
		})
	}
	return res
}

func bsonAccessRightsToPB(ars []AccessRight) []*pb.AccessRight {
	var res []*pb.AccessRight
	for _, ar := range ars {
		res = append(res, &pb.AccessRight{
			MenuCode:  ar.MenuCode,
			CanCreate: ar.CanCreate,
			CanRead:   ar.CanRead,
			CanUpdate: ar.CanUpdate,
			CanDelete: ar.CanDelete,
		})
	}
	return res
}

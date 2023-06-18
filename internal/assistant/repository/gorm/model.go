package gorm

import "github.com/defryheryanto/whatsapp-assistant/internal/assistant"

type Role struct {
	Id   int64  `gorm:"primaryKey;autoIncrement;column:id"`
	Name string `gorm:"column:name"`
}

func (Role) TableName() string {
	return "roles"
}

type RoleMember struct {
	Id     int64  `gorm:"primaryKey;autoIncrement;column:id"`
	RoleId int64  `gorm:"column:role_id"`
	Jid    string `gorm:"column:jid"`
}

func (RoleMember) TableName() string {
	return "role_members"
}

func convertToServiceRole(role *Role, members []*RoleMember) *assistant.Role {
	if role == nil {
		return nil
	}
	if members == nil {
		return nil
	}

	result := &assistant.Role{
		Name: role.Name,
	}

	for _, member := range members {
		result.MemberJIDs = append(result.MemberJIDs, member.Jid)
	}

	return result
}

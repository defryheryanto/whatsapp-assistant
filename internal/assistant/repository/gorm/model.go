package gorm

import "github.com/defryheryanto/whatsapp-assistant/internal/assistant"

type Role struct {
	Id       int64  `gorm:"primaryKey;autoIncrement;column:id"`
	Name     string `gorm:"column:name"`
	GroupJid string `gorm:"column:group_jid"`
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

type SavedText struct {
	Id       int64  `gorm:"primaryKey;autoIncrement;column:id"`
	GroupJid string `gorm:"column:group_jid"`
	Title    string `gorm:"column:title"`
	Content  string `gorm:"column:content"`
}

func (SavedText) TableName() string {
	return "saved_texts"
}

type Birthday struct {
	Id            int64  `gorm:"primaryKey;autoIncrement;column:id"`
	Name          string `gorm:"column:name"`
	BirthDate     int16  `gorm:"column:birth_date"`
	BirthMonth    int16  `gorm:"column:birth_month"`
	BirthYear     int16  `gorm:"column:birth_year"`
	TargetChatJid string `gorm:"column:target_chat_jid"`
}

func (Birthday) TableName() string {
	return "birthdays"
}

func (b *Birthday) ToServiceModel() *assistant.Birthday {
	if b == nil {
		return nil
	}

	return &assistant.Birthday{
		Name:          b.Name,
		BirthDate:     b.BirthDate,
		BirthMonth:    b.BirthMonth,
		BirthYear:     b.BirthYear,
		TargetChatJid: b.TargetChatJid,
	}
}

func convertToServiceBirthdays(data []*Birthday) []*assistant.Birthday {
	if data == nil {
		return nil
	}

	birthdays := make([]*assistant.Birthday, 0, len(data))
	for _, d := range data {
		birthdays = append(birthdays, d.ToServiceModel())
	}

	return birthdays
}

type PremiumUser struct {
	Id      int64  `gorm:"primaryKey;autoIncrement;column:id"`
	UserJid string `gorm:"column:user_jid"`
}

func (PremiumUser) TableName() string {
	return "premium_users"
}

func (p *PremiumUser) ToServiceModel() *assistant.PremiumUser {
	if p == nil {
		return nil
	}

	return &assistant.PremiumUser{
		UserJid: p.UserJid,
	}
}

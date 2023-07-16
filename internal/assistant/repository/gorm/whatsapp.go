package gorm

import (
	"context"

	"github.com/defryheryanto/whatsapp-assistant/internal/assistant"
	"gorm.io/gorm"
)

type WhatsAppAssistantRepository struct {
	db *gorm.DB
}

func NewWhatsAppAssistantRepository(db *gorm.DB) *WhatsAppAssistantRepository {
	return &WhatsAppAssistantRepository{db}
}

func (r *WhatsAppAssistantRepository) InsertRole(ctx context.Context, data *assistant.Role) error {
	if data == nil {
		return nil
	}
	if data.MemberJIDs == nil || len(data.MemberJIDs) == 0 {
		return nil
	}

	db := r.db.Begin()

	role := &Role{
		Name:     data.Name,
		GroupJid: data.GroupJid,
	}
	err := db.Create(&role).Error
	if err != nil {
		db.Rollback()
		return err
	}

	members := make([]*RoleMember, len(data.MemberJIDs))
	for i, jid := range data.MemberJIDs {
		members[i] = &RoleMember{
			RoleId: role.Id,
			Jid:    jid,
		}
	}

	err = db.Create(&members).Error
	if err != nil {
		db.Rollback()
		return err
	}

	db.Commit()
	return nil
}

func (r *WhatsAppAssistantRepository) FindRole(ctx context.Context, name, groupJid string) (*assistant.Role, error) {
	var role *Role

	err := r.db.Where("name = ? AND group_jid = ?", name, groupJid).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	var roleMembers []*RoleMember
	err = r.db.Where("role_id = ?", role.Id).Find(&roleMembers).Error
	if err != nil {
		return nil, err
	}

	return convertToServiceRole(role, roleMembers), nil
}

func (r *WhatsAppAssistantRepository) DeleteRole(ctx context.Context, name string) error {
	var role *Role

	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		return err
	}

	db := r.db.Begin()
	err = db.Where("role_id = ?", role.Id).Delete(&RoleMember{}).Error
	if err != nil {
		db.Rollback()
		return err
	}

	err = db.Where("id = ?", role.Id).Delete(&Role{}).Error
	if err != nil {
		db.Rollback()
		return err
	}

	db.Commit()
	return nil
}

func (r *WhatsAppAssistantRepository) SaveText(ctx context.Context, data *assistant.SavedText) error {
	savedText := &SavedText{
		GroupJid: data.GroupJid,
		Title:    data.Title,
		Content:  data.Content,
	}

	db := r.db.Begin()
	err := db.Create(&savedText).Error
	if err != nil {
		db.Rollback()
		return err
	}

	db.Commit()
	return nil
}

func (r *WhatsAppAssistantRepository) GetSavedText(ctx context.Context, groupJid, title string) (*assistant.SavedText, error) {
	var savedText *SavedText

	err := r.db.Where("group_jid = ? AND title = ?", groupJid, title).First(&savedText).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &assistant.SavedText{
		GroupJid: savedText.GroupJid,
		Title:    savedText.Title,
		Content:  savedText.Content,
	}, nil
}

func (r *WhatsAppAssistantRepository) DeleteSavedText(ctx context.Context, groupJid, title string) error {
	db := r.db.Begin()

	err := db.Where("group_jid = ? AND title = ?", groupJid, title).Delete(&SavedText{}).Error
	if err != nil {
		db.Rollback()
		return err
	}

	db.Commit()
	return nil
}

func (r WhatsAppAssistantRepository) InsertBirthday(ctx context.Context, birthday *assistant.Birthday) error {
	db := r.db.Begin()

	newBirthday := &Birthday{
		Name:          birthday.Name,
		BirthDate:     birthday.BirthDate,
		BirthMonth:    birthday.BirthMonth,
		BirthYear:     birthday.BirthYear,
		TargetChatJid: birthday.TargetChatJid,
	}

	err := db.Create(&newBirthday).Error
	if err != nil {
		db.Rollback()
		return err
	}

	db.Commit()
	return nil
}

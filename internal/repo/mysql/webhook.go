package mysql

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *GormRepository) UpsertContactFromWebhook(accountID uint64, amoID int64, name string, email *string, status string) (uint, error) {
	if accountID == 0 {
		return 0, errors.New("account_id must be > 0")
	}
	if amoID <= 0 {
		return 0, errors.New("amo_id must be > 0")
	}
	if strings.TrimSpace(status) == "" {
		return 0, errors.New("status is empty")
	}

	var normEmail *string
	if email != nil {
		v := strings.TrimSpace(*email)
		if v != "" {
			normEmail = &v
		}
	}

	m := Contact{
		AccountID: accountID,
		AmoID:     amoID,
		Name:      strings.TrimSpace(name),
		Email:     normEmail,
		Status:    status,
	}

	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "amo_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "email", "status", "updated_at"}),
	}).Create(&m).Error; err != nil {
		return 0, err
	}

	var got Contact
	if err := r.db.Where("account_id = ? and amo_id = ?", accountID, amoID).First(&got).Error; err != nil {
		return 0, err
	}

	return got.ID, nil
}

func (r *GormRepository) FindContactIDByAmoID(accountID uint64, amoID int64) (uint, bool, error) {
	if accountID == 0 || amoID <= 0 {
		return 0, false, nil
	}

	var got Contact
	err := r.db.Select("id").Where("account_id = ? and amo_id = ?", accountID, amoID).First(&got).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return got.ID, true, nil
}

func (r *GormRepository) AddSyncHistory(contactID uint, status string, message string) error {
	if contactID == 0 {
		return errors.New("contact_id must be > 0")
	}
	h := SyncHistory{
		ContactID: contactID,
		Status:    strings.TrimSpace(status),
		Message:   message,
		CreatedAt: time.Now(),
	}
	if h.Status == "" {
		h.Status = "unknown"
	}
	return r.db.Create(&h).Error
}

func (r *GormRepository) TrimSyncHistory(contactID uint, keepLast int) error {
	if contactID == 0 || keepLast <= 0 {
		return nil
	}

	return r.db.Exec(`
delete from sync_histories
where contact_id = ?
  and id not in (
    select id from (
      select id
      from sync_histories
      where contact_id = ?
      order by created_at desc, id desc
      limit ?
    ) t
  )
`, contactID, contactID, keepLast).Error
}

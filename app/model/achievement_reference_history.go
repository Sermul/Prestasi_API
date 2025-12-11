package model

import "time"

type AchievementReferenceHistory struct {
    ID            string     `json:"id"`             // UUID untuk history entry
    ReferenceID   string     `json:"referenceId"`    // ID dari achievement_reference
    OldStatus     string     `json:"oldStatus"`      // status sebelumnya
    NewStatus     string     `json:"newStatus"`      // status setelah perubahan
    Note          string     `json:"note,omitempty"` // untuk reject
    ChangedBy     string     `json:"changedBy"`      // user_id dari admin / dosen wali / mahasiswa
    ChangedByRole string     `json:"changedByRole"`  // role yang melakukan perubahan
    CreatedAt     time.Time  `json:"createdAt"`      // timestamp perubahan
}

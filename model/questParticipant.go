package model

/* クエスト参加者を管理するテーブルの定義 */

import "time"

type QuestParticipant struct {
	ID       uint      `json:"id" gorm:"primaryKey"` //! このIDを指定することはない
	JoinedAt time.Time `json:"joined_at"`
	User     User      `json:"user" gorm:"foreignKey:UserId;"` // ユーザーが削除されたら参加記録も削除
	UserId   uint      `json:"user_id" gorm:"not null"`
	Quest    Quest     `json:"quest" gorm:"foreignKey:QuestId;references:ID;"` // クエストが削除されたら参加記録も削除
	QuestId  uint      `json:"quest_id" gorm:"not null"`
}

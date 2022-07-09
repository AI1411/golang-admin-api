package models

import "time"

type MemberStatus string

const (
	Premium MemberStatus = "premium"
	Basic   MemberStatus = "basic"
)

type SubscriptionMember struct {
	ID                  string
	UserID              string
	MemberStatus        MemberStatus
	MemberStartDate     time.Time
	MemberEndDate       *time.Time
	MemberStopStartDate time.Time
	MemberStopEndDate   *time.Time
}

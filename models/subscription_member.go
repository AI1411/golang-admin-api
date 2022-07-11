package models

import "time"

type MemberStatus string

const (
	Premium  MemberStatus = "premium"
	Basic    MemberStatus = "basic"
	Inactive MemberStatus = "inactive"
	Stopped  MemberStatus = "stopped"
)

type SubscriptionMember struct {
	ID              string
	UserID          string
	MemberStatus    MemberStatus
	MemberStartDate time.Time
	MemberEndDate   time.Time
}

func NewSubscriptionMember(userID string, memberStatus MemberStatus, memberStartDate,
	memberEndDate time.Time) *SubscriptionMember {
	return &SubscriptionMember{
		UserID:          userID,
		MemberStatus:    memberStatus,
		MemberStartDate: memberStartDate,
		MemberEndDate:   memberEndDate,
	}
}

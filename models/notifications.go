package models

import (
	"sync"
)

type UserNotifications map[string][]Notification

type NotificationStore struct {
	Data UserNotifications
	Mu   sync.RWMutex
}

func (ns *NotificationStore) AddNotification(userID string, notification Notification) {
	ns.Mu.Lock()
	defer ns.Mu.Unlock()

	if _, ok := ns.Data[userID]; !ok {
		ns.Data[userID] = []Notification{}
	}

	ns.Data[userID] = append(ns.Data[userID], notification)
}

func (ns *NotificationStore) GetNotifications(userID string) []Notification {
	ns.Mu.RLock()
	defer ns.Mu.RUnlock()

	return ns.Data[userID]
}

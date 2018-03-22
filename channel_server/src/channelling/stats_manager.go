package channelling

import (
	"sync/atomic"
)

type HubStat struct {
	Rooms                 int                     `json:"rooms"`
	Connections           int                     `json:"connections"`
	Sessions              int                     `json:"sessions"`
	Users                 int                     `json:"users"`
	Count                 uint64                  `json:"count"`
	BroadcastChatMessages uint64                  `json:"broadcastchatmessages"`
	UnicastChatMessages   uint64                  `json:"unicastchatmessages"`
	IdsInRoom             map[string][]string     `json:"idsinroom,omitempty"`
	SessionsById          map[string]*DataSession `json:"sessionsbyid,omitempty"`
	UsersById             map[string]*DataUser    `json:"usersbyid,omitempty"`
	ConnectionsByIdx      map[string]string       `json:"connectionsbyidx,omitempty"`
}

type ConnectionCounter interface {
	CountConnection() uint64
}

type StatsCounter interface {
	CountBroadcastChat()
	CountUnicastChat()
}

type StatsGenerator interface {
	Stat(details bool) *HubStat
}

type StatsManager interface {
	ConnectionCounter
	StatsCounter
	StatsGenerator
}

type statsManager struct {
	ClientStats
	RoomStats
	UserStats
	connectionCount       uint64
	broadcastChatMessages uint64
	unicastChatMessages   uint64
}

func NewStatsManager(clientStats ClientStats, roomStats RoomStats, userStats UserStats) StatsManager {
	return &statsManager{clientStats, roomStats, userStats, 0, 0, 0}
}

func (stats *statsManager) CountConnection() uint64 {
	return atomic.AddUint64(&stats.connectionCount, 1)
}

func (stats *statsManager) CountBroadcastChat() {
	atomic.AddUint64(&stats.broadcastChatMessages, 1)
}

func (stats *statsManager) CountUnicastChat() {
	atomic.AddUint64(&stats.unicastChatMessages, 1)
}

func (stats *statsManager) Stat(details bool) *HubStat {
	roomCount, roomSessionInfo := stats.RoomInfo(details)
	clientCount, sessions, connections := stats.ClientInfo(details)
	userCount, users := stats.UserInfo(details)

	return &HubStat{
		Rooms:       roomCount,
		Connections: clientCount,
		Sessions:    clientCount,
		Users:       userCount,
		Count:       atomic.LoadUint64(&stats.connectionCount),
		BroadcastChatMessages: atomic.LoadUint64(&stats.broadcastChatMessages),
		UnicastChatMessages:   atomic.LoadUint64(&stats.unicastChatMessages),
		IdsInRoom:             roomSessionInfo,
		SessionsById:          sessions,
		UsersById:             users,
		ConnectionsByIdx:      connections,
	}
}

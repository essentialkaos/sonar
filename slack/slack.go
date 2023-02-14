package slack

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/essentialkaos/ek/v12/log"
	"github.com/essentialkaos/ek/v12/pluralize"
	"github.com/essentialkaos/ek/v12/timeutil"

	"github.com/slack-go/slack"

	"github.com/orcaman/concurrent-map"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Status consts
const (
	STATUS_UNKNOWN Status = iota
	STATUS_OFFLINE
	STATUS_ONLINE
	STATUS_DND
	STATUS_DND_OFFLINE
	STATUS_VACATION
	STATUS_ONCALL
	STATUS_DISABLED
)

// Status is Slack status
type Status uint8

// MAX_PRESENCE_CHECK_BATCH maximum number of users per batch for checking presence
const MAX_PRESENCE_CHECK_BATCH = 100

// MAX_DND_BATCH maximum number of users per batch for getting DND status
const MAX_DND_BATCH = 50

// ////////////////////////////////////////////////////////////////////////////////// //

type userMeta struct {
	Online     bool
	Vacation   bool
	OnCall     bool
	Disabled   bool
	DNDStart   int64
	DNDEnd     int64
	DNDUpdated int64

	Email    string
	RealName string

	mutex *sync.RWMutex
}

type dataStore struct {
	// Map mail -> meta
	MailIndex cmap.ConcurrentMap

	// Map ID -> meta
	IDIndex cmap.ConcurrentMap
}

// slackLogProxy is proxy logger for slack
type slackLogProxy struct {
	Prefix string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// slack client
var client *slack.Client

// rtm connection
var rtm *slack.RTM

// user info store
var store *dataStore

// mail mappings
var mappings map[string]string

// connection flag
var connected bool

// ////////////////////////////////////////////////////////////////////////////////// //

// StartObserver start status observer
func StartObserver(token string, mp map[string]string) error {
	client = slack.New(
		token,
		slack.OptionLog(&slackLogProxy{"SLACK:"}), // disabled by default, use for debug purposes
	)

	rtm = client.NewRTM(
		slack.RTMOptionConnParams(url.Values{
			"batch_presence_aware": {"1"},
		}),
	)

	mappings = mp

	store = &dataStore{cmap.New(), cmap.New()}

	err := fetchInitialInfo()

	if err != nil {
		return err
	}

	go rtmLoop()
	go presenceCheckLoop()

	return nil
}

// GetStatus return user status by name
func GetStatus(mail string) Status {
	if mappings != nil && mappings[mail] != "" {
		mail = mappings[mail]
	}

	data, ok := store.MailIndex.Get(mail)

	if !ok {
		log.Debug("Can't find info for user %s", mail)
		return STATUS_UNKNOWN
	}

	meta := data.(*userMeta)

	meta.mutex.RLock()
	defer meta.mutex.RUnlock()

	if meta.Disabled {
		return STATUS_DISABLED
	}

	if meta.Vacation {
		return STATUS_VACATION
	}

	if meta.Online {
		if meta.IsDND() {
			return STATUS_DND
		}

		if meta.OnCall {
			return STATUS_ONCALL
		}

		return STATUS_ONLINE
	}

	if meta.IsDND() {
		return STATUS_DND_OFFLINE
	}

	return STATUS_OFFLINE
}

// ////////////////////////////////////////////////////////////////////////////////// //

// IsDND return true if user in DND
func (s *userMeta) IsDND() bool {
	now := time.Now().Unix()

	if now > s.DNDStart && now < s.DNDEnd {
		return true
	}

	return false
}

// ////////////////////////////////////////////////////////////////////////////////// //

// rtm message handling loop
func rtmLoop() {
	go rtm.ManageConnection()

	for {
		select {
		case event := <-rtm.IncomingEvents:
			log.Debug("Got slack event %v", event)
			switch event.Data.(type) {

			case *slack.ConnectingEvent:
				log.Info("Connecting to Slack…")

			case *slack.ConnectedEvent:
				connected = true
				log.Info("Connected to Slack")

			case *slack.DisconnectedEvent:
				connected = false
				log.Warn("Disconnected from Slack")

			case *slack.ConnectionErrorEvent:
				connected = false
				log.Warn(
					"Slack connecting error: %s",
					event.Data.(*slack.ConnectionErrorEvent).Error(),
				)

			case *slack.HelloEvent:
				sendPresenceQuery()

			case *slack.DNDUpdatedEvent:
				ev := event.Data.(*slack.DNDUpdatedEvent)
				updateUserDND(ev.User, ev.Status)

			case *slack.PresenceChangeEvent:
				ev := event.Data.(*slack.PresenceChangeEvent)
				if ev.User != "" {
					updateUserPresence([]string{ev.User}, ev.Presence == "active")
				} else {
					updateUserPresence(ev.Users, ev.Presence == "active")
				}

			case *slack.UserChangeEvent:
				ev := event.Data.(*slack.UserChangeEvent)

				if !store.IDIndex.Has(ev.User.ID) {
					addNewUser(ev.User, nil)
				} else {
					updateUserStatus(ev.User)
				}
			}
		}
	}
}

// presenceCheckLoop is presence check loop
func presenceCheckLoop() {
	for range time.NewTicker(3 * time.Minute).C {
		if connected {
			sendPresenceQuery()
		}
	}
}

// sendPresenceQuery send presence query message
func sendPresenceQuery() {
	var ids []string

	keys := store.IDIndex.Keys()
	totalUsers := len(keys)
	counter := 0

	log.Info(
		"Sending presence query messages (%s per message)…",
		pluralize.P("%d %s", MAX_PRESENCE_CHECK_BATCH, "user", "users"),
	)

	for index, id := range keys {
		data, _ := store.IDIndex.Get(id)
		meta := data.(*userMeta)

		meta.mutex.RLock()

		if !meta.Disabled {
			ids = append(ids, id)
			counter++
		}

		meta.mutex.RUnlock()

		if len(ids) == MAX_PRESENCE_CHECK_BATCH || index+1 == totalUsers {
			rtm.SendMessage(&slack.OutgoingMessage{Type: "presence_sub", IDs: ids})
			time.Sleep(time.Second)
			ids = nil
		}
	}

	log.Info(
		"Presence query messages successfully sent (%s)",
		pluralize.P("%d %s", counter, "user", "users"),
	)
}

// fetchInitialInfo fetch initial info
func fetchInitialInfo() error {
	log.Info("Fetching initial data…")

	users, err := fetchTeamUsers()

	if err != nil {
		return fmt.Errorf("Can't fetch users: %v", err)
	}

	dndInfo, err := fetchDNDInfo(users)

	if err != nil {
		return fmt.Errorf("Can't fetch DND info: %v", err)
	}

	for _, user := range users {
		addNewUser(user, dndInfo)
	}

	log.Info(
		"Information about %d users added to storage (%d have DND)",
		len(users), len(dndInfo),
	)

	return nil
}

// addNewUser add new user to store
func addNewUser(user slack.User, dndInfo map[string]slack.DNDStatus) {
	if user.IsBot || user.Deleted {
		return
	}

	meta := &userMeta{mutex: &sync.RWMutex{}}

	meta.Email = user.Profile.Email
	meta.RealName = user.RealName

	if !user.Deleted {
		meta.Online = user.Presence == "active"
		meta.Vacation = strings.HasPrefix(user.RealName, "[")

		if dndInfo != nil {
			dnd, ok := dndInfo[user.ID]

			if ok {
				meta.DNDUpdated = time.Now().Unix()
				meta.DNDStart = int64(dnd.NextStartTimestamp)
				meta.DNDEnd = int64(dnd.NextEndTimestamp)
			}
		}
	} else {
		meta.Disabled = true
	}

	store.MailIndex.Set(user.Profile.Email, meta)
	store.IDIndex.Set(user.ID, meta)

	log.Info("Appended new user %s (%s - %s)", user.Profile.Email, user.ID, user.RealName)
}

// updateUserDND update user DND times
func updateUserDND(id string, status slack.DNDStatus) {
	data, ok := store.IDIndex.Get(id)

	if !ok {
		log.Warn("Can't find user with ID %s for DND update", id)
		return
	}

	meta := data.(*userMeta)

	meta.mutex.Lock()
	meta.DNDStart = int64(status.NextStartTimestamp)
	meta.DNDEnd = int64(status.NextEndTimestamp)
	meta.mutex.Unlock()

	if meta.DNDStart < 86400 {
		dndStart := timeutil.Format(time.Unix(meta.DNDStart, 0), "%Y/%m/%d %H:%M")
		dndEnd := timeutil.Format(time.Unix(meta.DNDEnd, 0), "%Y/%m/%d %H:%M")
		log.Info(
			"Updated DND (%s ↔ %s) for user %s (%s - %s)",
			dndStart, dndEnd, meta.Email, id, meta.RealName,
		)
	} else {
		log.Info(
			"Cleared DND for user %s (%s - %s)",
			meta.Email, id, meta.RealName,
		)
	}
}

// updateUserPresence update user presence
func updateUserPresence(ids []string, online bool) {
	for _, id := range ids {
		data, ok := store.IDIndex.Get(id)

		if !ok {
			log.Warn("Can't find user with ID %s for presence update", id)
			return
		}

		meta := data.(*userMeta)

		meta.mutex.RLock()

		if meta.Online == online {
			meta.mutex.RUnlock()
			continue
		}

		meta.mutex.RUnlock()

		meta.mutex.Lock()
		meta.Online = online
		meta.mutex.Unlock()

		checkUserDND(id, meta)

		log.Info(
			"Updated presence (online: %t) for user %s (%s - %s)",
			online, meta.Email, id, meta.RealName,
		)
	}
}

// checkUserDND check if we should update user dnd status
func checkUserDND(id string, meta *userMeta) {
	if meta.DNDStart == 0 {
		return
	}

	now := time.Now().Unix()

	if meta.DNDEnd > now || now-meta.DNDUpdated < 180 {
		return
	}

	status, err := client.GetDNDInfo(&id)

	if err != nil {
		log.Error("Can't check user DND status: %v", err)
		return
	}

	meta.mutex.Lock()
	meta.DNDStart = int64(status.NextStartTimestamp)
	meta.DNDEnd = int64(status.NextEndTimestamp)
	meta.DNDUpdated = now
	meta.mutex.Unlock()

	if meta.DNDStart > 86400 {
		dndStart := timeutil.Format(time.Unix(meta.DNDStart, 0), "%Y/%m/%d %H:%M")
		dndEnd := timeutil.Format(time.Unix(meta.DNDEnd, 0), "%Y/%m/%d %H:%M")
		log.Info(
			"Checked and updated DND (%s ↔ %s) for user %s (%s - %s)",
			dndStart, dndEnd, meta.Email, id, meta.RealName,
		)
	}
}

// updateUserStatus user vacation status
func updateUserStatus(user slack.User) {
	if user.IsBot || user.Deleted {
		return
	}

	data, ok := store.IDIndex.Get(user.ID)

	if !ok {
		log.Warn("Can't find user with ID %s for vacation update", user)
		return
	}

	meta := data.(*userMeta)

	meta.mutex.Lock()

	// Update vacation status
	meta.Vacation = strings.HasPrefix(user.RealName, "[")

	// Update on-call status
	if user.Profile.StatusEmoji == ":slack_call:" && user.Profile.StatusText == "On a call" {
		meta.OnCall = true
	} else {
		if meta.OnCall {
			meta.OnCall = false
		}
	}

	meta.mutex.Unlock()

	log.Info("Checked status for user %s (%s - %s)", user.Profile.Email, user.ID, user.RealName)
}

// fetchTeamUsers fetches users from Slack
func fetchTeamUsers() ([]slack.User, error) {
	users, err := client.GetUsers()

	if err != nil {
		return nil, err
	}

	var result []slack.User

	for _, user := range users {
		if user.IsBot || user.IsStranger {
			continue
		}

		result = append(result, user)
	}

	return result, nil
}

// fetchDNDInfo fetches DND info for all users
func fetchDNDInfo(users []slack.User) (map[string]slack.DNDStatus, error) {
	dndInfo := make(map[string]slack.DNDStatus)

	log.Info("Fetching DND info…")

	for i := 0; i < 1000; i++ {
		ids := getUsersBatch(users, MAX_DND_BATCH, i)

		if len(ids) == 0 {
			break
		}

		info, err := client.GetDNDTeamInfo(ids)

		if err != nil {
			return nil, err
		}

		appendDNDDdata(dndInfo, info)

		log.Info("Added DND info for %d users", (i+1)*MAX_DND_BATCH)

		time.Sleep(3 * time.Second)
	}

	return dndInfo, nil
}

// getUsersBatch returns batch with users IDs
func getUsersBatch(users []slack.User, size, index int) []string {
	var result []string

	start := size * index
	end := start + size

	if end >= len(users) {
		end = len(users) - 1
	}

	for i := start; i < end; i++ {
		if users[i].Deleted {
			continue
		}

		result = append(result, users[i].ID)
	}

	return result
}

// appendDNDDdata appends new dnd data
func appendDNDDdata(target, source map[string]slack.DNDStatus) {
	for id, status := range source {
		target[id] = status
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Output writes log message to default logger with prefix
func (s *slackLogProxy) Output(calldepth int, message string) error {
	return log.Info("%s %s", s.Prefix, message)
}

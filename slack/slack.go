package slack

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"strings"
	"sync"
	"time"

	"pkg.re/essentialkaos/ek.v9/log"

	"pkg.re/essentialkaos/slack.v2"

	"github.com/orcaman/concurrent-map"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Status consts
const (
	STATUS_UNKNOWN Status = iota
	STATUS_OFFLINE
	STATUS_ONLINE
	STATUS_DND
	STATUS_VACATION
	STATUS_ONCALL
)

// Status is Slack status
type Status uint8

// ////////////////////////////////////////////////////////////////////////////////// //

type userMeta struct {
	Online     bool
	Vacation   bool
	OnCall     bool
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

// ////////////////////////////////////////////////////////////////////////////////// //

// slack client
var client *slack.Client

// rtm connection
var rtm *slack.RTM

// user info store
var store *dataStore

// mail mappings
var mappings map[string]string

// ////////////////////////////////////////////////////////////////////////////////// //

// StartObserver start status observer
func StartObserver(token string, mp map[string]string) error {
	client = slack.New(token)
	client.Config.BatchPresenceAware = true
	client.Config.PresenceSub = true

	rtm = client.NewRTM()
	mappings = mp

	store = &dataStore{cmap.New(), cmap.New()}

	err := fetchInitialInfo()

	if err != nil {
		return err
	}

	go rtmLoop()

	return nil
}

// GetStatus return user status by name
func GetStatus(mail string) Status {
	if mappings != nil && mappings[mail] != "" {
		mail = mappings[mail]
	}

	data, ok := store.MailIndex.Get(mail)

	if !ok {
		log.Warn("Can't find info for user %s", mail)
		return STATUS_UNKNOWN
	}

	meta := data.(*userMeta)

	meta.mutex.RLock()
	defer meta.mutex.RUnlock()

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
				log.Info("Connecting to Slack...")

			case *slack.ConnectedEvent:
				log.Info("Connected to Slack")

			case *slack.HelloEvent:
				subscribeToPresenceEvents()

			case *slack.DisconnectedEvent:
				log.Warn("Disconnected from Slack")

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
					subscribeToPresenceEvents()
				} else {
					updateUserStatus(ev.User)
				}
			}
		}
	}
}

// subscribeToPresenceEvents subscribe bot to events about presence change for all users
func subscribeToPresenceEvents() {
	var usersToSub []string

	for _, id := range store.IDIndex.Keys() {
		log.Debug("Added %s to subscribe list", id)
		usersToSub = append(usersToSub, id)
	}

	err := rtm.PresenceSub(usersToSub)

	if err != nil {
		log.Error(err.Error())
	}

	log.Info("Subscribed to presence events for all users")
}

// fetchInitialInfo fetch initial info
func fetchInitialInfo() error {
	users, err := client.GetUsers()

	if err != nil {
		return err
	}

	dndInfo, err := client.GetDNDTeamInfo(nil)

	if err != nil {
		return err
	}

	for _, user := range users {
		addNewUser(user, dndInfo)
	}

	return nil
}

// addNewUser add new user to store
func addNewUser(user slack.User, dndInfo map[string]slack.DNDStatus) {
	if user.Deleted || user.IsBot {
		return
	}

	now := time.Now().Unix()
	meta := &userMeta{mutex: &sync.RWMutex{}}

	meta.Online = user.Presence == "active"
	meta.Vacation = strings.HasPrefix(user.RealName, "[")

	meta.Email = user.Profile.Email
	meta.RealName = user.RealName

	if dndInfo != nil {
		dnd, ok := dndInfo[user.ID]

		if ok {
			meta.DNDUpdated = now
			meta.DNDStart = int64(dnd.NextStartTimestamp)
			meta.DNDEnd = int64(dnd.NextEndTimestamp)
		}
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

	log.Info("Updated DND for user %s (%s - %s)", meta.Email, id, meta.RealName)
}

// updateUserPresence update user presence
func updateUserPresence(ids []string, online bool) {
	if len(ids) > 1 {
		log.Info("Got batched presence status (%d events)", len(ids))
	}

	for _, id := range ids {
		data, ok := store.IDIndex.Get(id)

		if !ok {
			log.Warn("Can't find user with ID %s for presence update", id)
			return
		}

		meta := data.(*userMeta)

		meta.mutex.Lock()
		meta.Online = online
		meta.mutex.Unlock()

		checkUserDND(id, meta)

		log.Info("Updated presence for user %s (%s - %s)", meta.Email, id, meta.RealName)
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

	log.Info("Checked and updated DND for user %s (%s - %s)", meta.Email, id, meta.RealName)
}

// updateUserStatus user vacation status
func updateUserStatus(user slack.User) {
	if user.IsBot {
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

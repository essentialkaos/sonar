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

	"github.com/nlopes/slack"
	"github.com/orcaman/concurrent-map"
)

const (
	STATUS_UNKNOWN Status = iota
	STATUS_OFFLINE
	STATUS_ONLINE
	STATUS_DND
	STATUS_VACATION
)

type Status uint8

// ////////////////////////////////////////////////////////////////////////////////// //

type userMeta struct {
	Online   bool
	Vacation bool
	DNDStart int64
	DNDEnd   int64

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

var (
	client *slack.Client
	rtm    *slack.RTM
	store  *dataStore
)

// ////////////////////////////////////////////////////////////////////////////////// //

// StartObserver start status observer
func StartObserver(token string) error {
	client = slack.New(token)
	rtm = client.NewRTM()

	store = &dataStore{cmap.New(), cmap.New()}

	err := fetchInitialInfo()

	if err != nil {
		return err
	}

	go rtmLoop()

	return nil
}

// GetStatus return user status by name
func GetStatus(email string) Status {
	log.Debug("Got status request for %s", email)

	data, ok := store.MailIndex.Get(strings.Replace(email, ".", "", -1))

	if !ok {
		log.Warn("Can't find info for user %s", email)
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

		return STATUS_ONLINE
	}

	return STATUS_OFFLINE
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *userMeta) IsDND() bool {
	now := time.Now().Unix()

	if now > s.DNDStart && now < s.DNDEnd {
		return true
	}

	return false
}

// ////////////////////////////////////////////////////////////////////////////////// //

func rtmLoop() {
	go rtm.ManageConnection()

	for {
		for {
			select {
			case event := <-rtm.IncomingEvents:
				switch event.Data.(type) {
				case *slack.ConnectingEvent:
					log.Info("Connecting to Slack...")
				case *slack.ConnectedEvent:
					log.Info("Connected to Slack")
				case *slack.DisconnectedEvent:
					log.Warn("Disconnected from Slack")

				case *slack.DNDUpdatedEvent:
					ev := event.Data.(*slack.DNDUpdatedEvent)
					updateUserDND(ev.User, ev.Status)
				case *slack.PresenceChangeEvent:
					ev := event.Data.(*slack.PresenceChangeEvent)
					updateUserPresence(ev.User, ev.Presence == "active")
				case *slack.UserChangeEvent:
					ev := event.Data.(*slack.UserChangeEvent)

					if !store.IDIndex.Has(ev.User.ID) {
						addNewUser(ev.User, nil)
					} else {
						updateUserVacation(ev.User)
					}
				}
			}
		}
	}
}

// fetchInitialInfo fetch initial info
func fetchInitialInfo() error {
	users, err := client.GetUsers()

	if err != nil {
		return err
	}

	dndInfo, err := client.GetDNDTeamInfo(nil)

	if err != nil {
		return nil
	}

	for _, user := range users {
		addNewUser(user, dndInfo)
	}

	return nil
}

// addNewUser add new user to store
func addNewUser(user slack.User, dndInfo map[string]slack.DNDStatus) {
	if user.Deleted {
		return
	}

	meta := &userMeta{mutex: &sync.RWMutex{}}

	meta.Online = user.Presence == "active"
	meta.Vacation = strings.HasPrefix(user.RealName, "[")

	meta.Email = user.Profile.Email
	meta.RealName = user.RealName

	if dndInfo != nil {
		dnd, ok := dndInfo[user.ID]

		if ok {
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

	log.Debug("Updated DND for user %s (%s - %s)", meta.Email, id, meta.RealName)
}

// updateUserPresence update user presence
func updateUserPresence(id string, online bool) {
	data, ok := store.IDIndex.Get(id)

	if !ok {
		log.Warn("Can't find user with ID %s for presence update", id)
		return
	}

	meta := data.(*userMeta)

	meta.mutex.Lock()
	meta.Online = online
	meta.mutex.Unlock()

	log.Debug("Updated presence for user %s (%s - %s)", meta.Email, id, meta.RealName)
}

// updateUserVacation user vacation status
func updateUserVacation(user slack.User) {
	data, ok := store.IDIndex.Get(user.ID)

	if !ok {
		log.Warn("Can't find user with ID %s for vacation update", user)
		return
	}

	meta := data.(*userMeta)

	meta.mutex.Lock()
	meta.Vacation = strings.HasPrefix(user.RealName, "[")
	meta.mutex.Unlock()

	log.Debug("Updated vacation status for user %s (%s - %s)", user.Profile.Email, user.ID, user.RealName)
}

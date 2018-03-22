package messages



// BundleRef is a lightweight record describing the bundle name
// and version installed on a Relay
type BundleRef struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// Announcement describes the online/offline status of a Relay
type Announcement struct {
	ID      string      `json:"announcement_id,omitempty" valid:"-"`
	RelayID string      `json:"relay" valid:"required"`
	Online  bool        `json:"online" valid:"bool,required"`
	Bundles []BundleRef `json:"bundles,omitempty"`
	// Deprecated
	Snapshot bool   `json:"snapshot" valid:"bool,required"`
	ReplyTo  string `json:"reply_to,omitempty" valid:"-"`
}



// AnnouncementEnvelope is a wrapper around an Announcement directive.
type AnnouncementEnvelope struct {
	Announcement *Announcement `json:"announce" valid:"required"`
}


// NewOfflineAnnouncement builds an Announcement informing Cog the Relay is offline
func NewOfflineAnnouncement(relayID string, replyTo string) *AnnouncementEnvelope {
	return &AnnouncementEnvelope{
		Announcement: &Announcement{
			ID:       "0",
			RelayID:  relayID,
			Online:   false,
			Snapshot: true,
			ReplyTo:  replyTo,
		},
	}
}
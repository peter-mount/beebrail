package status

import (
	"encoding/xml"
)

// snapshot of the current status
type SnapshotCategory struct {
	Name    string        `json:"name" xml:"name,attr"`                  // Name of this category
	Title   string        `json:"title" xml:"title,attr,omitempty"`      // Title of this category
	Port    uint16        `json:"port,omitempty" xml:"port,omitempty"`   // Optional port for this category
	Stats   Stats         `json:"stats" xml:"stats"`                     // Statistics
	Entries []*Connection `json:"connection,omitempty" xml:"connection"` // Connections in this category
}

// Snapshot wrapper to give a consistent XML document
type Snapshot struct {
	XMLName xml.Name            `json:"-" xml:"status"`
	Entries []*SnapshotCategory `json:"category" xml:"category"`
}

func (s *Status) Snapshot() *Snapshot {
	ret := &Snapshot{}

	_ = s.ForEach(func(category *Category) error {
		cat := &SnapshotCategory{
			Name:  category.Name,
			Title: category.Title,
			Port:  category.Port,
			Stats: category.Stats,
		}
		ret.Entries = append(ret.Entries, cat)
		return nil
	})

	for _, cat := range ret.Entries {
		s.GetCategory(cat.Name).ForEach(func(con *Connection) error {
			cat.Entries = append(cat.Entries, con.clone())
			return nil
		})
	}

	return ret
}

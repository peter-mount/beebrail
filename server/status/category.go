package status

import (
	"log"
	"sort"
	"time"
)

// Common structure used to manage times, duration, idle & bytes processed
type Stats struct {
	Time       time.Time     `json:"time" xml:"time,attr"`         // When the connection was made
	Duration   time.Duration `json:"duration" xml:"duration,attr"` // How long the connecton has been running
	Idle       time.Duration `json:"idle" xml:"idle,attr"`         // When last active
	LastActive time.Time     `json:"-" xml:"-"`                    // Last active time, not exported
	BytesIn    int           `json:"bytesIn" xml:"bytesIn"`        // Bytes received
	BytesOut   int           `json:"bytesOut" xml:"bytesOut"`      // Bytes received
}

// Category of connections
type Category struct {
	Name            string              // Name of this category
	Title           string              // Title of this category
	Port            uint16              // Optional port identifier
	ConnectionCount int                 // Counter of number of connections
	Stats           Stats               // Statistics
	s               *Status             // Link to underlying Status
	connections     map[int]*Connection // map of connections in this Category
}

func (c *Category) Status() *Status {
	return c.s
}

// Adds a category to the system if it doesn't already exist
func (s *Status) AddCategory(name, title string) *Category {
	var cat *Category

	_ = s.Invoke(func() error {
		if existing, exists := s.categories[name]; exists {
			cat = existing
		} else {
			now := time.Now()
			cat = &Category{
				Name:        name,
				Title:       title,
				s:           s,
				connections: make(map[int]*Connection),
				Stats: Stats{
					Time:       now,
					LastActive: now,
				},
			}

			s.categories[name] = cat
			log.Println("Adding status category", name, title)
		}

		return nil
	})

	return cat
}

// Returns a named Category or nil if not present
func (s *Status) GetCategory(name string) *Category {
	var cat *Category
	_ = s.Invoke(func() error {
		cat = s.categories[name]
		return nil
	})
	return cat
}

// ForEach calls a function for each connection currently active.
func (c *Category) ForEach(f func(connection *Connection) error) error {
	return c.s.Invoke(func() error {
		// Sorted list of connections
		var cons []int
		for k, _ := range c.connections {
			cons = append(cons, k)
		}
		sort.SliceStable(cons, func(i, j int) bool {
			return cons[i] < cons[j]
		})

		for _, k := range cons {
			err := f(c.connections[k])
			if err != nil {
				return err
			}
		}

		return nil
	})
}

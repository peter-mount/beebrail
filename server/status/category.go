package status

import (
	"log"
	"sort"
)

// Category of connections
type Category struct {
	Name            string              // Name of this category
	Title           string              // Title of this category
	Port            uint16              // Optional port identifier
	ConnectionCount int                 // Counter of number of connections
	s               *Status             // Link to underlying Status
	connections     map[int]*Connection // map of connections in this Category
}

// Adds a category to the system if it doesn't already exist
func (s *Status) AddCategory(name, title string) *Category {
	var cat *Category

	_ = s.invoke(func() error {
		if existing, exists := s.categories[name]; exists {
			cat = existing
		} else {
			cat = &Category{
				Name:        name,
				Title:       title,
				s:           s,
				connections: make(map[int]*Connection),
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
	_ = s.invoke(func() error {
		cat = s.categories[name]
		return nil
	})
	return cat
}

// ForEach calls a function for each connection currently active.
func (c *Category) ForEach(f func(connection *Connection) error) error {
	return c.s.invoke(func() error {
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

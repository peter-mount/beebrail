package status

import (
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/golib/rest"
	"github.com/peter-mount/sortfold"
	"sort"
	"sync"
	"time"
)

// Status presents a JSON/HTML endpoint which we can use to provide monitoring of the service
type Status struct {
	mutex       sync.Mutex           // Mutex to allow concurrent access
	restService *rest.Server         // Rest http server
	categories  map[string]*Category // map of categories by name
	started     time.Time            // Time started
}

func (s *Status) Name() string {
	return "Status"
}

func (s *Status) Init(k *kernel.Kernel) error {
	s.categories = make(map[string]*Category)

	service, err := k.AddService(&rest.Server{})
	if err != nil {
		return err
	}
	s.restService = (service).(*rest.Server)

	return nil
}

func (s *Status) PostInit() error {
	s.started = time.Now().UTC()

	s.restService.Handle("/status.json", s.jsonHandler)
	s.restService.Handle("/status.xml", s.xmlHandler)
	s.restService.Handle("/", s.htmlHandler)
	return nil
}

// invoke invokes a function within our mutex
func (s *Status) invoke(f func() error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return f()
}

// ForEach calls a function for each connection currently active.
func (s *Status) ForEach(f func(category *Category) error) error {
	// Sorted list of categories by Title
	var cats []*Category
	_ = s.invoke(func() error {
		for _, cat := range s.categories {
			cats = append(cats, cat)
		}
		return nil
	})

	sort.SliceStable(cats, func(i, j int) bool {
		r := sortfold.CompareFold(cats[i].Title, cats[j].Title)
		if r == 0 {
			r = int(cats[j].Port - cats[i].Port)
		}
		return r < 0
	})

	for _, cat := range cats {
		err := f(cat)
		if err != nil {
			return err
		}
	}

	return nil
}

package observability

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type counter struct{ v atomic.Int64 }

type Registry struct {
	mu       sync.RWMutex
	counters map[string]*counter
}

var M = NewRegistry()

func NewRegistry() *Registry { return &Registry{counters: map[string]*counter{}} }

func (r *Registry) Inc(name string) {
	r.mu.RLock()
	c := r.counters[name]
	r.mu.RUnlock()
	if c == nil {
		r.mu.Lock()
		if r.counters[name] == nil {
			r.counters[name] = &counter{}
		}
		c = r.counters[name]
		r.mu.Unlock()
	}
	c.v.Add(1)
}

func (r *Registry) Add(name string, d int64) {
	r.mu.RLock()
	c := r.counters[name]
	r.mu.RUnlock()
	if c == nil {
		r.mu.Lock()
		if r.counters[name] == nil {
			r.counters[name] = &counter{}
		}
		c = r.counters[name]
		r.mu.Unlock()
	}
	c.v.Add(d)
}

func (r *Registry) PrometheusText() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := make([]string, 0, len(r.counters))
	for k := range r.counters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		name := sanitize(k)
		v := r.counters[k].v.Load()
		b.WriteString(fmt.Sprintf("# TYPE %s counter\n%s %d\n", name, name, v))
	}
	return b.String()
}

func sanitize(s string) string {
	rep := strings.NewReplacer(":", "_", "-", "_", ".", "_")
	return rep.Replace(s)
}

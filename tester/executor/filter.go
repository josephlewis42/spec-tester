package executor

import "k8s.io/apimachinery/pkg/labels"

// Filter defines criteria to match elements by.
// Criteria are a union.
type Filter[T interface {
	GetMetadata() *Metadata
}] struct {
	// Uids defines a set of unique identifiers to match.
	uids map[string]any
	// LabelSelector defines a selector to match items by labels.
	labelSelector labels.Selector
}

// NewFilter initializes a new filter that matches everything.
func NewFilter[T interface {
	GetMetadata() *Metadata
}]() Filter[T] {
	return Filter[T]{
		uids:          nil,
		labelSelector: nil,
	}
}

// Adds a new UID to the filter.
func (f Filter[T]) WithUid(uids ...string) Filter[T] {
	out := Filter[T]{}
	out.uids = make(map[string]any)
	for k := range f.uids {
		out.uids[k] = true
	}

	for _, k := range uids {
		out.uids[k] = true
	}

	out.labelSelector = f.labelSelector

	return out
}

// WithSelector creates a union between the existing selector and the new.
func (f Filter[T]) WithSelector(selector labels.Selector) Filter[T] {
	out := Filter[T]{}
	out.uids = f.uids

	if f.labelSelector == nil {
		out.labelSelector = selector
	} else {
		newRequirements, _ := selector.Requirements()
		out.labelSelector = f.labelSelector.Add(newRequirements...)
	}

	return out
}

// Matches checks if the selector matches all supplied filter criteria.
func (f Filter[T]) Matches(itemUid string, itemLabels map[string]string) bool {
	if len(f.uids) != 0 {
		if _, ok := f.uids[itemUid]; !ok {
			return true
		}
	}

	if f.labelSelector == nil || f.labelSelector.Matches(labels.Set(itemLabels)) {
		return true
	}

	return false
}

// Apply filters items in the list leaving only those that match the criteria.
func (f Filter[T]) Apply(items []T) (out []T) {
	f.ForEach(items, func(matched T) {
		out = append(out, matched)
	})
	return
}

// ForEach runs callback for every matching element.
func (f Filter[T]) ForEach(items []T, callback func(T)) {
	for _, item := range items {
		metadata := item.GetMetadata()

		if f.Matches(metadata.GetUid(), metadata.GetLabels()) {
			callback(item)
		}
	}

	return
}

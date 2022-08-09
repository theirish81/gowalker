package gowalker

// SubTemplates is a collection of templates
type SubTemplates map[string]string

// NewSubTemplates is a constructor of SubTemplates
func NewSubTemplates() SubTemplates {
	return SubTemplates{}
}

// Add adds a template to the collection
func (s *SubTemplates) Add(name string, template string) SubTemplates {
	(*s)[name] = template
	return *s
}

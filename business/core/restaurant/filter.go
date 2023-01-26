package restaurant

// QueryFilter holds the available fields to search for restaurants
type QueryFilter struct {
	ID       *string `validate:"omitempty,uuid4"`
	Name     *string `validate:"omitempty,min=3"`
	Location *string `validate:"omitempty,min=3"`
	Distance *int    `validate:"omitempty,numeric"`
}

// ByID sets the ID of the QueryFilter value.
func (f *QueryFilter) ByID(id string) {
	if id != "" {
		f.ID = &id
	}
}

// ByName sets the Name field of the QueryFilter value
func (f *QueryFilter) ByName(name string) {
	if name != "" {
		f.Name = &name
	}
}

// ByLocation sets the Location field of the QueryFilter value.
func (f *QueryFilter) ByLocation(location string) {
	if location != "" {
		f.Location = &location
	}
}

// ByDistance sets the Distance field of the QueryFilter value.
func (f *QueryFilter) ByDistance(distance int) {
	if distance != 0 {
		f.Distance = &distance
	}
}

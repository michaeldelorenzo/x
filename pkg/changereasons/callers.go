package changereasons

type ReasonList []Reason

// Find takes an id and returns a reason for change or nil
func (list ReasonList) Find(id string) *Reason {
	for _, r := range list {
		if r.ID == id {
			return &r
		}
	}

	return nil
}

// IsValidReason determines if a provide id is a valid reason for change
func IsValidReason(category Category, id string) bool {
	if reasonList, ok := ReasonsForChangeByCategoryAndAction[category]; ok {
		reason := reasonList.Find(id)
		if reason != nil {
			return true
		}
	}

	return false
}

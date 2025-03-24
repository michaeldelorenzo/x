package changereasons_test

import (
	"testing"

	"github.com/koneksahealth/x/pkg/changereasons"
)

func TestFindReasonSuccess(t *testing.T) {
	reasonList := changereasons.ReasonList{
		changereasons.AssessmentManuallyLaunched,
		changereasons.DataModification,
		changereasons.DeviceChange,
	}

	selectedReason := reasonList.Find(changereasons.DeviceChange.ID)

	if selectedReason == nil {
		t.Error("expected a reason but got nil")
		return
	}

	if selectedReason.ID != changereasons.DeviceChange.ID {
		t.Errorf("unexpected reason returned: %s", selectedReason.ID)
	}
}

func TestNoReasonFound(t *testing.T) {
	reasonList := changereasons.ReasonList{
		changereasons.AssessmentManuallyLaunched,
		changereasons.DataModification,
		changereasons.DeviceChange,
	}

	selectedReason := reasonList.Find("invalid-reason-id")

	if selectedReason != nil {
		t.Errorf("unexpected reason found: %s", selectedReason.ID)
	}
}

func TestReasons(t *testing.T) {
	t.Run("valid reason for change", func(t *testing.T) {
		validCategory := changereasons.CategoryParticipantStatusRemove
		validReasonID := changereasons.StudyStaffModification.ID

		valid := changereasons.IsValidReason(validCategory, validReasonID)
		if !valid {
			t.Error("expected reason for change to be valid")
		}
	})

	t.Run("invalid category", func(t *testing.T) {
		invalidCategory := changereasons.CategoryUnknown
		validReasonID := changereasons.ParticipantModification.ID

		valid := changereasons.IsValidReason(invalidCategory, validReasonID)
		if valid {
			t.Error("expected category to be invalid")
		}
	})

	t.Run("invalid reason", func(t *testing.T) {
		validCategory := changereasons.CategoryParticipantStatusRemove
		invalidReasonID := "invalid-reason-id"

		valid := changereasons.IsValidReason(validCategory, invalidReasonID)
		if valid {
			t.Error("expected reason for change to be valid")
		}
	})
}

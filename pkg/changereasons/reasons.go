package changereasons

// Reason defines a reason for changing/updating/deleting a resource
type Reason struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type Category string

const (
	CategoryUnknown Category = ""
	// CategoryParticipantStatusRevert is to be used when a status is either getting rolled back
	// or a participant is switching between terminal states. It is an error category used when fixing a mistake.
	CategoryParticipantStatusRevert Category = "participant-status-revert"
	// CategoryParticipantStatusRemove is to be used when removing a participant pre-enrollment or post-enrollment
	// without them completing the study
	CategoryParticipantStatusRemove Category = "participant-status-remove"
	// CategoryParticipantStatusAdvance is to be used when moving a participant from pre-enrollment to enrolled
	CategoryParticipantStatusAdvance Category = "participant-status-advance"
	// CategoryParticipantStatusComplete is to be used when a participant completes the study
	CategoryParticipantStatusComplete Category = "participant-status-complete"
	// CategoryGenericCreation is to be used for general entities (Orgs, Sponsors, Etc.) to validate the `InitialCreation` reason
	// over time, some entities may opt out of this category for something more specific if the need arises
	CategoryGenericCreation Category = "generic-creation"
	// CategoryGenericUpdate is to be used for general entities (Orgs, Sponsors, Etc.) to validate the `Data Modification` reason
	// over time, some entities may opt out of this category for something more specific if the need arises
	CategoryGenericUpdate Category = "generic-update"

	CategoryParticipantDetailsCreate              Category = "participant-details-create"
	CategoryParticipantDetailsUpdate              Category = "participant-details-update"
	CategoryParticipantStudyScheduleCreate        Category = "participant-study-schedule-create"
	CategoryParticipantStudyScheduleUpdate        Category = "participant-study-schedule-update"
	CategoryParticipantStudyAssessmentLaunched    Category = "participant-study-schedule-assessment-launch"
	CategoryParticipantMobileAppCredentialsUpdate Category = "participant-mobile-app-credentials-update"

	CategoryParticipantVendorAccountManagementCreate Category = "participant-vendor-account-management-create"
	CategoryParticipantVendorAccountManagementUpdate Category = "participant-vendor-account-management-update"
	CategoryParticipantVendorDeviceAssignmentCreate  Category = "participant-vendor-device-assignment-create"
	CategoryParticipantVendorDeviceAssignmentUpdate  Category = "participant-vendor-device-assignment-update"

	CategoryUserDetailsCreate            Category = "user-details-create"
	CategoryUserDetailsUpdate            Category = "user-details-update"
	CategoryUserLogEventsSessionCreate   Category = "user-log-events-session-create"
	CategoryUserLogEventsSessionDelete   Category = "user-log-events-session-delete"
	CategoryUserLogEventsPasswordChanged Category = "user-log-events-password-changed"
	CategoryUserLogEventsStatusChanged   Category = "user-status-change"
)

var ErrorCorrection = Reason{
	ID:          "error-correction",
	DisplayName: "Error Correction",
}

var ParticipantModification = Reason{
	ID:          "participant-modification",
	DisplayName: "Participant Modification",
}

var StudyStaffModification = Reason{
	ID:          "study-staff-modification",
	DisplayName: "Study Staff Modification",
}

var TechnicalReasonModification = Reason{
	ID:          "technical-reason-modification",
	DisplayName: "Technical Reason Modification",
}

var DataModification = Reason{
	ID:          "data-modification",
	DisplayName: "Data Modification",
}

var ParticipantOnboarding = Reason{
	ID:          "participant-onboarding",
	DisplayName: "Participant Onboarding",
}

var AssessmentManuallyLaunched = Reason{
	ID:          "assessment-manually-launched",
	DisplayName: "Assessment Manually Launched",
}

var StudyParticipationConcluded = Reason{
	ID:          "study-participation-concluded",
	DisplayName: "Study Participant Concluded",
}

var ExternalAccountCreation = Reason{
	ID:          "external-account-creation",
	DisplayName: "External Account Creation",
}

var DeviceAssignment = Reason{
	ID:          "device-assignment",
	DisplayName: "Device Assignment",
}

var DeviceChange = Reason{
	ID:          "device-change",
	DisplayName: "Device Change",
}

var DeviceMalfunction = Reason{
	ID:          "device-malfunction",
	DisplayName: "Device Malfunction",
}

var DeviceNoLongerNeeded = Reason{
	ID:          "device-no-longer-needed",
	DisplayName: "Device No Longer Needed",
}

var PasscodeReset = Reason{
	ID:          "passcode-reset",
	DisplayName: "Passcode Reset",
}

var LoginInformationExpired = Reason{
	ID:          "login-information-expired",
	DisplayName: "Login Information Expired",
}

var UserLogin = Reason{
	ID:          "user-login",
	DisplayName: "User Login",
}

var UserLogout = Reason{
	ID:          "user-logout",
	DisplayName: "User Logout",
}

var UserPasswordUpdate = Reason{
	ID:          "user-password-update",
	DisplayName: "User Password Update",
}

var UserStatusChange = Reason{
	ID:          "user-status-change",
	DisplayName: "User Status Change",
}

var InitialCreation = Reason{
	ID:          "initial-creation",
	DisplayName: "Initial Creation",
}

var InitialUserAccountCreation = Reason{
	ID:          "initial-user-account-creation",
	DisplayName: "Initial User Account Creation",
}

var ReportRequested = Reason{
	ID:          "report-requested",
	DisplayName: "Report Requested",
}

var ReasonsForChangeByCategoryAndAction = map[Category]ReasonList{
	CategoryParticipantStatusAdvance: {
		ParticipantOnboarding,
		ParticipantModification,
		StudyStaffModification,
	},
	CategoryParticipantStatusRevert: {
		ErrorCorrection,
		TechnicalReasonModification,
	},
	CategoryParticipantStatusComplete: {
		StudyParticipationConcluded,
	},
	CategoryParticipantStatusRemove: {
		StudyStaffModification,
	},
	CategoryParticipantDetailsCreate: {
		ParticipantOnboarding,
	},
	CategoryParticipantDetailsUpdate: {
		ErrorCorrection,
		ParticipantModification,
	},
	CategoryParticipantStudyScheduleCreate: {
		ParticipantOnboarding,
	},
	CategoryParticipantStudyScheduleUpdate: {
		ErrorCorrection,
		ParticipantModification,
		StudyStaffModification,
		TechnicalReasonModification,
	},
	CategoryParticipantStudyAssessmentLaunched: {
		AssessmentManuallyLaunched,
	},
	CategoryGenericCreation: {
		InitialCreation,
	},
	CategoryGenericUpdate: {
		DataModification,
		ErrorCorrection,
	},

	CategoryUserDetailsCreate: {
		InitialUserAccountCreation,
	},
	CategoryUserDetailsUpdate: {
		DataModification,
	},
	CategoryUserLogEventsSessionCreate: {
		UserLogin,
	},
	CategoryUserLogEventsSessionDelete: {
		UserLogout,
	},
	CategoryUserLogEventsPasswordChanged: {
		UserPasswordUpdate,
	},
	CategoryUserLogEventsStatusChanged: {
		UserStatusChange,
	},
	CategoryParticipantMobileAppCredentialsUpdate: {
		DeviceAssignment,
		DeviceChange,
		DeviceMalfunction,
		PasscodeReset,
		LoginInformationExpired,
	},
	CategoryParticipantVendorAccountManagementCreate: {
		ExternalAccountCreation,
	},
	CategoryParticipantVendorAccountManagementUpdate: {
		ParticipantModification,
		ErrorCorrection,
	},
	CategoryParticipantVendorDeviceAssignmentCreate: {
		DeviceAssignment,
	},
	CategoryParticipantVendorDeviceAssignmentUpdate: {
		DeviceAssignment,
		DeviceChange,
		DeviceMalfunction,
		DeviceNoLongerNeeded,
	},
}

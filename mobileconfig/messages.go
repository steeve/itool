package mobileconfig

type RequestBase struct {
	RequestType string `plist:"RequestType"`
}

type ResponseBase struct {
	Status     string `plist:"Status"`
	ErrorChain []struct {
		ErrorCode            int    `plist:"ErrorCode"`
		ErrorDomain          string `plist:"ErrorDomain"`
		LocalizedDescription string `plist:"LocalizedDescription"`
		USEnglishDescription string `plist:"USEnglishDescription"`
	} `plist:"ErrorChain"`
}

type InstallProfileRequest struct {
	RequestBase
	Payload []byte `plist:"Payload"`
}

type InstallProfileResponse struct {
	ResponseBase
}

type GetProfileListRequest struct {
	RequestBase
}

type GetProfileListResponse struct {
	ResponseBase
	OrderedIdentifiers []string                    `plist:"OrderedIdentifiers"`
	ProfileManifest    map[string]*ProfileManifest `plist:"ProfileManifest"`
	ProfileMetadata    map[string]*ProfileMetadata `plist:"ProfileMetadata"`
}

type RemoveProfilePayload struct {
	PayloadType       string `plist:"PayloadType"`
	PayloadIdentifier string `plist:"PayloadIdentifier"`
	PayloadUUID       string `plist:"PayloadUUID"`
	PayloadVersion    int    `plist:"PayloadVersion"`
}

type RemoveProfileRequest struct {
	RequestBase
	ProfileIdentifier []byte `plist:"ProfileIdentifier"`
}

type RemoveProfileResponse struct {
	ResponseBase
}

type ProfileManifest struct {
	Description string `plist:"Description"`
	IsActive    bool   `plist:"IsActive"`
}

type ProfileMetadata struct {
	PayloadDescription       string `plist:"PayloadDescription"`
	PayloadDisplayName       string `plist:"PayloadDisplayName"`
	PayloadOrganization      string `plist:"PayloadOrganization"`
	PayloadRemovalDisallowed bool   `plist:"PayloadRemovalDisallowed"`
	PayloadUUID              string `plist:"PayloadUUID"`
	PayloadVersion           int    `plist:"PayloadVersion"`
}

type Profile struct {
	Identifier string
	Manifest   *ProfileManifest
	Metadata   *ProfileMetadata
}

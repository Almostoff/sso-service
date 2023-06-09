package model

type KYCResponse struct {
	ClientId    string `json:"clientId"`
	ScanRef     string `json:"scanRef"`
	ExternalRef string `json:"externalRef"`
	Platform    string `json:"platform"`
	StartTime   int    `json:"startTime"`
	FinishTime  int    `json:"finishTime"`
	Status      struct {
		Overall          string        `json:"overall"`
		SuspicionReasons []interface{} `json:"suspicionReasons"`
		MismatchTags     []interface{} `json:"mismatchTags"`
		AutoDocument     string        `json:"autoDocument"`
		AutoFace         string        `json:"autoFace"`
		ManualDocument   string        `json:"manualDocument"`
		ManualFace       string        `json:"manualFace"`
	} `json:"status"`
	Data struct {
		SelectedCountry     string      `json:"selectedCountry"`
		DocFirstName        string      `json:"docFirstName"`
		DocLastName         string      `json:"docLastName"`
		DocNumber           string      `json:"docNumber"`
		DocPersonalCode     string      `json:"docPersonalCode"`
		DocExpiry           string      `json:"docExpiry"`
		DocDob              string      `json:"docDob"`
		DocType             string      `json:"docType"`
		DocSex              string      `json:"docSex"`
		DocNationality      string      `json:"docNationality"`
		DocIssuingCountry   string      `json:"docIssuingCountry"`
		ManuallyDataChanged bool        `json:"manuallyDataChanged"`
		OrgFirstName        string      `json:"orgFirstName"`
		OrgLastName         string      `json:"orgLastName"`
		OrgNationality      string      `json:"orgNationality"`
		OrgBirthPlace       string      `json:"orgBirthPlace"`
		OrgAuthority        interface{} `json:"orgAuthority"`
		OrgAddress          interface{} `json:"orgAddress"`
	} `json:"data"`
	FileUrls struct {
		FRONT string `json:"FRONT"`
		BACK  string `json:"BACK"`
		FACE  string `json:"FACE"`
	} `json:"fileUrls"`
	AML []struct {
		Status struct {
			ServiceSuspected bool   `json:"serviceSuspected"`
			CheckSuccessful  bool   `json:"checkSuccessful"`
			ServiceFound     bool   `json:"serviceFound"`
			ServiceUsed      bool   `json:"serviceUsed"`
			OverallStatus    string `json:"overallStatus"`
		} `json:"status"`
		Data             []interface{} `json:"data"`
		ServiceName      string        `json:"serviceName"`
		ServiceGroupType string        `json:"serviceGroupType"`
		Uid              string        `json:"uid"`
		ErrorMessage     interface{}   `json:"errorMessage"`
	} `json:"AML"`
	LID []struct {
		Status struct {
			ServiceSuspected bool   `json:"serviceSuspected"`
			CheckSuccessful  bool   `json:"checkSuccessful"`
			ServiceFound     bool   `json:"serviceFound"`
			ServiceUsed      bool   `json:"serviceUsed"`
			OverallStatus    string `json:"overallStatus"`
		} `json:"status"`
		Data             []interface{} `json:"data"`
		ServiceName      string        `json:"serviceName"`
		ServiceGroupType string        `json:"serviceGroupType"`
		Uid              string        `json:"uid"`
		ErrorMessage     interface{}   `json:"errorMessage"`
	} `json:"LID"`
}

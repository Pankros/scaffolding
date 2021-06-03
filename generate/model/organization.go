package model

type OrganizationOutput struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	SiteID string `json:"site_id"`
}

type OrganizationCreate struct {
	Name   string `json:"name"`
	SiteID string `json:"site_id"`
}

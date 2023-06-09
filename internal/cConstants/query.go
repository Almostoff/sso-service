package cConstants

const (
	GetInnerConnectionQuery string = "SELECT * FROM inner_connection WHERE name = $1"
	GetServiceByPublicQuery string = "SELECT * FROM inner_connection WHERE public = $1"
)

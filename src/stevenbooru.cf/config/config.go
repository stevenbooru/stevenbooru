package config

// Config is the parent configuration for Stevenbooru.
type Config struct {
	Database      Database
	Storage       Storage
	Elasticsearch Elasticsearch
	Redis         Redis
	HTTP          HTTP
	Site          Site
	Roles         map[string]Role
}

// Database is the database configuration.
type Database struct {
	Kind     string
	Username string
	Password string
	Database string
	Host     string
	Port     int
}

// Storage is the storage configuration for Stevenbooru uploaded content.
type Storage struct {
	Kind string
	Path string
}

// Elasticsearch is the configuration for SB elasticsearch usage.
type Elasticsearch struct {
	Host  string
	Index string
}

// Redis is the redis configuration for SB.
type Redis struct {
	Host string
	Port string
}

// HTTP is the http configuration for SB.
type HTTP struct {
	Bindhost string
	Port     string
}

// Site contains other information about Stevenbooru.
type Site struct {
	ShowAds bool     // Show advertisments?
	Pepper  string   // String to pepper password hashes with
	Roles   []string // Site roles to load and evaluate
}

/*
Role defines a role a Stevenbooru user can have.

A sane set of defaults follows:

	; Administrators have all access to everything.
	[role "admin"]
	name=Administrator
	canhide=true
	canharddelete=true
	canban=true
	canpermaban=true
	canseerawip=true
	canadmintags=true
	canlock=true
	canuserlink=true
	canedit=true
	siteadmin=true

	; Moderators are like administrators, but without some devastating effects.
	[role "moderator"]
	name=Moderator
	canhide=true
	canharddelete=false
	canban=true
	canpermaban=false
	canseerawip=true
	canadmintags=true
	canlock=true
	canuserlink=true
	canedit=true
	siteadmin=false

	; Assistants have the least amount of permissions allocated to them.
	[role "assistant"]
	name=Assistant
	canhide=true
	canharddelete=false
	canban=false
	canpermaban=false
	canseerawip=false
	canadmintags=true
	canlock=true
	canuserlink=true
	canedit=true
	siteadmin=false
*/
type Role struct {
	Name          string // Human-readable role name
	CanHide       bool   // Can hide things?
	CanHardDelete bool   // Can hard-delete things?
	CanBan        bool   // Can place bans?
	CanPermaban   bool   // Can place permanent bans?
	CanSeeRawIP   bool   // Can see raw IP data? If false it will show a cryptographic hash.
	CanAdminTags  bool   // Can administer tags?
	CanLock       bool   // Can lock things from additional comments/posts?
	CanUserLink   bool   // Can manage user links?
	CanEdit       bool   // Can this user edit things they otherwise couldn't?
	SiteAdmin     bool   // Can this user see site admin areas?
}

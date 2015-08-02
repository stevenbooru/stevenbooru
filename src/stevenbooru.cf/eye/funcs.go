package eye

import (
	"html/template"

	"stevenbooru.cf/globals"
)

var (
	funcs template.FuncMap
)

func init() {
	funcs = template.FuncMap{
		"Can": Can, // can you do the can-can?
	}
}

func Can(role, permission string) bool {
	perms, ok := globals.Config.Role[role]
	if !ok {
		return false
	}

	switch permission {
	case "canhide":
		return perms.CanHide
	case "canharddelete":
		return perms.CanHardDelete
	case "canban":
		return perms.CanBan
	case "canpermaban":
		return perms.CanPermaban
	case "canseerawip":
		return perms.CanSeeRawIP
	case "canadmintags":
		return perms.CanAdminTags
	case "canlock":
		return perms.CanLock
	case "canuserlink":
		return perms.CanUserLink
	case "canedit":
		return perms.CanEdit
	case "siteadmin":
		return perms.SiteAdmin
	default:
		return false
	}
}

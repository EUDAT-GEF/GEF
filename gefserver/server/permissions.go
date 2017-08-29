package server

import (
	"net/http"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
)

type Authorization struct {
	s *Server
	w http.ResponseWriter
	r *http.Request
}

// allowSuperAdmin returns the logged in user or writes errors into the http stream
// First returned value is true if the user is superadmin
func (a Authorization) getUserInfo() (bool, *db.User) {
	user := a.s.getUserOrWriteError(a.w, a.r)
	if user == nil {
		return false, nil
	}
	return a.s.isSuperAdmin(user), user
}

func (a Authorization) allowCreateToken() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// allow any logged in user to create their own tokens
	allow = true
	return
}

func (a Authorization) allowGetTokens() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// allow any user to list their own tokens
	allow = true
	return
}

func (a Authorization) allowDeleteToken(token db.Token) (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	if user.ID == token.UserID {
		// allow logged in user which owns the token to delete it
		allow = true
		return
	}
	Response{a.w}.Forbidden("This token is owned by someone else")
	return
}

func (a Authorization) allowListRoles() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// allow any logged in user to list roles
	allow = true
	return
}

func (a Authorization) allowListRoleUsers() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// only superadmins can list role users
	Response{a.w}.Forbidden("Only superadministrators can list role users")
	return
}

func (a Authorization) allowNewRoleUser() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// only superadmins can create role user
	Response{a.w}.Forbidden("Only superadministrators can create role user")
	return
}

func (a Authorization) allowDeleteRoleUser() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// only superadmins can delete role user
	Response{a.w}.Forbidden("Only superadministrators can delete role user")
	return
}

func (a Authorization) allowCreateBuild() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// only superadmins can create builds
	Response{a.w}.Forbidden("Only superadministrators can create builds")
	return
}

func (a Authorization) allowUploadIntoBuild() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// only superadmins can upload data into builds
	Response{a.w}.Forbidden("Only superadministrators can create services")
	return
}

func (a Authorization) allowListServices() (allow bool, user *db.User) {
	allow = true // anybody can see the list of services
	return
}

func (a Authorization) allowInspectService(serviceID db.ServiceID) (allow bool, user *db.User) {
	allow = true // anybody can inspect a service
	return
}

func (a Authorization) allowEditService(serviceID db.ServiceID) (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	if a.s.db.IsServiceOwner(user.ID, serviceID) {
		allow = true // a service's owner can remove the job
		return
	}
	Response{a.w}.Forbidden("A service can only be edited by its owner")
	return
}

func (a Authorization) allowRemoveService(serviceID db.ServiceID) (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	if a.s.db.IsServiceOwner(user.ID, serviceID) {
		allow = true // a service's owner can remove the job
		return
	}
	Response{a.w}.Forbidden("A service can only be removed by its owner")
	return
}

func (a Authorization) allowCreateJob() (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	// roles, err := a.s.db.GetUserRoles(user.ID)
	// if err == nil {
	// 	for _, r := range roles {
	// 		if r.Name == db.CommunityMemberRoleName {
	// 			allow = true // community members can create jobs
	// 			return
	// 		}
	// 	}
	// }
	// Response{a.w}.Forbidden("Only community members can create jobs")
	allow = true // anyone can create jobs (for the moment?)
	return
}

func (a Authorization) allowListJobs() (allow bool, user *db.User) {
	allow = true // anybody can see the list of jobs
	return
}

func (a Authorization) allowInspectJob(jobID db.JobID) (allow bool, user *db.User) {
	allow = true // anybody can inspect a job (metadata only)
	return
}

func (a Authorization) allowRemoveJob(jobID db.JobID) (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	if a.s.db.IsJobOwner(user.ID, jobID) {
		allow = true // a job's owner can remove the job
		return
	}
	Response{a.w}.Forbidden("A job can only be removed by its owner")
	return
}

func (a Authorization) allowGetJobData(jobID db.JobID) (allow bool, user *db.User) {
	allow, user = a.getUserInfo()
	if user == nil || allow {
		return
	}
	if a.s.db.IsJobOwner(user.ID, jobID) {
		allow = true // a job's owner can retrieve the job data
		return
	}
	Response{a.w}.Forbidden("Job data can only be retrieved by its owner")
	return
}

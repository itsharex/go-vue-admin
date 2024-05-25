package system

import "github/shansec/go-vue-admin/utils"

var (
	CreateRoleVerify = utils.Rules{"RoleId": {utils.NotEmpty()}, "RoleName": {utils.NotEmpty()}}
	DeleteRoleVerify = utils.Rules{"RoleId": {utils.NotEmpty()}}
)

package role

import (
	"b2b.nati011.github.com/internal/core/application/role"
)

func allRolesMapper(in role.GetAllResponse) GetAllRoleResponse {
	var response GetAllRoleResponse
	for _, i := range in.List {
		response.List = append(response.List, (GetRoleResponse)(i))
	}
	return response
}

func allResourcesMapper(in role.GetAllResourcesResponse) GetAllResourcesResponse {
	var response GetAllResourcesResponse
	response.List = append(response.List, in.List...)
	return response
}

package policies

import (
	"backend-brevet/models"
)

// CanBeAssignedAsTeacher checks if user has role "teacher"
func CanBeAssignedAsTeacher(user *models.User) bool {
	return user.RoleType == models.RoleTypeGuru
}

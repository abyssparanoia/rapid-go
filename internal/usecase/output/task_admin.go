package output

type TaskCreateAdmin struct {
	AdminID  string
	AuthUID  string
	Password string //nolint:gosec
}

func NewTaskCreateAdmin(
	adminID string,
	authUID string,
	password string,
) *TaskCreateAdmin {
	return &TaskCreateAdmin{
		AdminID:  adminID,
		AuthUID:  authUID,
		Password: password,
	}
}

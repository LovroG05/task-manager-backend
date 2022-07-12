package inputs

type TaskInput struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Time         string `json:"time" binding:"required"`
	Days         uint8  `json:"days" binding:"required"`
	OneTime      bool   `json:"one_time"`
	AssigneesIDs []uint `json:"assignees_ids"`
}

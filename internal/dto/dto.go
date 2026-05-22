package dto

type GetClassListPayload struct {
}

type LogginPayload struct {
	StudentID      string `json:"student_id"`
	SuccessClassID []int  `json:"success_class_id"`
}

type RegisterPayload struct {
	ClassID    string `json:"class_id"`
	StudentID  string `json:"student_id"`  //Need to be array
	ResponseCh string `json:"response_ch"` // identifier / correlation ID for the response
}

package dto

type GetClassListPayload struct {
	StudentID         string `bson:"student_id"`
	RegisteredClassID []int  `json:"success_class_id"`
}

type LogginPayload struct {
	Action         int    `json:"action" bson:"-"`
	StudentID      string `json:"student_id" bson:"student_id"`
	SuccessClassID []int  `json:"success_class_id" bson:"success_class_id"`
}

type RegisterPayload struct {
	ClassID    string `json:"class_id"`
	StudentID  string `json:"student_id"`  //Need to be array
	ResponseCh string `json:"response_ch"` // identifier / correlation ID for the response
}

// type Unregister struct {
// 	StudentID  string `json:"student_id" bson:"student_id"`
// 	Unregister []int  `json:"unregister_class_id"`
// }

package ssg

type ContentForm struct {
	Heading string `form:"heading" required:"true"`
	Body    string `form:"body"`
	Status  string `form:"status"`
}

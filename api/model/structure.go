package model

type Structure struct {
	ID         int       `json:"id"`
	ParentID   *int      `json:"parentId" db:"parent_id"`
	Title      string    `json:"title"`
	SyllabusID *int      `json:"-" db:"syllabus_id"`
	Syllabus   *Syllabus `json:"syllabus" db:"-"`
}

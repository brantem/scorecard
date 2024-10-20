package model

type SyllabusStructure struct {
	ID     int    `json:"id"`
	PrevID *int   `json:"prevId" db:"prev_id"`
	Title  string `json:"title"`
}

type Syllabus struct {
	ID          int    `json:"id"`
	ParentID    *int   `json:"parentId" db:"parent_id"`
	StructureID *int   `json:"structureId,omitempty" db:"structure_id"`
	Title       string `json:"title"`
}

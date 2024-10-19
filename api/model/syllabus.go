package model

type SyllabusStructure struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Syllabus struct {
	ID          int    `json:"id"`
	StructureID *int   `json:"structureId,omitempty" db:"structure_id"`
	Title       string `json:"title"`
}

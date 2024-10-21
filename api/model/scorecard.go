package model

type ScorecardStructure struct {
	ID         int                         `json:"id"`
	ParentID   *int                        `json:"parentId" db:"parent_id"`
	Title      string                      `json:"title"`
	SyllabusID *int                        `json:"-" db:"syllabus_id"`
	Syllabus   *ScorecardStructureSyllabus `json:"syllabus" db:"-"`
}

type ScorecardStructureSyllabus struct {
	BaseSyllabus
	IsAssignment bool `json:"isAssignment" db:"is_assignment"`
}

type Scorecard struct {
	ID          int              `json:"id"`
	UserID      int              `json:"-" db:"user_id"`
	User        *User            `json:"user" db:"-"`
	Score       float64          `json:"score"`
	Items       []*ScorecardItem `json:"items"`
	IsOutdated  bool             `json:"isOutdated" db:"is_outdated"`
	GeneratedAt Time             `json:"generatedAt" db:"generated_at"`
}

type ScorecardItem struct {
	ID          int     `json:"id"`
	StructureID int     `json:"structureId" db:"structure_id"`
	Score       float64 `json:"score"`
}

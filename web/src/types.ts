export type User = {
  id: number;
  name: string;
};

export type SyllabusStructure = {
  id: number;
  prevId: number | null;
  title: string;
};

export type BaseSyllabus = {
  id: number;
  title: string;
};

export type Syllabus = BaseSyllabus & {
  parentId: number | null;
  structureId?: number;
};

export type Score = {
  user: User;
  score: number | null;
};

export type Structure = {
  id: number;
  parentId: number | null;
  title: string;
};

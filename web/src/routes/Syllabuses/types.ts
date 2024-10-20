export type Structure = {
  id: number;
  prevId: number | null;
  title: string;
};

export type Syllabus = {
  id: number;
  parentId: number | null;
  structureId: number;
  title: string;
};

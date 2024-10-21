import type { BaseSyllabus } from './syllabus';

type Syllabus = BaseSyllabus & {
  isAssignment: boolean;
};

export type Structure = {
  id: number;
  parentId: number | null;
  title: string;
  syllabus: Syllabus | null;
};

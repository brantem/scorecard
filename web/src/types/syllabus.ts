import type { User } from './user';

export type Structure = {
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

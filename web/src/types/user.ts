import type { BaseSyllabus } from './syllabus';

export type User = {
  id: number;
  name: string;
};

export type Score = {
  syllabus: BaseSyllabus & {
    parents: BaseSyllabus[];
  };
  score: number | null;
};

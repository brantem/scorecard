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

export type Scorecard = {
  id: number;
  user: {
    id: number;
    name: string;
  };
  score: number;
  items: ScorecardItem[];
  isOutdated: boolean;
  generatedAt: string;
};

export type ScorecardItem = {
  id: number;
  structureId: number;
  score: number;
};

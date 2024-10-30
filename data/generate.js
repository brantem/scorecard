const { faker } = require('@faker-js/faker');
const fs = require('fs');

const path = 'data.sql';

if (fs.existsSync(path)) fs.unlinkSync(path);

const PROGRAMS = 10;
const USERS = 10;

const SYLLABUS_STRUCTURES_MIN = 1;
const SYLLABUS_STRUCTURES_MAX = 3;

const SYLLABUSES_MIN = 2;
const SYLLABUSES_MAX = 3;

// programs
for (let i = 1; i <= PROGRAMS; i++) {
  const titles = [
    faker.food.fruit(),
    faker.food.ingredient(),
    faker.food.meat(),
    faker.food.spice(),
    faker.food.vegetable(),
  ];
  const title = titleCase(titles[Math.floor(Math.random() * titles.length)]);
  fs.appendFileSync(path, `INSERT INTO programs (id, title) VALUES (${i}, '${title}');\n`);
}

fs.appendFileSync(path, '\n');

// users
const userIdsByProgramId = {};

let users_id = 1;
for (let i = 1; i <= PROGRAMS; i++) {
  for (let j = 1; j <= USERS; j++) {
    const id = users_id++;
    const name = titleCase(faker.person.fullName().replace("'"));
    fs.appendFileSync(path, `INSERT INTO users (id, program_id, name) VALUES (${id}, ${i}, '${name}');\n`);

    userIdsByProgramId[i] = [...(userIdsByProgramId[i] || []), id];
  }
}

fs.appendFileSync(path, '\n');

// syllabus_structures
const structureIdsByProgramId = {};
const assignmentStructureIdByProgramId = {};

let syllabus_structures_id = 1;
for (let i = 1; i <= PROGRAMS; i++) {
  const STRUCTURES = getRandomNumber(SYLLABUS_STRUCTURES_MIN, SYLLABUS_STRUCTURES_MAX);

  const id = syllabus_structures_id++;
  fs.appendFileSync(
    path,
    `INSERT INTO syllabus_structures (id, program_id, prev_id, title) VALUES (${id}, ${i}, -1, 'Assignment');\n`,
  );
  structureIdsByProgramId[i] = [...(structureIdsByProgramId[i] || []), id];
  assignmentStructureIdByProgramId[i] = id;

  for (let j = 0; j < STRUCTURES; j++) {
    const id = syllabus_structures_id++;
    const prevId = j === 0 ? 'NULL' : syllabus_structures_id - 2;
    const title = titleCase(faker.lorem.word());
    fs.appendFileSync(
      path,
      `INSERT INTO syllabus_structures (id, program_id, prev_id, title) VALUES (${id}, ${i}, ${prevId}, '${title}');\n`,
    );

    structureIdsByProgramId[i] = [...(structureIdsByProgramId[i] || []), id];
  }
}

fs.appendFileSync(path, '\n');

// syllabuses
const assignmentIdsByProgramId = {};

let syllabuses_id = 1;
for (const programId in structureIdsByProgramId) {
  const SYLLABUSES = getRandomNumber(SYLLABUSES_MIN, SYLLABUSES_MAX);
  let parentIds = new Array(SYLLABUSES).fill('NULL');

  let structures = structureIdsByProgramId[programId];
  structures = [...structures.slice(1), structures[0]];

  for (let i = 0; i < structures.length; i++) {
    const structureId = structures[i];

    let tempParentIds = [];
    for (const parentId of parentIds) {
      for (let j = 0; j < (parentId === 'NULL' ? 1 : SYLLABUSES); j++) {
        const id = syllabuses_id++;
        const title = titleCase(faker.lorem.words({ min: 1, max: 5 }));
        fs.appendFileSync(
          path,
          `INSERT INTO syllabuses (id, parent_id, structure_id, title) VALUES (${id}, ${parentId}, ${structureId}, '${title}');\n`,
        );

        if (assignmentStructureIdByProgramId[programId] === structureId) {
          assignmentIdsByProgramId[programId] = [...(assignmentIdsByProgramId[programId] || []), id];
        }
        tempParentIds.push(id);
      }
    }
    parentIds = tempParentIds;
  }
}

fs.appendFileSync(path, '\n');

// user_scores
for (const programId in assignmentIdsByProgramId) {
  for (const syllabusId of assignmentIdsByProgramId[programId]) {
    for (const userId of userIdsByProgramId[programId]) {
      const score = getRandomNumber(0, 100);
      fs.appendFileSync(
        path,
        `INSERT INTO user_scores (user_id, syllabus_id, score) VALUES (${userId}, ${syllabusId}, ${score});\n`,
      );
    }
  }
}

function titleCase(s) {
  return s.replace(/\w\S*/g, (text) => text.charAt(0).toUpperCase() + text.substring(1).toLowerCase());
}

function getRandomNumber(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

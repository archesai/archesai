import "tsconfig-paths/register";

import { resetDatabase } from "../prisma/seed";

module.exports = async () => {
  await resetDatabase(); // Reset the database before tests run
};

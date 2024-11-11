import { resetDatabase } from "../prisma/seed";

module.exports = async () => {
  await resetDatabase(); // Reset the database before tests run
};

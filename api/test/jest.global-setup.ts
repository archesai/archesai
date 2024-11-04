import { resetDatabase } from "./util";

module.exports = async () => {
  await resetDatabase(); // Reset the database before tests run
};

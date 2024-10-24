export default async () => {
  await fetch(
    "http://arches-firebase:9099/emulator/v1/projects/filechat-io/accounts",
    {
      headers: {
        Authorization: "Bearer owner",
      },
      method: "DELETE",
    }
  );
};

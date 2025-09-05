export function formatDate(
  date: Date | number | string | undefined,
  opts: Intl.DateTimeFormatOptions = {}
): string {
  if (!date) return ""

  try {
    return new Intl.DateTimeFormat("en-US", {
      day: opts.day ?? "numeric",
      month: opts.month ?? "long",
      year: opts.year ?? "numeric",
      ...opts
    }).format(new Date(date))
  } catch {
    return ""
  }
}

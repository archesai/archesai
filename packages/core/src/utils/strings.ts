/**
 * Converts a string into "camelCase".
 *
 * Example: "API Token Example" => "apiTokenExample"
 */
export function toCamelCase(str: string): string {
  const words = str
    // Replace any non-alphanumeric with space
    .replace(/[^a-zA-Z0-9]+/g, ' ')
    .trim()
    .split(/\s+/)
  if (words.length === 0) return ''

  // Lowercase the first word entirely
  const first = words[0]!.toLowerCase()

  // Capitalize first letter of each subsequent word and lowercase the rest
  const rest = words
    .slice(1)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())

  return [first, ...rest].join('')
}

/**
 * Converts a string into "kebab-case" (dash-separated).
 *
 * Example: "Api Token Example" => "api-token-example"
 */
export function toKebabCase(str: string): string {
  return (
    str
      // Replace any non-alphanumeric/underscore/dash with a space
      .replace(/[^a-zA-Z0-9]+/g, ' ')
      .trim()
      // Split by spaces, lowercase, then join with dashes
      .split(/\s+/)
      .map((word) => word.toLowerCase())
      .join('-')
  )
}

/**
 * Transforms a string into "Sentence Case":
 * 1. Replaces underscores with spaces.
 * 2. Inserts a space before uppercase letters.
 * 3. Converts to lowercase.
 * 4. Removes extra spaces.
 * 5. Capitalizes each word.
 *
 * Example: "api_TOKENExample" => "Api Token Example"
 */
export function toSentenceCase(str: string): string {
  return str
    .replace(/_/g, ' ') // Replace underscores with spaces
    .replace(/([A-Z])/g, ' $1') // Add space before capital letters
    .toLowerCase() // Convert the string to lowercase
    .replace(/\s+/g, ' ') // Remove extra spaces
    .trim() // Remove leading/trailing spaces
    .split(' ') // Split into words
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1)) // Capitalize each word
    .join(' ')
}

/**
 * Converts a string into "snake_case" (underscore-separated).
 *
 * Example: "Api Token Example" => "api_token_example"
 */
export function toSnakeCase(str: string): string {
  return (
    str
      // Replace any non-alphanumeric/underscore/dash with a space
      .replace(/[^a-zA-Z0-9]+/g, ' ')
      .trim()
      // Split by spaces, lowercase, then join with underscores
      .split(/\s+/)
      .map((word) => word.toLowerCase())
      .join('_')
  )
}

/**
 * Converts a string into "Title Case" (each word capitalized, rest lower).
 *
 * Example: "my small example" => "My Small Example"
 */
export function toTitleCase(str: string): string {
  return (
    str
      // Replace underscores/dashes with spaces
      .replace(/[_-]+/g, ' ')
      .trim()
      .split(/\s+/)
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ')
  )
}

/**
 * Converts a string into "TitleCase" (each word capitalized, rest lower).
 *
 * Example: "my small example" => "My SmallExample"
 */
export function toTitleCaseNoSpaces(str: string): string {
  return (
    str
      // Replace underscores/dashes with spaces
      .replace(/[_-]+/g, ' ')
      .trim()
      .split(/\s+/)
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join('')
  )
}

/**
 * Determines if the given word starts with a vowel and returns 'n' if it does.
 * Otherwise, returns an empty string.
 *
 * Examples:
 *    vf("apple") => "n"
 *    vf("banana") => ""
 */
export function vf(word: string): string {
  const char = word.charAt(0).toLowerCase()
  if (['a', 'e', 'i', 'o', 'u'].includes(char)) {
    return 'n'
  }
  return ''
}

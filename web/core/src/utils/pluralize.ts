export const pluralize = (word: string, count: number): string => {
  return count === 1 ? word : `${word}s`
}

export const singularize = (word: string): string => {
  return word.endsWith('s') ? word.slice(0, -1) : word
}

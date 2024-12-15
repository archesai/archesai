import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function toSentenceCase(str: string) {
  return str
    .replace(/_/g, ' ') // Replace underscores with spaces
    .replace(/([A-Z])/g, ' $1') // Add space before capital letters
    .toLowerCase() // Convert the string to lowercase
    .replace(/\s+/g, ' ') // Remove extra spaces
    .trim() // Remove leading/trailing spaces
    .split(' ') // Split the string into words
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1)) // Capitalize each word
    .join(' ') // Join the words back into a single string
}

export function toCapitalized(str: string): string {
  return str.replace(/\b\w/g, (char) => char.toUpperCase())
}

import type { ClassValue } from 'clsx'

import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function stringToColor(str: string) {
  if (str.includes('video')) {
    return 'text-green-600'
  } else if (
    str.includes('audio') ||
    str.includes('music') ||
    str.includes('speech')
  ) {
    return 'text-purple-600'
  } else if (str.includes('image')) {
    return 'text-blue-900'
  } else if (str.includes('application')) {
    return 'text-blue-800'
  } else if (str.includes('text')) {
    return 'text-green-600'
  } else {
    let hash = 0
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash)
    }
    let color = '#'
    for (let i = 0; i < 3; i++) {
      const value = (hash >> (i * 8)) & 0xff
      color += ('00' + value.toString(16)).slice(-2)
    }
    return color
  }
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

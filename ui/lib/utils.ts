import { ContentControllerFindAllResponse } from '@/generated/archesApiComponents'
import { queryKeyFn } from '@/generated/archesApiContext'
import { ContentEntity } from '@/generated/archesApiSchemas'
import { QueryClient } from '@tanstack/react-query'
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

export const streamContent = (
  orgname: string,
  labelId: string,
  content: ContentEntity,
  queryClient: QueryClient
) => {
  queryClient.setQueryData(
    queryKeyFn({
      operationId: 'contentControllerFindAll',
      path: '/organizations/{orgname}/content',
      variables: {
        pathParams: {
          orgname: orgname
        },
        queryParams: {
          sortBy: 'createdAt',
          sortDirection: 'desc'
        }
      }
    }),
    (oldData: ContentControllerFindAllResponse) => {
      if (!oldData) {
        oldData = {
          aggregates: [],
          metadata: { limit: 100, offset: 0, totalResults: 0 },
          results: []
        }
      }
      const prevStreamedMessage = oldData.results?.find(
        (i) => i.id === content.id
      )
      if (prevStreamedMessage) {
        return {
          ...oldData,
          results: [
            { ...prevStreamedMessage, answer: content.text },
            ...(oldData.results || [])
              .filter((i) => i.createdAt !== prevStreamedMessage?.createdAt)
              .filter((i) => i.id !== 'pending')
          ]
        }
      } else {
        return {
          ...oldData,
          results: [content, ...(oldData.results || [])]
        }
      }
    }
  )
}

import { atom } from 'jotai'

import type { SearchQuery } from '@archesai/core'
import type { BaseEntity } from '@archesai/schemas'

// Define atoms
export const searchQueryAtom = atom<SearchQuery<BaseEntity>>()

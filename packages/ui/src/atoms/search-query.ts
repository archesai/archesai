import { atom } from 'jotai'

import type { SearchQuery } from '@archesai/core'
import type { BaseEntity } from '@archesai/domain'

// Define atoms
export const searchQueryAtom = atom<SearchQuery<BaseEntity>>()

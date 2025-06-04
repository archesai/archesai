import type { ConfigArray } from 'typescript-eslint'

import base from '@archesai/eslint/base'
import jest from '@archesai/eslint/jest'

export default [...base, ...jest] satisfies ConfigArray

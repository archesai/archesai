import { render } from '@testing-library/react'

import Page from '#app/(app)/chat/page'

describe('Page', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<Page />)
    expect(baseElement).toBeTruthy()
  })
})

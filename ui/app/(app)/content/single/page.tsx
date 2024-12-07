import { Button } from '@/components/ui/button'
import { archesClient } from '@/lib/api/index'
import { QueryErrorResetBoundary } from '@tanstack/react-query'
import { ErrorBoundary } from 'react-error-boundary'

type ContentProps = {
  contentId: string | undefined
  defaultOrgname: string | undefined
}

const ContentSkeleton = () => (
  <div>
    <h1 style={{ backgroundColor: '#eee', height: '1.5em', width: '50%' }} />
    <p style={{ backgroundColor: '#eee', height: '1em', width: '80%' }} />
  </div>
)

const ContentComponent: React.FC<ContentProps> = ({ contentId, defaultOrgname }) => {
  const { data: content } = archesClient.useSuspenseQuery(
    'get',
    '/organizations/{orgname}/content/{id}',
    {
      params: {
        path: {
          id: contentId ?? '',
          orgname: defaultOrgname ?? ''
        }
      }
    },
    {
      // placeholderData: {
      //   description: '', // Mock description placeholder
      //   title: <ContentSkeleton /> // Mock title placeholder
      // }
    }
  )

  return (
    <div>
      <h1>{content?.title ?? 'No Title'}</h1>
      <p>{content?.description ?? 'No Description'}</p>
    </div>
  )
}

export const App = () => {
  const contentId = '123'
  const defaultOrgname = 'my-organization'

  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ErrorBoundary
          fallbackRender={({ resetErrorBoundary }) => (
            <div>
              There was an error!
              <Button onClick={() => resetErrorBoundary()}>Try again</Button>
            </div>
          )}
          onReset={reset}
        >
          <ContentComponent contentId={contentId} defaultOrgname={defaultOrgname} />
        </ErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  )
}

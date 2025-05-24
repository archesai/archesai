'use client'

import type { ChangeEvent, KeyboardEvent } from 'react'

import { useEffect, useRef, useState } from 'react'

import { createLabel, createRun, useFindManyContents } from '@archesai/client'
import { RefreshCcw } from '@archesai/ui/components/custom/icons'
import { Button } from '@archesai/ui/components/shadcn/button'
import { ScrollArea } from '@archesai/ui/components/shadcn/scroll-area'
import { toast } from '@archesai/ui/components/shadcn/sonner'
import { Textarea } from '@archesai/ui/components/shadcn/textarea'
import { useAuth } from '@archesai/ui/hooks/use-auth'
import { cn } from '@archesai/ui/lib/utils'

export default function Chat() {
  const [labelId, setLabelId] = useState<string>('')
  const { defaultOrgname } = useAuth()
  const [message, setMessage] = useState<string>('')

  const { data } = useFindManyContents({
    filter: {
      orgname: {
        equals: defaultOrgname
      }
    }
  })
  const messages = data?.data.data

  const messagesEndRef = useRef<HTMLDivElement | null>(null)
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'end'
      })
    }
  }, [messages])

  const handleSend = async () => {
    if (!message.trim()) return // Prevent sending empty messages

    if (!labelId) {
      try {
        const label = await createLabel({
          name: 'Chat'
        })
        setLabelId(label.data.data.id)
      } catch (error: unknown) {
        console.error(error)
        toast('Failed to create label')
      }
    }

    setMessage('')
    // streamContent(
    //   defaultOrgname,
    //   currentLabelId,
    //   {
    //     children: [],
    //     consumers: [],
    //     createdAt: new Date().toISOString(),
    //     credits: 0,
    //     description: 'Pending',
    //     id: 'pending',
    //     labels: [],
    //     mimeType: 'text/plain',
    //     name: 'Pending',
    //     orgname: defaultOrgname,
    //     parent: {
    //       id: '',
    //       name: 'Chat'
    //     },
    //     parentId: null,
    //     previewImage: null,
    //     producer: {
    //       id: '',
    //       name: 'Chat'
    //     },
    //     producerId: null,
    //     text: message.trim(),
    //     url: null
    //   },
    //   queryClient
    // )

    try {
      await createRun({
        pipelineId: 'chat'
      })
    } catch (error) {
      if (error instanceof Error) {
        console.error(error)
      }
      toast('Failed to send message')
    }
  }

  const handleKeyDown = async (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      await handleSend()
    }
  }

  const handleChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    if (e.target.value.endsWith('@')) {
      setOpen(true)
    } else {
      setOpen(false)
    }
    setMessage(e.target.value)
  }
  const [, setOpen] = useState(false)

  return (
    <div className='relative flex h-full gap-6'>
      {/* Refresh Button */}
      <div className='absolute top-0 left-0 z-10 hidden flex-col gap-2 bg-transparent md:flex'>
        <Button
          className='text-muted-foreground hover:text-primary'
          onClick={() => {
            setLabelId('')
          }}
          size='icon'
          variant={'ghost'}
        >
          <RefreshCcw className='h-5 w-5' />
        </Button>
      </div>

      {/* Chat Body */}
      <div className='flex flex-1 flex-col'>
        {/* Message Area */}
        <ScrollArea className='flex-1 p-4'>
          <div className='flex flex-col gap-4 px-8 xl:px-52'>
            {messages
              ?.slice()
              .reverse()
              .map((msg) => (
                <div
                  className='flex flex-col gap-2'
                  key={msg.id}
                >
                  {/* User Message */}
                  <div className='flex justify-end py-2'>
                    <div className='rounded-lg bg-gray-200 px-4 py-2 text-gray-800 dark:bg-gray-800 dark:text-gray-200'>
                      {msg.attributes.text}
                    </div>
                  </div>
                  {/* Bot Response */}
                  <div className='flex items-start gap-2 py-2'>
                    {/* <Avatar>
                        <ArchesLogo scale={0.124} size="sm" />
                      </Avatar> */}
                    {msg.id === 'pending' ? (
                      <div className='flex items-center justify-center'>
                        <div className='pulse h-5 w-5 rounded-full bg-black'></div>
                      </div>
                    ) : (
                      <div className='rounded-lg py-2'>
                        {msg.attributes
                          .text!.replaceAll(' -', '\n-')
                          .split(/(```[\s\S]*?```)/g)
                          .map((segment, index) => {
                            const replaced = segment
                              .split(/(\*\*[^*]+\*\*|`[^`]+`|\n)/g)
                              .map((part, partIndex) => {
                                if (
                                  part.startsWith('**') &&
                                  part.endsWith('**')
                                ) {
                                  return (
                                    <b key={partIndex}>{part.slice(2, -2)}</b>
                                  )
                                } else if (
                                  part.startsWith('`') &&
                                  part.endsWith('`')
                                ) {
                                  return (
                                    <span
                                      className='markdown-code'
                                      key={partIndex}
                                    >
                                      {part.slice(1, -1)}
                                    </span>
                                  )
                                } else if (part === '\n') {
                                  return <br key={partIndex} />
                                } else {
                                  return part
                                }
                              })

                            return <span key={index}>{replaced}</span>
                          })}
                      </div>
                    )}
                  </div>
                </div>
              ))}
            <div ref={messagesEndRef} />
          </div>
        </ScrollArea>

        {/* Input Form */}
        <form
          onSubmit={async (e) => {
            e.preventDefault()
            await handleSend()
          }}
        >
          <div className='flex items-center gap-2'>
            <Textarea
              className='text-md max-h-40 flex-1 resize-none rounded-lg bg-muted/50 text-gray-800 focus-visible:ring-0 focus-visible:ring-transparent focus-visible:ring-offset-0 dark:text-gray-200'
              onChange={handleChange}
              onInput={(e) => {
                const target = e.target as HTMLTextAreaElement
                target.style.height = 'auto'
                target.style.height = `${target.scrollHeight.toString()}px`
              }}
              onKeyDown={handleKeyDown}
              placeholder='Type your message...'
              rows={1}
              // Auto-resize functionality
              style={{
                height: 'auto',
                overflow: 'hidden'
              }}
              value={message}
            />

            <Button
              className={cn(
                'flex items-center justify-center p-2',
                !message.trim() && 'cursor-not-allowed opacity-50'
              )}
              disabled={!message.trim()}
              type='submit'
            >
              <svg
                className='h-5 w-5 text-white'
                viewBox='0 0 20 20'
                xmlns='http://www.w3.org/2000/svg'
              >
                <title>Send Message</title>
                <path
                  d='M15.44 1.68c.69-.05 1.47.08 2.13.74.66.67.8 1.45.75 2.14-.03.47-.15 1-.25 1.4l-.09.35a43.7 43.7 0 01-3.83 10.67A2.52 2.52 0 019.7 17l-1.65-3.03a.83.83 0 01.14-1l3.1-3.1a.83.83 0 10-1.18-1.17l-3.1 3.1a.83.83 0 01-.99.14L2.98 10.3a2.52 2.52 0 01.04-4.45 43.7 43.7 0 0111.02-3.9c.4-.1.92-.23 1.4-.26Z'
                  fill='currentColor'
                ></path>
              </svg>
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}

import React from 'react'
import { createFileRoute } from '@tanstack/react-router'

import type { LucideIcon } from '@archesai/ui/components/custom/icons'

import {
  AlertCircle,
  ArrowDownLeft,
  ArrowRight,
  ArrowUpRight,
  Calendar,
  CheckCircle2,
  CreditCard,
  PiggyBank,
  Plus,
  QrCode,
  SendHorizontal,
  ShoppingCart,
  Timer,
  TrendingUp,
  Wallet
} from '@archesai/ui/components/custom/icons'
import { cn } from '@archesai/ui/lib/utils'

interface AccountItem {
  balance: string
  description?: string
  id: string
  title: string
  type: 'checking' | 'debt' | 'investment' | 'savings'
}

interface List01Props {
  accounts?: AccountItem[]
  className?: string
  totalBalance?: string
}

const ACCOUNTS: AccountItem[] = [
  {
    balance: '$8,459.45',
    description: 'Personal savings',
    id: '1',
    title: 'Main Savings',
    type: 'savings'
  },
  {
    balance: '$2,850.00',
    description: 'Daily expenses',
    id: '2',
    title: 'Checking Account',
    type: 'checking'
  },
  {
    balance: '$15,230.80',
    description: 'Stock & ETFs',
    id: '3',
    title: 'Investment Portfolio',
    type: 'investment'
  },
  {
    balance: '$1,200.00',
    description: 'Pending charges',
    id: '4',
    title: 'Credit Card',
    type: 'debt'
  },
  {
    balance: '$3,000.00',
    description: 'Emergency fund',
    id: '5',
    title: 'Savings Account',
    type: 'savings'
  }
]

interface List02Props {
  className?: string
  transactions?: Transaction[]
}

interface Transaction {
  amount: string
  category: string
  icon: LucideIcon
  id: string
  status: 'completed' | 'failed' | 'pending'
  timestamp: string
  title: string
  type: 'incoming' | 'outgoing'
}

const TRANSACTIONS: Transaction[] = [
  {
    amount: '$999.00',
    category: 'shopping',
    icon: ShoppingCart,
    id: '1',
    status: 'completed',
    timestamp: 'Today, 2:45 PM',
    title: 'Apple Store Purchase',
    type: 'outgoing'
  },
  {
    amount: '$4,500.00',
    category: 'transport',
    icon: Wallet,
    id: '2',
    status: 'completed',
    timestamp: 'Today, 9:00 AM',
    title: 'Salary Deposit',
    type: 'incoming'
  },
  {
    amount: '$15.99',
    category: 'entertainment',
    icon: CreditCard,
    id: '3',
    status: 'pending',
    timestamp: 'Yesterday',
    title: 'Netflix Subscription',
    type: 'outgoing'
  },
  {
    amount: '$999.00',
    category: 'shopping',
    icon: ShoppingCart,
    id: '4',
    status: 'completed',
    timestamp: 'Today, 2:45 PM',
    title: 'Apple Store Purchase',
    type: 'outgoing'
  },
  {
    amount: '$15.99',
    category: 'entertainment',
    icon: CreditCard,
    id: '5',
    status: 'pending',
    timestamp: 'Yesterday',
    title: 'Supabase Subscription',
    type: 'outgoing'
  },
  {
    amount: '$15.99',
    category: 'entertainment',
    icon: CreditCard,
    id: '6',
    status: 'pending',
    timestamp: 'Yesterday',
    title: 'Vercel Subscription',
    type: 'outgoing'
  }
]

export const Route = createFileRoute('/_app/')({
  component: Dashboard
})

interface List03Props {
  className?: string
  items?: ListItem[]
}

interface ListItem {
  amount?: string
  date: string
  icon: LucideIcon
  iconStyle: string
  id: string
  progress?: number
  status: 'completed' | 'in-progress' | 'pending'
  subtitle: string
  time?: string
  title: string
}

export function Dashboard() {
  return (
    <div className='space-y-4'>
      <div className='grid grid-cols-1 gap-6 lg:grid-cols-2'>
        <div className='flex flex-col rounded-xl border border-gray-200 bg-white p-6 dark:border-[#1F1F23] dark:bg-[#0F0F12]'>
          <h2 className='mb-4 flex items-center gap-2 text-left text-lg font-bold text-gray-900 dark:text-white'>
            <Wallet className='h-3.5 w-3.5 text-zinc-900 dark:text-zinc-50' />
            Accounts
          </h2>
          <div className='flex-1'>
            <List01 className='h-full' />
          </div>
        </div>
        <div className='flex flex-col rounded-xl border border-gray-200 bg-white p-6 dark:border-[#1F1F23] dark:bg-[#0F0F12]'>
          <h2 className='mb-4 flex items-center gap-2 text-left text-lg font-bold text-gray-900 dark:text-white'>
            <CreditCard className='h-3.5 w-3.5 text-zinc-900 dark:text-zinc-50' />
            Recent Transactions
          </h2>
          <div className='flex-1'>
            <List02 className='h-full' />
          </div>
        </div>
      </div>

      <div className='flex flex-col items-start justify-start rounded-xl border border-gray-200 bg-white p-6 dark:border-[#1F1F23] dark:bg-[#0F0F12]'>
        <h2 className='mb-4 flex items-center gap-2 text-left text-lg font-bold text-gray-900 dark:text-white'>
          <Calendar className='h-3.5 w-3.5 text-zinc-900 dark:text-zinc-50' />
          Upcoming Events
        </h2>
        <List03 />
      </div>
    </div>
  )
}

export function List01({
  accounts = ACCOUNTS,
  className,
  totalBalance = '$26,540.25'
}: List01Props) {
  return (
    <div
      className={cn(
        'mx-auto w-full max-w-xl',
        'bg-white dark:bg-zinc-900/70',
        'border border-zinc-100 dark:border-zinc-800',
        'rounded-xl shadow-sm backdrop-blur-xl',
        className
      )}
    >
      {/* Total Balance Section */}
      <div className='border-b border-zinc-100 p-4 dark:border-zinc-800'>
        <p className='text-xs text-zinc-600 dark:text-zinc-400'>
          Total Balance
        </p>
        <h1 className='text-2xl font-semibold text-zinc-900 dark:text-zinc-50'>
          {totalBalance}
        </h1>
      </div>

      {/* Accounts List */}
      <div className='p-3'>
        <div className='mb-2 flex items-center justify-between'>
          <h2 className='text-xs font-medium text-zinc-900 dark:text-zinc-100'>
            Your Accounts
          </h2>
        </div>

        <div className='space-y-1'>
          {accounts.map((account) => (
            <div
              className={cn(
                'group flex items-center justify-between',
                'rounded-lg p-2',
                'hover:bg-zinc-100 dark:hover:bg-zinc-800/50',
                'transition-all duration-200'
              )}
              key={account.id}
            >
              <div className='flex items-center gap-2'>
                <div
                  className={cn('rounded-lg p-1.5', {
                    'bg-blue-100 dark:bg-blue-900/30':
                      account.type === 'checking',
                    'bg-emerald-100 dark:bg-emerald-900/30':
                      account.type === 'savings',
                    'bg-purple-100 dark:bg-purple-900/30':
                      account.type === 'investment'
                  })}
                >
                  {account.type === 'savings' && (
                    <Wallet className='h-3.5 w-3.5 text-emerald-600 dark:text-emerald-400' />
                  )}
                  {account.type === 'checking' && (
                    <QrCode className='h-3.5 w-3.5 text-blue-600 dark:text-blue-400' />
                  )}
                  {account.type === 'investment' && (
                    <ArrowUpRight className='h-3.5 w-3.5 text-purple-600 dark:text-purple-400' />
                  )}
                  {account.type === 'debt' && (
                    <CreditCard className='h-3.5 w-3.5 text-red-600 dark:text-red-400' />
                  )}
                </div>
                <div>
                  <h3 className='text-xs font-medium text-zinc-900 dark:text-zinc-100'>
                    {account.title}
                  </h3>
                  {account.description && (
                    <p className='text-[11px] text-zinc-600 dark:text-zinc-400'>
                      {account.description}
                    </p>
                  )}
                </div>
              </div>

              <div className='text-right'>
                <span className='text-xs font-medium text-zinc-900 dark:text-zinc-100'>
                  {account.balance}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Updated footer with four buttons */}
      <div className='border-t border-zinc-100 p-2 dark:border-zinc-800'>
        <div className='grid grid-cols-4 gap-2'>
          <button
            className={cn(
              'flex items-center justify-center gap-2',
              'rounded-lg px-3 py-2',
              'text-xs font-medium',
              'bg-zinc-900 dark:bg-zinc-50',
              'text-zinc-50 dark:text-zinc-900',
              'hover:bg-zinc-800 dark:hover:bg-zinc-200',
              'shadow-sm hover:shadow',
              'transition-all duration-200'
            )}
            type='button'
          >
            <Plus className='h-3.5 w-3.5' />
            <span>Add</span>
          </button>
          <button
            className={cn(
              'flex items-center justify-center gap-2',
              'rounded-lg px-3 py-2',
              'text-xs font-medium',
              'bg-zinc-900 dark:bg-zinc-50',
              'text-zinc-50 dark:text-zinc-900',
              'hover:bg-zinc-800 dark:hover:bg-zinc-200',
              'shadow-sm hover:shadow',
              'transition-all duration-200'
            )}
            type='button'
          >
            <SendHorizontal className='h-3.5 w-3.5' />
            <span>Send</span>
          </button>
          <button
            className={cn(
              'flex items-center justify-center gap-2',
              'rounded-lg px-3 py-2',
              'text-xs font-medium',
              'bg-zinc-900 dark:bg-zinc-50',
              'text-zinc-50 dark:text-zinc-900',
              'hover:bg-zinc-800 dark:hover:bg-zinc-200',
              'shadow-sm hover:shadow',
              'transition-all duration-200'
            )}
            type='button'
          >
            <ArrowDownLeft className='h-3.5 w-3.5' />
            <span>Top-up</span>
          </button>
          <button
            className={cn(
              'flex items-center justify-center gap-2',
              'rounded-lg px-3 py-2',
              'text-xs font-medium',
              'bg-zinc-900 dark:bg-zinc-50',
              'text-zinc-50 dark:text-zinc-900',
              'hover:bg-zinc-800 dark:hover:bg-zinc-200',
              'shadow-sm hover:shadow',
              'transition-all duration-200'
            )}
            type='button'
          >
            <ArrowRight className='h-3.5 w-3.5' />
            <span>More</span>
          </button>
        </div>
      </div>
    </div>
  )
}

export function List02({
  className,
  transactions = TRANSACTIONS
}: List02Props) {
  return (
    <div
      className={cn(
        'mx-auto w-full max-w-xl',
        'bg-white dark:bg-zinc-900/70',
        'border border-zinc-100 dark:border-zinc-800',
        'rounded-xl shadow-sm backdrop-blur-xl',
        className
      )}
    >
      <div className='p-4'>
        <div className='mb-3 flex items-center justify-between'>
          <h2 className='text-sm font-semibold text-zinc-900 dark:text-zinc-100'>
            Recent Activity
            <span className='ml-1 text-xs font-normal text-zinc-600 dark:text-zinc-400'>
              (23 transactions)
            </span>
          </h2>
          <span className='text-xs text-zinc-600 dark:text-zinc-400'>
            This Month
          </span>
        </div>

        <div className='space-y-1'>
          {transactions.map((transaction) => (
            <div
              className={cn(
                'group flex items-center gap-3',
                'rounded-lg p-2',
                'hover:bg-zinc-100 dark:hover:bg-zinc-800/50',
                'transition-all duration-200'
              )}
              key={transaction.id}
            >
              <div
                className={cn(
                  'rounded-lg p-2',
                  'bg-zinc-100 dark:bg-zinc-800',
                  'border border-zinc-200 dark:border-zinc-700'
                )}
              >
                <transaction.icon className='h-4 w-4 text-zinc-900 dark:text-zinc-100' />
              </div>

              <div className='flex min-w-0 flex-1 items-center justify-between'>
                <div className='space-y-0.5'>
                  <h3 className='text-xs font-medium text-zinc-900 dark:text-zinc-100'>
                    {transaction.title}
                  </h3>
                  <p className='text-[11px] text-zinc-600 dark:text-zinc-400'>
                    {transaction.timestamp}
                  </p>
                </div>

                <div className='flex items-center gap-1.5 pl-3'>
                  <span
                    className={cn(
                      'text-xs font-medium',
                      transaction.type === 'incoming' ?
                        'text-emerald-600 dark:text-emerald-400'
                      : 'text-red-600 dark:text-red-400'
                    )}
                  >
                    {transaction.type === 'incoming' ? '+' : '-'}
                    {transaction.amount}
                  </span>
                  {transaction.type === 'incoming' ?
                    <ArrowDownLeft className='h-3.5 w-3.5 text-emerald-600 dark:text-emerald-400' />
                  : <ArrowUpRight className='h-3.5 w-3.5 text-red-600 dark:text-red-400' />
                  }
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className='border-t border-zinc-100 p-2 dark:border-zinc-800'>
        <button
          className={cn(
            'flex w-full items-center justify-center gap-2',
            'rounded-lg px-3 py-2',
            'text-xs font-medium',
            'bg-gradient-to-r from-zinc-900 to-zinc-800',
            'dark:from-zinc-50 dark:to-zinc-200',
            'text-zinc-50 dark:text-zinc-900',
            'hover:from-zinc-800 hover:to-zinc-700',
            'dark:hover:from-zinc-200 dark:hover:to-zinc-300',
            'shadow-sm hover:shadow',
            'transform transition-all duration-200',
            'hover:-translate-y-0.5',
            'active:translate-y-0',
            'focus:ring-2 focus:outline-none',
            'focus:ring-zinc-500 dark:focus:ring-zinc-400',
            'focus:ring-offset-2 dark:focus:ring-offset-zinc-900'
          )}
          type='button'
        >
          <span>View All Transactions</span>
          <ArrowRight className='h-3.5 w-3.5' />
        </button>
      </div>
    </div>
  )
}

const iconStyles = {
  debt: 'bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100',
  investment: 'bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100',
  savings: 'bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100'
}

const statusConfig = {
  completed: {
    bg: 'bg-emerald-100 dark:bg-emerald-900/30',
    class: 'text-emerald-600 dark:text-emerald-400',
    icon: CheckCircle2
  },
  'in-progress': {
    bg: 'bg-blue-100 dark:bg-blue-900/30',
    class: 'text-blue-600 dark:text-blue-400',
    icon: AlertCircle
  },
  pending: {
    bg: 'bg-amber-100 dark:bg-amber-900/30',
    class: 'text-amber-600 dark:text-amber-400',
    icon: Timer
  }
}

const ITEMS: ListItem[] = [
  {
    amount: '$15,000',
    date: 'Target: Dec 2024',
    icon: PiggyBank,
    iconStyle: 'savings',
    id: '1',
    progress: 65,
    status: 'in-progress',
    subtitle: '3 months of expenses saved',
    title: 'Emergency Fund'
  },
  {
    amount: '$50,000',
    date: 'Target: Jun 2024',
    icon: TrendingUp,
    iconStyle: 'investment',
    id: '2',
    progress: 30,
    status: 'pending',
    subtitle: 'Tech sector investment plan',
    title: 'Stock Portfolio'
  },
  {
    amount: '$25,000',
    date: 'Target: Mar 2025',
    icon: CreditCard,
    iconStyle: 'debt',
    id: '3',
    progress: 45,
    status: 'in-progress',
    subtitle: 'Student loan payoff plan',
    title: 'Debt Repayment'
  }
]

function List03({ className, items = ITEMS }: List03Props) {
  return (
    <div className={cn('scrollbar-none w-full overflow-x-auto', className)}>
      <div className='flex min-w-full gap-3 p-1'>
        {items.map((item) => (
          <div
            className={cn(
              'flex flex-col',
              'w-[280px] shrink-0',
              'bg-white dark:bg-zinc-900/70',
              'rounded-xl',
              'border border-zinc-100 dark:border-zinc-800',
              'hover:border-zinc-200 dark:hover:border-zinc-700',
              'transition-all duration-200',
              'shadow-sm backdrop-blur-xl'
            )}
            key={item.id}
          >
            <div className='space-y-3 p-4'>
              <div className='flex items-start justify-between'>
                <div
                  className={cn(
                    'rounded-lg p-2',
                    iconStyles[item.iconStyle as keyof typeof iconStyles]
                  )}
                >
                  <item.icon className='h-4 w-4' />
                </div>
                <div
                  className={cn(
                    'flex items-center gap-1.5 rounded-full px-2 py-1 text-xs font-medium',
                    statusConfig[item.status].bg,
                    statusConfig[item.status].class
                  )}
                >
                  {React.createElement(statusConfig[item.status].icon, {
                    className: 'w-3.5 h-3.5'
                  })}
                  {item.status.charAt(0).toUpperCase() + item.status.slice(1)}
                </div>
              </div>

              <div>
                <h3 className='mb-1 text-sm font-medium text-zinc-900 dark:text-zinc-100'>
                  {item.title}
                </h3>
                <p className='line-clamp-2 text-xs text-zinc-600 dark:text-zinc-400'>
                  {item.subtitle}
                </p>
              </div>

              {typeof item.progress === 'number' && (
                <div className='space-y-1.5'>
                  <div className='flex items-center justify-between text-xs'>
                    <span className='text-zinc-600 dark:text-zinc-400'>
                      Progress
                    </span>
                    <span className='text-zinc-900 dark:text-zinc-100'>
                      {item.progress}%
                    </span>
                  </div>
                  <div className='h-1.5 overflow-hidden rounded-full bg-zinc-100 dark:bg-zinc-800'>
                    <div
                      className='h-full rounded-full bg-zinc-900 dark:bg-zinc-100'
                      style={{ width: `${item.progress.toString()}%` }}
                    />
                  </div>
                </div>
              )}

              {item.amount && (
                <div className='flex items-center gap-1.5'>
                  <span className='text-sm font-medium text-zinc-900 dark:text-zinc-100'>
                    {item.amount}
                  </span>
                  <span className='text-xs text-zinc-600 dark:text-zinc-400'>
                    target
                  </span>
                </div>
              )}

              <div className='flex items-center text-xs text-zinc-600 dark:text-zinc-400'>
                <Calendar className='mr-1.5 h-3.5 w-3.5' />
                <span>{item.date}</span>
              </div>
            </div>

            <div className='mt-auto border-t border-zinc-100 dark:border-zinc-800'>
              <button
                className={cn(
                  'flex w-full items-center justify-center gap-2',
                  'px-3 py-2.5',
                  'text-xs font-medium',
                  'text-zinc-600 dark:text-zinc-400',
                  'hover:text-zinc-900 dark:hover:text-zinc-100',
                  'hover:bg-zinc-100 dark:hover:bg-zinc-800/50',
                  'transition-colors duration-200'
                )}
              >
                View Details
                <ArrowRight className='h-3.5 w-3.5' />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

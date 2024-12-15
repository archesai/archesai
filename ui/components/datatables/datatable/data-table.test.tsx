import { render, screen } from '@testing-library/react'
import { DataTable } from './data-table'
import { AccessorKeyColumnDef } from '@tanstack/react-table'
import { toSentenceCase } from '@/lib/utils'
import userEvent from '@testing-library/user-event'
import * as hooks from '@/hooks/useSelectItems' // Import hooks module for spying

jest.mock('@/hooks/useSelectItems') // Mock the module

describe('DataTable Component', () => {
  const mockColumns: AccessorKeyColumnDef<any>[] = [
    {
      accessorKey: 'name',
      cell: ({ row }) => <div>{row.original.name}</div>
    },
    {
      accessorKey: 'createdAt',
      cell: ({ row }) => <div>{row.original.createdAt}</div>
    }
  ]
  const mockData = [
    { id: '1', name: 'Item 1', createdAt: new Date().toISOString() },
    { id: '2', name: 'Item 2', createdAt: new Date(0).toISOString() }
  ]
  const mockUseFindAll = jest.fn(() => ({
    data: {
      metadata: { limit: 10, offset: 0, totalResults: 2 },
      results: mockData
    },
    isLoading: false,
    isPlaceholderData: false
  }))
  const mockUseRemove = jest.fn(() => ({
    mutateAsync: jest.fn().mockResolvedValueOnce(undefined)
  }))
  const mockHandleSelect = jest.fn()

  const mockToggleSelectAll = jest.fn()
  const mockToggleSelection = jest.fn()
  const mockSetSelectedItems = jest.fn()
  beforeEach(() => {
    ;(hooks.useSelectItems as jest.Mock).mockImplementation(() => ({
      selectedItems: [],
      selectedAllItems: false,
      selectedSomeItems: false,
      setSelectedItems: mockSetSelectedItems,
      toggleSelectAll: mockToggleSelectAll,
      toggleSelection: mockToggleSelection
    }))
  })

  it('renders table columns and rows', () => {
    render(
      <DataTable
        columns={mockColumns}
        itemType='Item'
        handleSelect={mockHandleSelect}
        useFindAll={mockUseFindAll}
        useRemove={mockUseRemove}
        dataIcon={<div />}
        findAllPathParams={{ limit: 10, offset: 0 }}
      />
    )

    // Check if column headers are rendered
    mockColumns.forEach((column) => {
      expect(
        screen.getByText(toSentenceCase(column.accessorKey.toString()))
      ).toBeInTheDocument()
    })

    // Check if rows are rendered
    mockData.forEach((item) => {
      expect(screen.getByText(item.name)).toBeInTheDocument()
      expect(screen.getByText(item.createdAt)).toBeInTheDocument()
    })
  })

  it('handles sorting correctly', async () => {
    const user = userEvent.setup()
    render(
      <DataTable
        columns={mockColumns}
        itemType='Item'
        handleSelect={mockHandleSelect}
        useFindAll={mockUseFindAll}
        useRemove={mockUseRemove}
        dataIcon={<div />}
        findAllPathParams={{}}
      />
    )

    // Ensure sorting state is updated (mock logic)
    expect(mockUseFindAll).toHaveBeenCalledWith(
      expect.objectContaining({
        queryParams: expect.objectContaining({
          sortBy: 'createdAt',
          sortDirection: 'desc'
        })
      })
    )

    // Simulate clicking on a sortable column header
    const columnHeader = screen.getByText('Name').closest('button')

    // Assert that the column header is a button
    expect(columnHeader).toBeInTheDocument()
    await user.click(columnHeader!)

    const ascendingButton = await screen.findByText('Asc')
    await user.click(ascendingButton)

    // Ensure sorting state is updated (mock logic)
    expect(mockUseFindAll).toHaveBeenCalledWith(
      expect.objectContaining({
        queryParams: expect.objectContaining({
          sortBy: 'name',
          sortDirection: 'asc'
        })
      })
    )
  })

  it('calls useSelectItems with correct parameters', () => {
    render(
      <DataTable
        columns={mockColumns}
        itemType='Item'
        handleSelect={mockHandleSelect}
        useFindAll={mockUseFindAll}
        useRemove={mockUseRemove}
        dataIcon={<div />}
        findAllPathParams={{ limit: 10, offset: 0 }}
      />
    )

    // Verify that useSelectItems was called with the correct items
    expect(hooks.useSelectItems).toHaveBeenCalledWith({
      items: mockData
    })
  })

  it('allows selecting rows', async () => {
    const user = userEvent.setup()
    render(
      <DataTable
        columns={mockColumns}
        itemType='Item'
        handleSelect={mockHandleSelect}
        useFindAll={mockUseFindAll}
        useRemove={mockUseRemove}
        dataIcon={<div />}
        findAllPathParams={{ limit: 10, offset: 0 }}
      />
    )

    // Simulate select all
    const checkbox = screen.getAllByRole('checkbox')[0]
    expect(checkbox).toBeInTheDocument()
    await user.click(checkbox!)
    expect(mockToggleSelectAll).toHaveBeenCalled()

    // Similate deselect all
    await user.click(checkbox!)
    expect(mockToggleSelectAll).toHaveBeenCalled()

    // Simulate selecting a single row
    const rowCheckbox = screen.getAllByRole('checkbox')[1]
    expect(rowCheckbox).toBeInTheDocument()
    await user.click(rowCheckbox!)

    // Ensure handleSelect is called with correct parameters
    expect(mockToggleSelection).toHaveBeenCalledWith(mockData[0]?.id)
  })

  it('handles delete action', async () => {
    const user = userEvent.setup()

    render(
      <DataTable
        columns={mockColumns}
        itemType='Item'
        handleSelect={mockHandleSelect}
        useFindAll={mockUseFindAll}
        useRemove={mockUseRemove}
        getDeleteVariablesFromItem={(item) => ({ id: item.id })}
        dataIcon={<div />}
        findAllPathParams={{ limit: 10, offset: 0 }}
      />
    )

    // Open dropdown menu for first row
    const dropdownButton = screen.getAllByRole('button', {
      name: 'Expand row options'
    })[0]
    await user.click(dropdownButton!)

    // Click delete button
    const deleteButton = await screen.findByText('Delete')
    expect(deleteButton).toBeInTheDocument()
    await user.click(deleteButton)

    // Confirm popup
    const confirmDialog = await screen.findByRole('dialog')
    expect(confirmDialog).toBeInTheDocument()

    // Click confirm button
    const confirmButton = screen.getByRole('button', { name: 'Delete' })
    expect(confirmButton).toBeInTheDocument()
    await user.click(confirmButton)

    // Ensure delete function is called
    expect(
      mockUseRemove.mock.results[0]?.value.mutateAsync
    ).toHaveBeenCalledWith({
      id: '1'
    })
  })

  it('renders grid view', () => {
    render(
      <DataTable
        columns={mockColumns}
        itemType='Item'
        handleSelect={mockHandleSelect}
        useFindAll={mockUseFindAll}
        useRemove={mockUseRemove}
        defaultView='grid'
        dataIcon={<div />}
        findAllPathParams={{ limit: 10, offset: 0 }}
      />
    )

    // Check that GridView is rendered
    expect(screen.getByText('Item 1')).toBeInTheDocument()
    expect(screen.getByText('Item 2')).toBeInTheDocument()
  })
})

// import type { AccessorKeyColumnDef } from '@tanstack/react-table'

// import { render, screen } from '@testing-library/react' // FIXME (these are uninstalled)
// import userEvent from '@archesai/jest/user-event' // FIXME (these are uninstalled)

// import type { BaseEntity } from '@archesai/schemas'
// import { toSentenceCase } from '#lib/utils'

// import { DataTable } from './data-table'

// jest.mock('#hooks/use-select-items') // Mock the module

// describe('DataTable Component', () => {
//   const mockColumns: AccessorKeyColumnDef<BaseEntity>[] = [
//     {
//       accessorKey: 'name',
//       cell: ({ row }) => <div>{row.original.name}</div>
//     },
//     {
//       accessorKey: 'createdAt',
//       cell: ({ row }) => <div>{row.original.createdAt}</div>
//     }
//   ]
//   const mockData = [
//     { createdAt: new Date().toISOString(), id: '1', name: 'Item 1' },
//     { createdAt: new Date(0).toISOString(), id: '2', name: 'Item 2' }
//   ]
//   const mockUseFindAll = jest.fn(() => ({
//     data: {
//       data: mockData,
//       meta: { total_records: 2 }
//     },
//     isFetched: true,
//     isLoading: false,
//     isPlaceholderData: false
//   }))
//   const mockUseRemove = jest.fn(() => ({
//     mutateAsync: jest.fn().mockResolvedValueOnce(undefined)
//   }))
//   const mockHandleSelect = jest.fn()

//   const mockToggleSelectAll = jest.fn()
//   const mockToggleSelection = jest.fn()
//   const mockSetSelectedItems = jest.fn()
//   beforeEach(() => {
//     ;(hooks.useSelectItems as jest.Mock).mockImplementation(() => ({
//       selectedAllItems: false,
//       selectedItems: [],
//       selectedSomeItems: false,
//       setSelectedItems: mockSetSelectedItems,
//       toggleSelectAll: mockToggleSelectAll,
//       toggleSelection: mockToggleSelection
//     }))
//   })

//   it('renders table columns and rows', () => {
//     render(
//       <DataTable
//         columns={mockColumns}
//         dataIcon={<div />}
//         findManyPathParams={{ limit: 10, offset: 0 }}
//         handleSelect={mockHandleSelect}
//         itemType='Item'
//         useFindAll={mockUseFindAll}
//         useRemove={mockUseRemove}
//       />
//     )

//     // Check if column headers are rendered
//     mockColumns.forEach((column) => {
//       expect(
//         screen.getByText(toSentenceCase(column.accessorKey.toString()))
//       ).toBeInTheDocument()
//     })

//     // Check if rows are rendered
//     mockData.forEach((item) => {
//       expect(screen.getByText(item.name)).toBeInTheDocument()
//       expect(screen.getByText(item.createdAt)).toBeInTheDocument()
//     })
//   })

//   it('handles sorting correctly', async () => {
//     const user = userEvent.setup()
//     render(
//       <DataTable
//         columns={mockColumns}
//         dataIcon={<div />}
//         findManyPathParams={{}}
//         handleSelect={mockHandleSelect}
//         itemType='Item'
//         useFindAll={mockUseFindAll}
//         useRemove={mockUseRemove}
//       />
//     )

//     // Ensure sorting state is updated (mock logic)
//     // expect(mockUseFindAll).toHaveBeenCalledWith(
//     //   expect.objectContaining({
//     //     queryParams: expect.objectContaining({
//     //       sortBy: 'createdAt',
//     //       sortDirection: 'desc'
//     //     })
//     //   })
//     // )

//     // Simulate clicking on a sortable column header
//     const columnHeader = screen.getByText('Name').closest('button')

//     // Assert that the column header is a button
//     expect(columnHeader).toBeInTheDocument()
//     if (!columnHeader) throw new Error('Column header not found')
//     await user.click(columnHeader)

//     const ascendingButton = await screen.findByText('Asc')
//     await user.click(ascendingButton)

//     // Ensure sorting state is updated (mock logic)
//     // expect(mockUseFindAll).toHaveBeenCalledWith(
//     //   expect.objectContaining({
//     //     queryParams: expect.objectContaining({
//     //       sortBy: 'name',
//     //       sortDirection: 'asc'
//     //     })
//     //   })
//     // )
//   })

//   it('calls useSelectItems with correct parameters', () => {
//     render(
//       <DataTable
//         columns={mockColumns}
//         dataIcon={<div />}
//         findManyPathParams={{ limit: 10, offset: 0 }}
//         handleSelect={mockHandleSelect}
//         itemType='Item'
//         useFindAll={mockUseFindAll}
//         useRemove={mockUseRemove}
//       />
//     )

//     // Verify that useSelectItems was called with the correct items
//     expect(hooks.useSelectItems).toHaveBeenCalledWith({
//       items: mockData
//     })
//   })

//   it('allows selecting rows', () => {
//     // const user = userEvent.setup()
//     render(
//       <DataTable
//         columns={mockColumns}
//         dataIcon={<div />}
//         findManyPathParams={{ limit: 10, offset: 0 }}
//         handleSelect={mockHandleSelect}
//         itemType='Item'
//         useFindAll={mockUseFindAll}
//         useRemove={mockUseRemove}
//       />
//     )

//     // Simulate select all
//     // const checkbox = screen.getAllByRoleTypeEnum('checkbox')[0]
//     // expect(checkbox).toBeInTheDocument()
//     // if (!checkbox) throw new Error('Checkbox not found')
//     // await user.click(checkbox)
//     // expect(mockToggleSelectAll).toHaveBeenCalled()

//     // // Similate deselect all
//     // await user.click(checkbox)
//     // expect(mockToggleSelectAll).toHaveBeenCalled()

//     // // Simulate selecting a single row
//     // const rowCheckbox = screen.getAllByRoleTypeEnum('checkbox')[1]
//     // expect(rowCheckbox).toBeInTheDocument()
//     // if (!rowCheckbox) throw new Error('Checkbox not found')
//     // await user.click(rowCheckbox)

//     // Ensure handleSelect is called with correct parameters
//     expect(mockToggleSelection).toHaveBeenCalledWith(mockData[0]?.id)
//   })

//   // it('handles delete action', async () => {
//   //   const user = userEvent.setup()

//   //   render(
//   //     <DataTable
//   //       columns={mockColumns}
//   //       dataIcon={<div />}
//   //       findManyPathParams={{ limit: 10, offset: 0 }}
//   //       getDeleteVariablesFromItem={(item) => ({ id: item.id })}
//   //       handleSelect={mockHandleSelect}
//   //       itemType='Item'
//   //       useFindAll={mockUseFindAll}
//   //       useRemove={mockUseRemove}
//   //     />
//   //   )

//   //   // Open dropdown menu for first row
//   //   const dropdownButton = screen.getAllByRoleTypeEnum('button', {
//   //     name: 'Expand row options'
//   //   })[0]
//   //   expect(dropdownButton).toBeInTheDocument()
//   //   if (!dropdownButton) throw new Error('Dropdown button not found')
//   //   await user.click(dropdownButton)

//   //   // Click delete button
//   //   const deleteButton = await screen.findByText('Delete')
//   //   expect(deleteButton).toBeInTheDocument()
//   //   await user.click(deleteButton)

//   //   // Confirm popup
//   //   const confirmDialog = await screen.findByRoleTypeEnum('dialog')
//   //   expect(confirmDialog).toBeInTheDocument()

//   //   // Click confirm button
//   //   const confirmButton = screen.getByRoleTypeEnum('button', { name: 'Delete' })
//   //   expect(confirmButton).toBeInTheDocument()
//   //   await user.click(confirmButton)

//   //   // Ensure delete function is called
//   //   expect(
//   //     mockUseRemove.mock.results[0]?.value.mutateAsync
//   //   ).toHaveBeenCalledWith({
//   //     id: '1'
//   //   })
//   // })

//   it('renders grid view', () => {
//     render(
//       <DataTable
//         columns={mockColumns}
//         dataIcon={<div />}
//         defaultView='grid'
//         findManyPathParams={{ limit: 10, offset: 0 }}
//         handleSelect={mockHandleSelect}
//         itemType='Item'
//         useFindAll={mockUseFindAll}
//         useRemove={mockUseRemove}
//       />
//     )

//     // Check that GridView is rendered
//     expect(screen.getByText('Item 1')).toBeInTheDocument()
//     expect(screen.getByText('Item 2')).toBeInTheDocument()
//   })
// })

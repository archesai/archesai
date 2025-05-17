import { useAtom } from 'jotai' // Adjust path as necessary

import type { BaseEntity } from '@archesai/domain'

import { selectedItemsAtom } from '#atoms/selected'

export const useSelectItems = ({ items }: { items: BaseEntity[] }) => {
  const [selectedItems, setSelectedItems] = useAtom(selectedItemsAtom)

  // Says if some items are selected
  const selectedSomeItems =
    selectedItems.length > 0 && selectedItems.length < items.length

  // Says if all items are selected
  const selectedAllItems =
    selectedItems.length === items.length && items.length > 0

  // Select or deselect all items
  const toggleSelectAll = (): void => {
    setSelectedItems(!selectedAllItems ? items.map((item) => item.id) : [])
  }

  // Toggle selection of one item
  const toggleSelection = (itemId: string): void => {
    if (!selectedItems.includes(itemId)) {
      setSelectedItems((prevSelected) => [...prevSelected, itemId])
    } else {
      setSelectedItems((prevSelected) =>
        prevSelected.filter((id) => id !== itemId)
      )
    }
  }

  return {
    selectedAllItems,
    selectedItems,
    selectedSomeItems,

    setSelectedItems,
    toggleSelectAll,
    toggleSelection
  }
}

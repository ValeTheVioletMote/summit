package player

import "github.com/paalgyula/summit/pkg/wow"

const InventorySlotBagEnd = 23

type InventoryItem struct {
	DisplayInfoID uint32            `yaml:"display_info_id"`
	InventoryType wow.InventoryType `yaml:"inventory_type"`
	EnchantSlot   uint32            `yaml:"enchant_slot"`
}

type Inventory struct {
	InventorySlots []*InventoryItem `yaml:"inventory_slots"`
}

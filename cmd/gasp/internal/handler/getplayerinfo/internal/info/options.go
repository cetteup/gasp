package info

import (
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/internal/dto"
	"github.com/cetteup/gasp/internal/constraints"
)

type ResolveOptions struct {
	DefaultKeys  []string
	ArmyIDs      []uint8
	FieldIDs     []uint16
	KitIDs       []uint8
	VehicleIDs   []uint8
	WeaponIDs    []uint8
	EquipmentIDs []uint8
}

func NewResolveOptions() *ResolveOptions {
	return &ResolveOptions{
		DefaultKeys:  nil,
		ArmyIDs:      dto.ArmyIDs,
		FieldIDs:     dto.FieldIDs,
		KitIDs:       dto.KitIDs,
		VehicleIDs:   dto.VehicleIDs,
		WeaponIDs:    dto.WeaponIDs,
		EquipmentIDs: dto.EquipmentIDs,
	}
}

func (o *ResolveOptions) SetDefaultKeys(keys ...string) *ResolveOptions {
	o.DefaultKeys = make([]string, 0, len(keys))
	for _, key := range keys {
		o.DefaultKeys = append(o.DefaultKeys, key)
	}
	return o
}

func (o *ResolveOptions) MaybeSetArmyIDs(ids ...*uint8) *ResolveOptions {
	o.ArmyIDs = pick(o.ArmyIDs, ids)
	return o
}

func (o *ResolveOptions) MaybeSetFieldIDs(ids ...*uint16) *ResolveOptions {
	o.FieldIDs = pick(o.FieldIDs, ids)
	return o
}

func (o *ResolveOptions) MaybeSetKitIDs(ids ...*uint8) *ResolveOptions {
	o.KitIDs = pick(o.KitIDs, ids)
	return o
}

func (o *ResolveOptions) MaybeSetVehicleIDs(ids ...*uint8) *ResolveOptions {
	o.VehicleIDs = pick(o.VehicleIDs, ids)
	return o
}

func (o *ResolveOptions) MaybeSetWeaponIDs(ids ...*uint8) *ResolveOptions {
	o.WeaponIDs = pick(o.WeaponIDs, ids)
	return o
}

func (o *ResolveOptions) MaybeSetEquipmentIDs(ids ...*uint8) *ResolveOptions {
	o.EquipmentIDs = pick(o.EquipmentIDs, ids)
	return o
}

func pick[T constraints.UnsignedInteger](definitely []T, maybe []*T) []T {
	l := len(definitely)
	for _, m := range maybe {
		if m != nil {
			definitely = append(definitely, *m)
		}
	}
	// Remove previously existing entries if any were added
	if len(definitely) > l {
		definitely = definitely[l:]
	}

	return definitely
}

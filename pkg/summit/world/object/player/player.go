package player

import (
	"math"

	"github.com/paalgyula/summit/pkg/summit/world/object"
	"github.com/paalgyula/summit/pkg/wow"
)

type WorldLocation struct {
	X           float32 `yaml:"position_x"`
	Y           float32 `yaml:"position_y"`
	Z           float32 `yaml:"position_z"`
	Map         uint32  `yaml:"position_map"`
	Zone        uint32  `yaml:"position_zone"`
	Orientation float32 `yaml:"position_orientation"`
}

// Location returns the X, Y, Z coordinates and map ID of a WorldLocation.
//
// No parameters are required.
// Returns a float32 for X, Y, and Z coordinates, and a uint32 for the map ID.
func (l WorldLocation) Location() (float32, float32, float32, uint32) {
	return l.X, l.Y, l.Z, l.Map
}

// Distance calculates the distance between two locations
func (loc *WorldLocation) Distance(point *WorldLocation) float64 {
	dx := loc.X - point.X
	dy := loc.Y - point.Y
	dz := loc.Z - point.Z

	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

func NewPlayer() *Player {
	p := &Player{}

	return p
}

func CreatePlayer() {
	p := NewPlayer()
	// uint8 powertype = cEntry->powerType;

	var unitfield uint32
	powertype := wow.PowerTypeRage

	_, _, _ = p, unitfield, powertype

	// switch powertype {
	// case wow.PowerTypeEnergy,wow.PowerTypeMana:
	// 	unitfield = 0x00000000
	// case wow.PowerTypeRage:
	// 	unitfield = 0x00110000
	// default:
	// 	log.Warn().Msgf("Invalid default powertype %s for player (class %T)", powertype)
	// 	return
	// }

	// p.Object.SetFloatValue(object.UnitFieldBoundingradius), DEFAULT_WORLD_OBJECT_SIZE);
	// p.Object.SetFloatValue(UNIT_FIELD_COMBATREACH, DEFAULT_COMBAT_REACH);

	// switch (gender)
	// {
	//     case GENDER_FEMALE:
	//         SetDisplayId(info->displayId_f);
	//         SetNativeDisplayId(info->displayId_f);
	//         break;
	//     case GENDER_MALE:
	//         SetDisplayId(info->displayId_m);
	//         SetNativeDisplayId(info->displayId_m);
	//         break;
	//     default:
	//         sLog.outError("Invalid gender %u for player", gender);
	//         return false;
	//         break;
	// }

	// setFactionForRace(race);

	// RaceClassGender uint32 = (race) | (class_ << 8) | (gender << 16);

	// p.SetUInt32Value(UNIT_FIELD_BYTES_0, (RaceClassGender | (powertype << 24)));
	// SetUInt32Value(UNIT_FIELD_BYTES_1, unitfield);
	// SetByteValue(UNIT_FIELD_BYTES_2, 1, UNIT_BYTE2_FLAG_SANCTUARY | UNIT_BYTE2_FLAG_UNK5);
	// SetUInt32Value(UNIT_FIELD_FLAGS, UNIT_FLAG_PVP_ATTACKABLE);
	// SetFloatValue(UNIT_MOD_CAST_SPEED, 1.0f);               // fix cast time showed in spell tooltip on client

	// // -1 is default value
	// SetInt32Value(PLAYER_FIELD_WATCHED_FACTION_INDEX, uint32(-1));

	// SetUInt32Value(PLAYER_BYTES, (skin | (face << 8) | (hairStyle << 16) | (hairColor << 24)));
	// SetUInt32Value(PLAYER_BYTES_2, (facialHair | (0x00 << 8) | (0x00 << 16) | (0x02 << 24)));
	// SetByteValue(PLAYER_BYTES_3, 0, gender);

	// SetUInt32Value(PLAYER_GUILDID, 0);
	// SetUInt32Value(PLAYER_GUILDRANK, 0);
	// SetUInt32Value(PLAYER_GUILD_TIMESTAMP, 0);

	// for (int i = 0; i < KNOWN_TITLES_SIZE; ++i)
	//     SetUInt64Value(PLAYER__FIELD_KNOWN_TITLES + i, 0);  // 0=disabled
	// SetUInt32Value(PLAYER_CHOSEN_TITLE, 0);

	// SetUInt32Value(PLAYER_FIELD_KILLS, 0);
	// SetUInt32Value(PLAYER_FIELD_LIFETIME_HONORABLE_KILLS, 0);
	// SetUInt32Value(PLAYER_FIELD_TODAY_CONTRIBUTION, 0);
	// SetUInt32Value(PLAYER_FIELD_YESTERDAY_CONTRIBUTION, 0);

	// // set starting level
	// uint32 start_level = sWorld.getConfig(CONFIG_START_PLAYER_LEVEL);

	// if (GetSession()->GetSecurity() >= SEC_MODERATOR)
	// {
	//     uint32 gm_level = sWorld.getConfig(CONFIG_START_GM_LEVEL);
	//     if (gm_level > start_level)
	//         start_level = gm_level;
	// }

	// SetUInt32Value(UNIT_FIELD_LEVEL, start_level);
	// SetUInt32Value (PLAYER_FIELD_COINAGE, sWorld.getConfig(CONFIG_START_PLAYER_MONEY));
	// SetUInt32Value (PLAYER_FIELD_HONOR_CURRENCY, sWorld.getConfig(CONFIG_START_HONOR_POINTS));
	// SetUInt32Value (PLAYER_FIELD_ARENA_CURRENCY, sWorld.getConfig(CONFIG_START_ARENA_POINTS));

	// // start with every map explored
	// if (sWorld.getConfig(CONFIG_START_ALL_EXPLORED))
	// {
	//     for (uint8 i = 0; i < 64; i++)
	//         SetFlag(PLAYER_EXPLORED_ZONES_1 + i, 0xFFFFFFFF);
	// }
}

type Player struct {
	*object.Object `yaml:"-"`
	*object.Unit   `yaml:"-"`

	ID     uint32           `yaml:"id"`
	Name   string           `yaml:"name"`
	Race   wow.PlayerRace   `yaml:"race"`
	Class  wow.PlayerClass  `yaml:"class"`
	Gender wow.PlayerGender `yaml:"gender"`

	Skin       uint8 `yaml:"skin"`
	Face       uint8 `yaml:"face"`
	HairStyle  uint8 `yaml:"hairstyle"`
	HairColor  uint8 `yaml:"haircolor"`
	FacialHair uint8 `yaml:"facialhair"`
	OutfitID   uint8 `yaml:"outfit_id"`

	Location     WorldLocation `yaml:"location"`
	BindLocation WorldLocation `yaml:"bind_location"`

	Level uint8 `yaml:"level"`

	Inventory *Inventory `yaml:"inventory"`
	GuildID   uint32     `yaml:"guild_id"`

	// CharFlags for example dead, and display ghost
	CharFlags uint32 `yaml:"char_flags"`

	// Recustomization flags (change name, look, etc)
	// Needs some research
	Recustomization uint32 `yaml:"recustomization"`

	// Boolean, but uint8 :D
	FirstLogin uint8 `yaml:"had_first_login"`

	Pet Pet `yaml:"pet"`
}

func (p *Player) Guid() wow.GUID {
	return p.Object.Guid()
}

// Initializes an empty inventory
func (p *Player) InitInventory() {
	if p.Inventory != nil {
		return
	}

	p.Inventory = &Inventory{
		InventorySlots: []*InventoryItem{},
	}

	for i := 0; i < InventorySlotBagEnd; i++ {
		p.Inventory.InventorySlots = append(p.Inventory.InventorySlots, &InventoryItem{})
	}
}

func (p *Player) SetFloatValue() {

}

func (p *Player) GUID() wow.GUID {
	return wow.NewPlayerGUID(p.ID)
}

func (p *Player) Init() {
	p.InitInventory()
}

func (p *Player) Transport() *object.Transport {
	return nil
}

func (p *Player) BuildCreateUpdateForPlayer(target *Player) {
	updatetype := wow.UpdateTypeCreateObject
	flags := p.Object.UpdateFlags()

	_ = updatetype

	/** lower flag1 **/
	if target != nil { // building packet for oneself
		flags |= wow.UpdateFlagSelf
	}

	if flags&wow.UpdateFlagHasPosition != 0 {
		// UPDATETYPE_CREATE_OBJECT2 dynamic objects, corpses...
		// if isType(TYPEMASK_DYNAMICOBJECT) || isType(TYPEMASK_CORPSE) || isType(TYPEMASK_PLAYER) {
		// 	updatetype = wow.UpdateTypeCreateObject2
		// }

		// UPDATETYPE_CREATE_OBJECT2 for pets...
		// if target.GetPetGUID() == p.GetGUID() {
		//     updatetype = UPDATETYPE_CREATE_OBJECT2
		// }

		// UPDATETYPE_CREATE_OBJECT2 for some gameobject types...
		// if (isType(TYPEMASK_GAMEOBJECT))
		// {
		//     switch (((GameObject*)this)->GetGoType())
		//     {
		//     case GAMEOBJECT_TYPE_TRAP:
		//     case GAMEOBJECT_TYPE_DUEL_ARBITER:
		//     case GAMEOBJECT_TYPE_FLAGSTAND:
		//     case GAMEOBJECT_TYPE_FLAGDROP:
		//         updatetype = UPDATETYPE_CREATE_OBJECT2;
		//         break;
		//     case GAMEOBJECT_TYPE_TRANSPORT:
		//         flags |= UPDATEFLAG_TRANSPORT;
		//         break;
		//     default:
		//         break;
		//     }
		// }
	}
}

func (p *Player) WriteToLogin(w *wow.Packet) {
	w.Write(p.GUID())
	w.WriteString(p.Name)
	w.Write(p.Race)
	w.Write(p.Class)
	w.Write(p.Gender)

	w.Write(p.Skin)
	w.Write(p.Face)
	w.Write(p.HairStyle)
	w.Write(p.HairColor)
	w.Write(p.FacialHair)

	w.Write(p.Level)

	w.Write(p.Location.Zone)
	w.Write(p.Location.Map)

	w.Write(p.Location.X)
	w.Write(p.Location.Y)
	w.Write(p.Location.Z)

	w.Write(p.GuildID)

	// Character flags
	w.Write(p.CharFlags)
	w.Write(p.Recustomization)

	// First login
	// *data << uint8(atLoginFlags & AT_LOGIN_FIRST ? 1 : 0);
	w.Write(p.FirstLogin)

	// Player Pet section
	w.Write(p.Pet.DisplayID)
	w.Write(p.Pet.PetLevel)
	w.Write(p.Pet.PetFamilly)

	for _, slot := range p.Inventory.InventorySlots {
		w.Write(slot.DisplayInfoID)
		w.Write(slot.InventoryType)
		w.Write(slot.EnchantSlot)
	}

	// Yipeee
}

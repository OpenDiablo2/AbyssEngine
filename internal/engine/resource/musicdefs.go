package resource

import "github.com/OpenDiablo2/AbyssEngine/internal/engine/common/enum"

// MusicDef stores the music definitions of a region
type MusicDef struct {
	Region    enum.RegionIdType
	InTown    bool
	MusicFile string
}

func getMusicDefs() []MusicDef {
	return []MusicDef{
		{enum.RegionAct1Town, false, BGMAct1Town1},
		{enum.RegionAct1Wilderness, false, BGMAct1Wild},
		{enum.RegionAct1Cave, false, BGMAct1Caves},
		{enum.RegionAct1Crypt, false, BGMAct1Crypt},
		{enum.RegionAct1Monestary, false, BGMAct1Monastery},
		{enum.RegionAct1Courtyard, false, BGMAct1Monastery},
		{enum.RegionAct1Barracks, false, BGMAct1Monastery},
		{enum.RegionAct1Jail, false, BGMAct1Monastery},
		{enum.RegionAct1Cathedral, false, BGMAct1Monastery},
		{enum.RegionAct1Catacombs, false, BGMAct1Monastery},
		{enum.RegionAct1Tristram, false, BGMAct1Tristram},
		{enum.RegionAct2Town, false, BGMAct2Town2},
		{enum.RegionAct2Sewer, false, BGMAct2Sewer},
		{enum.RegionAct2Harem, false, BGMAct2Harem},
		{enum.RegionAct2Basement, false, BGMAct2Harem},
		{enum.RegionAct2Desert, false, BGMAct2Desert},
		{enum.RegionAct2Tomb, false, BGMAct2Tombs},
		{enum.RegionAct2Lair, false, BGMAct2Lair},
		{enum.RegionAct2Arcane, false, BGMAct2Sanctuary},
		{enum.RegionAct3Town, false, BGMAct3Town3},
		{enum.RegionAct3Jungle, false, BGMAct3Jungle},
		{enum.RegionAct3Kurast, false, BGMAct3Kurast},
		{enum.RegionAct3Spider, false, BGMAct3Spider},
		{enum.RegionAct3Dungeon, false, BGMAct3KurastSewer},
		{enum.RegionAct3Sewer, false, BGMAct3KurastSewer},
		{enum.RegionAct4Town, false, BGMAct4Town4},
		{enum.RegionAct4Mesa, false, BGMAct4Mesa},
		{enum.RegionAct4Lava, false, BGMAct4Mesa},
		{enum.RegonAct5Town, false, BGMAct5XTown},
		{enum.RegionAct5Siege, false, BGMAct5Siege},
		{enum.RegionAct5Barricade, false, BGMAct5Siege}, // ?
		{enum.RegionAct5Temple, false, BGMAct5XTemple},
		{enum.RegionAct5IceCaves, false, BGMAct5IceCaves},
		{enum.RegionAct5Baal, false, BGMAct5Baal},
		{enum.RegionAct5Lava, false, BGMAct5Nihlathak}, // ?
	}
}

// GetMusicDef returns the MusicDef of the given region
func GetMusicDef(regionType enum.RegionIdType) *MusicDef {
	musicDefs := getMusicDefs()
	for idx := range musicDefs {
		if musicDefs[idx].Region != regionType {
			continue
		}

		return &musicDefs[idx]
	}

	return &musicDefs[0]
}

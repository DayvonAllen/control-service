package domain

import (
	"fmt"
	"strings"
)

type Tag struct {
	Value          string `bson:"value" json:"value"`
	CreepyPasta    bool   `bson:"-" json:"-"`
	TrueScaryStory bool   `bson:"-" json:"-"`
	CampFire       bool   `bson:"-" json:"-"`
	Paranormal     bool   `bson:"-" json:"-"`
	GhostStory     bool   `bson:"-" json:"-"`
	Other          bool   `bson:"-" json:"-"`
}

func (t Tag) ValidateTag(tagValidator *Tag) error {
	switch tag := strings.ToLower(t.Value); tag {
	case "creepypasta":
		if !tagValidator.CreepyPasta {
			tagValidator.CreepyPasta = true
			return nil
		}
		return fmt.Errorf("no duplicate tags")
	case "truescarystory":
		if !tagValidator.TrueScaryStory {
			tagValidator.TrueScaryStory = true
			return nil
		}
		return fmt.Errorf("no duplicate tags")
	case "campfire":
		if !tagValidator.CampFire {
			tagValidator.CampFire = true
			return nil
		}
		return fmt.Errorf("no duplicate tags")
	case "ghoststory":
		if !tagValidator.GhostStory {
			tagValidator.GhostStory = true
			return nil
		}
		return fmt.Errorf("no duplicate tags")
	case "paranormal":
		if !tagValidator.Paranormal {
			tagValidator.Paranormal = true
			return nil
		}
		return fmt.Errorf("no duplicate tags")
	case "other":
		if !tagValidator.Other {
			tagValidator.Other = true
			return nil
		}
		return fmt.Errorf("no duplicate tags")
	default:
		return fmt.Errorf("invalid tag")
	}
}

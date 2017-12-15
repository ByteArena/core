package utils

type Tag uint64 // limited to 64 components

func (tag Tag) IsIncludedIn(biggertag Tag) bool {
	return tag&biggertag == tag
}

func (tag Tag) Includes(smallertag Tag) bool {
	return tag&smallertag == smallertag
}

func BuildTag(tags ...Tag) Tag {

	tag := Tag(0)

	for _, othertag := range tags {
		tag |= othertag
	}

	return tag
}

func MakeTag(rank uint64) Tag {
	if rank > 63 {
		panic("Cannot make tag with rank over 63")
	}

	var value uint64 = (1 << rank)

	return Tag(value)
}

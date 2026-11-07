package clothes

import "github.com/boundedinfinity/canonical-go/idiomatic/ider"

type SockHeight struct {
	ID          ider.Id
	Name        string
	Aliases     []string
	Description string
}

var SockHeightData = []SockHeight{
	{
		ID:          ider.MustParse("2A7A4080-CF3B-49B3-9AF3-274A98B04C29"),
		Name:        "Toe Cover",
		Aliases:     []string{"Mute", "Half Sock"},
		Description: "",
	},
	{
		ID:   ider.MustParse("D233F757-B4F8-46EC-81F7-14EFA0238829"),
		Name: "No Show Hidden",
	},
	{
		ID:   ider.MustParse("6020F3AA-0C51-4FA4-97C0-D700399A8846"),
		Name: "No Show",
	},
	{
		ID:      ider.MustParse("DE40D506-0DFD-4A39-9D0D-EF37E8D67519"),
		Name:    "Ankle",
		Aliases: []string{"1/4"},
	},
	{
		ID:   ider.MustParse("5C7F33FC-9A66-4EAD-BFCD-94D3811DEA7A"),
		Name: "Micro Crew",
	},
	{
		ID:   ider.MustParse("1D609ABA-FD35-40F3-83B9-9FA58E09C8F0"),
		Name: "Crew",
	},
	{
		ID:   ider.MustParse("7B64F702-20F0-4459-8BF3-E5E8380DA5DF"),
		Name: "Boot Sock",
	},
	{
		ID:   ider.MustParse("B0C2CC04-344E-4D35-8A77-275838A2A6FA"),
		Name: "Mid-Calf",
	},
	{
		ID:   ider.MustParse("18C9EC7F-2F7F-43A4-B818-FEAEF4D79D47"),
		Name: "Over-The-Calf",
	},
	{
		ID:   ider.MustParse("A9A09161-8989-4349-A315-D9F1B1214494"),
		Name: "Knee-high",
	},
	{
		ID:   ider.MustParse("93611678-0134-4069-B2B4-B753672F21AD"),
		Name: "Over-this-knee",
	},
	{
		ID:   ider.MustParse("31E6B0A8-3219-42E9-8F23-3C1D573107BB"),
		Name: "Thigh-high",
	},
}

type SockWeight struct {
	ID          ider.Id
	Name        string
	Aliases     []string
	Description string
}

var SockWeightData = []SockHeight{
	{
		ID:   ider.MustParse("32F2FB7A-DE71-4A4B-BE3D-265E7A5CBCF0"),
		Name: "Ultra-then",
	},
	{
		ID:   ider.MustParse("1713E9E6-FF46-4E95-9321-D01182E2837D"),
		Name: "Lightweight",
	},
	{
		ID:   ider.MustParse("E2BDFAF1-F8BB-4514-9FB1-D0D669B5AD29"),
		Name: "Midweight",
	},
	{
		ID:      ider.MustParse("C290A6CC-EB0B-4C99-852C-4878EC8EED34"),
		Name:    "Heavyweight",
		Aliases: []string{"Thermal"},
	},
}

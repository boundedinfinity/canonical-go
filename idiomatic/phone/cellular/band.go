package cellular

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/location/region"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/boundedinfinity/go-commoner/idiomatic/mather"
)

type Band struct {
	Name          string
	Alias         optioner.Option[string]
	Frequency     float32
	UplinkSpeed   mather.UpperLowerRange[float32]
	DownlinkSpeed mather.UpperLowerRange[float32]
	Regions       optioner.Option[[]region.Region]
}

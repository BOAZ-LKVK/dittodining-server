package fiberfx

import (
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"go.uber.org/fx"
)

func AsAPIController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(apicontroller.APIController)),
		fx.ResultTags(`group:"api_controllers"`),
	)
}

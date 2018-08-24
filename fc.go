package fc
import "github.com/j-griffith/csi-connectors/fibrechannel"

type fcDevice struct {
	*fibrechannel.Connector
	disk string
	dm string
}

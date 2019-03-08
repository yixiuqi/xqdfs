package plugin

const(
	HttpTextPlain			=	"text/plain"
	HttpApplicationJson	=	"application/json"
)

type Invocation struct {
	Body []byte
	Method string
	ContentType string
}
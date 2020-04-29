package maps

type Maps struct {
	URL    string
	ApiKey string
}

var (
	apiKey string
	url    string
)

func Init(maps Maps) {
	apiKey = maps.ApiKey
	url = maps.URL
}

package fetchers

type Fetcher interface {
	Fetch(ticker string) error
}

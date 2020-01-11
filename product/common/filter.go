package common

import "net/http"

type FilterHandle func(w http.ResponseWriter, r *http.Request) error

type Filter struct {
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{make(map[string]FilterHandle)}
}

func (f *Filter) RegisterFilterUri(url string, handler FilterHandle) {
	f.filterMap[url] = handler
}

func (f *Filter) GetFilterHandle(url string) FilterHandle {
	return f.filterMap[url]
}

type WebHandle func(w http.ResponseWriter, r *http.Request)

func (f *Filter) Handle(webHandle WebHandle) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for url, handle := range f.filterMap {
			if url == r.RequestURI {
				err := handle(w, r)
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
				break
			}
		}

		webHandle(w, r)
	}
}

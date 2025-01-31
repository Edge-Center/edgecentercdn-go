package resources

import (
	"fmt"
	"net/url"
	"strings"
)

func buildListPath(offset uint, size uint, filter *ListFilterRequest) string {
	queryParams := url.Values{}

	if offset > 0 {
		queryParams.Set("offset", fmt.Sprintf("%d", offset))
	}

	if size > 0 {
		queryParams.Set("size", fmt.Sprintf("%d", size))
	}

	if filter != nil {
		if len(filter.Fields) > 0 {
			queryParams.Set("fields", strings.Join(filter.Fields, ","))
		}

		if filter.Ordering != "" {
			queryParams.Set("ordering", filter.Ordering)
		}

		if filter.Search != "" {
			queryParams.Set("search", filter.Search)
		}

		if filter.Deleted {
			queryParams.Set("deleted", "true")
		}

		if len(filter.Status) > 0 {
			statusValues := make([]string, len(filter.Status))
			for i, status := range filter.Status {
				statusValues[i] = strings.ToLower(string(status))
			}
			queryParams.Set("status", strings.Join(statusValues, ","))
		}

		if filter.Cname != "" {
			queryParams.Set("cname", filter.Cname)
		}
	}

	if len(queryParams) > 0 {
		return fmt.Sprintf("/cdn/resources?%s", queryParams.Encode())
	}

	return "/cdn/resources"
}

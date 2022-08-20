package datastruct

type SCRelatedPage []struct {
	Elements []struct {
		Elements []struct {
			Elements []struct {
				Elements []struct {
					Elements []struct {
						Attributes struct {
							AriaLabel string `json:"aria-label"`
						} `json:"attributes"`
					} `json:"elements"`
				} `json:"elements,omitempty"`
			} `json:"elements"`
		} `json:"elements"`
	}
}

package apiv2

func nilOnEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nilOnNilOrEmpty(p *string) *string {
	if p == nil {
		return nil
	} else if *p == "" {
		return nil
	}
	return p
}

func emptyOnNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

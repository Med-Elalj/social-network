package structs

// var categoriesRegex = regexp.MustCompile(`¤\w{3,20}`)

// func (p Ptitle) IsValid() error {
// 	p = Ptitle(html.EscapeString(strings.TrimSpace(string(p))))
// 	if len(p) < 3 || len(p) > 30 {
// 		return errors.New("title must be between 3 and 30 characters")
// 	}
// 	return nil
// }

// func (p *Pbody) IsValid() error {
// 	*p = Pbody(html.EscapeString(strings.TrimSpace(string(*p))))
// 	if len(*p) < 10 || len(*p) > 1024 {
// 		return errors.New("content must be between 10 and 1500 characters")
// 	}
// 	return nil
// }

// func (c *Pcategories) IsValid() error {
// 	return nil
// }

// func (p PostCreate) ParseCategories() (Pcategories, error) {
// 	categories := categoriesRegex.FindAllString(string(p.Content), -1)
// 	if len(categories) == 0 {
// 		return Pcategories{}, nil
// 	}
// 	for i, cat := range categories {
// 		categories[i] = cat[1:]
// 	}
// 	return categories, nil
// }

// func (p ID) IsValid() error {
// 	if p <= 0 {
// 		return errors.New("post ID must be a positive integer")
// 	}
// 	return nil
// }

// func (c CommentContent) IsValid() error {
// 	c = CommentContent(html.EscapeString(strings.TrimSpace(string(c))))
// 	if len(c) < 1 || len(c) > 300 {
// 		return errors.New("comment must be between 1 and 300 characters")
// 	}
// 	return nil
// }

// func (p PostCreate) Validate() error {
// 	return validateStruct(p)
// }

// func (c CommentInfo) Validate() error {
// 	return validateStruct(c)
// }

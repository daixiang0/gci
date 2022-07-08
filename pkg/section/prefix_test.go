package section

// func TestPrefixSpecificity(t *testing.T) {
// 	testCases := []specificityTestData{
// 		{`"foo/pkg/bar"`, Prefix{"", nil, nil}, specificity.MisMatch{}},
// 		{`"foo/pkg/bar"`, Prefix{"foo", nil, nil}, specificity.Match{3}},
// 		{`"foo/pkg/bar"`, Prefix{"bar", nil, nil}, specificity.MisMatch{}},
// 		{`"foo/pkg/bar"`, Prefix{"github.com/foo/bar", nil, nil}, specificity.MisMatch{}},
// 		{`"foo/pkg/bar"`, Prefix{"github.com/foo", nil, nil}, specificity.MisMatch{}},
// 		{`"foo/pkg/bar"`, Prefix{"github.com/bar", nil, nil}, specificity.MisMatch{}},
// 	}
// 	testSpecificity(t, testCases)
// }

// func TestPrefixParsing(t *testing.T) {
// 	testCases := []sectionTestData{
// 		{"pkgPREFIX", Custom{"", nil, nil}, nil},
// 		{"prefix(test.com)", Custom{"test.com", nil, nil}, nil},
// 	}
// 	testSectionParser(t, testCases)
// }

// func TestPrefixToString(t *testing.T) {
// 	testSectionToString(t, Custom{})
// 	testSectionToString(t, Custom{"", nil, nil})
// 	testSectionToString(t, Custom{"abc.org", nil, nil})
// 	testSectionToString(t, Custom{"abc.org", nil, CommentLine{"a"}})
// 	testSectionToString(t, Custom{"abc.org", CommentLine{"a"}, NewLine{}})
// }

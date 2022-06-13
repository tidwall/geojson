package geojson

func parseWKB(data string, opts *ParseOptions) (Object, error) {
	// if len(data) == 0 {
	// 	// 0x00 or 0x01 must be the first bytes
	// 	return nil, errDataInvalid
	// }
	// well-known binary is not supported yet
	return nil, errDataInvalid
}

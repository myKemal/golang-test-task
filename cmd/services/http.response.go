package services

func CreateHttpStatusCode(err error) int {
	//need improve for custom error
	if err != nil {
		return 400
	} else {
		return 200
	}
}

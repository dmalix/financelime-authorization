package system

//noinspection GoSnakeCaseUsage
type MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var MockData MockDescription

func (a *MockDescription) version() (string, string, error) {
	return "version", "build", MockData.Expected.Error
}

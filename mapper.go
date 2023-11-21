package goutils

func MapSlice[I interface{}, O interface{}](inputData []I, mapper func(I) O) []O {
	if inputData == nil {
		return nil
	}
	outputData := make([]O, len(inputData))
	for i := 0; i < len(inputData); i++ {
		outputData[i] = mapper(inputData[i])
	}
	return outputData
}

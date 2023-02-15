package db

import (
	"strconv"
	"strings"
)

const separator = ";"

func textPositionToXY(textPosition string) (int, int, error) {
	arr := strings.Split(textPosition, separator)

	x, err := strconv.Atoi(arr[0])
	if err != nil {
		return -1, -1, err
	}

	y, err := strconv.Atoi(arr[1])
	if err != nil {
		return -1, -1, err
	}

	return x, y, nil
}

func XYPositionToText(x, y int) string {
	return strconv.Itoa(x) + separator + strconv.Itoa(y)
}

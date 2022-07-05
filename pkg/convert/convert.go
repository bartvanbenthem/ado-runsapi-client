package convert

import (
	"log"
	"strconv"
)

func StringToInt32(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}

	return int32(i)
}

func StringPointerToInt32(s *string) int32 {
	i, err := strconv.Atoi(*s)
	if err != nil {
		log.Fatal(err)
	}

	return int32(i)
}

func StringPointerToBool(s *string) bool {
	i, err := strconv.ParseBool(*s)
	if err != nil {
		log.Fatal(err)
	}

	return bool(i)
}

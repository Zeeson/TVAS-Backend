package utils

import (
	"reflect"
)

func Contains(s []string, searchterms []string) bool {
	var check = false
	for j, _ := range s {
		check = in_array(s[j], searchterms)
		if !check {
			return false
		}
	}
	return check
}

func in_array(val interface{}, array interface{}) (exists bool) {
    exists = false

    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                exists = true
                return
            }
        }
    }

    return
}

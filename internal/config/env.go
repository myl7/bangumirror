package config

import (
	"os"
	"strconv"
)

var BangumiApiHost = func() string {
	if e := os.Getenv("BANGUMI_API_HOST"); e != "" {
		return e
	} else {
		return "https://mirror.api.bgm.rincat.ch"
	}
}()

var JobNum = func() int {
	if e := os.Getenv("JOB_NUM"); e != "" {
		if n, err := strconv.Atoi(e); err == nil {
			return n
		}
	}
	return 12
}()

var MongoDb = func() string {
	if e := os.Getenv("MONGO_DB"); e != "" {
		return e
	} else {
		return "bangumi"
	}
}()

var MongoSubjectColl = func() string {
	if e := os.Getenv("MONGO_SUBJECT_COLL"); e != "" {
		return e
	} else {
		return "subjects"
	}
}()

var SubjectStart = func() int {
	if e := os.Getenv("SUBJECT_START"); e != "" {
		if n, err := strconv.Atoi(e); err == nil {
			return n
		}
	}
	return 1
}()

var SubjectEnd = func() int {
	if n, err := strconv.Atoi(os.Getenv("SUBJECT_END")); err == nil {
		return n
	} else {
		panic(err)
	}
}()

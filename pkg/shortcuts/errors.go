package shortcuts

import "log"

func FatalIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

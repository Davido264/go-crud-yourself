package assert

func Assert(cond bool) {
	if !cond {
		panic("Assersion failed")
	}
}

func AssertErrNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

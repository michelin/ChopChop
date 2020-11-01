package core_test

import "gochopchop/core"

func createInt32(x int32) *int32 {
	return &x
}

var fakeCheckStatus200 = &core.Check{
	StatusCode: createInt32(200),
}

var fakeCheckStatus404 = &core.Check{
	StatusCode: createInt32(404),
}

var fakeCheckStatus500 = &core.Check{
	StatusCode: createInt32(500),
}

var fakePlugin200 = &core.Plugin{
	Endpoint: "/200",
	Checks: []*core.Check{
		fakeCheckStatus200,
	},
}

var fakePlugin404 = &core.Plugin{
	Endpoint: "/404",
	Checks: []*core.Check{
		fakeCheckStatus404,
	},
}

var fakePlugin500 = &core.Plugin{
	Endpoint: "/500",
	Checks: []*core.Check{
		fakeCheckStatus500,
	},
}

var FakeSignatures = &core.Signatures{
	Plugins: []*core.Plugin{
		fakePlugin200,
		fakePlugin404,
		fakePlugin500,
	},
}

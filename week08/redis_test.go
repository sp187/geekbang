package week08

import (
	"testing"
)

var (
	string10, string20, string50, string100, string200, string1k, string5k string
)

func TestMain(m *testing.M) {
	string10 = randomStr(10)
	string20 = randomStr(20)
	string50 = randomStr(50)
	string100 = randomStr(100)
	string200 = randomStr(200)
	string1k = randomStr(1024)
	string5k = randomStr(5120)
	InitClient()
	m.Run()
	FlushAll()
}

func BenchmarkRedisSet10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetKey("10", string10)
	}
}

func BenchmarkRedisSet20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetKey("20", string20)
	}
}

func BenchmarkRedisSet50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetKey("50", string50)
	}
}

func BenchmarkRedisSet100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetKey("100", string100)
	}
}

func BenchmarkRedisSet200(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetKey("200", string200)
	}
}

func BenchmarkRedisSet1k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetKey("1k", string1k)
	}
}

func BenchmarkRedisSet5k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetKey("5k", string5k)
	}
}

func BenchmarkRedisGet10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey("10")
	}
}

func BenchmarkRedisGet20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey("20")
	}
}

func BenchmarkRedisGet50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey("50")
	}
}

func BenchmarkRedisGet100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey("100")
	}
}

func BenchmarkRedisGet200(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey("200")
	}
}

func BenchmarkRedisGet1k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey("1k")
	}
}

func BenchmarkRedisGet5k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey("5k")
	}
}

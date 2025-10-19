package goutils

import (
	"os"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

const testEnvName = "TEST_ENV"

func TestEnv(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		t.Run("int", func(t *testing.T) {
			env := Env("SOME_NOT_EXISTING_ENV", 123)
			assert2.Equal(t, 123, env)
		})

		t.Run("string", func(t *testing.T) {
			env := Env("SOME_NOT_EXISTING_ENV", "myStr")
			assert2.Equal(t, "myStr", env)
		})

		t.Run("bool", func(t *testing.T) {
			env := Env("SOME_NOT_EXISTING_ENV", true)
			assert2.Equal(t, true, env)
		})

	})

	t.Run("defined value", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			envVal := "xyz"
			err := os.Setenv(testEnvName, envVal)
			assert2.NoError(t, err)
			env := Env(testEnvName, "myStr")
			assert2.Equal(t, envVal, env)
		})

		t.Run("int", func(t *testing.T) {
			err := os.Setenv(testEnvName, "123")
			assert2.NoError(t, err)
			env := Env(testEnvName, 0)
			assert2.Equal(t, 123, env)
		})

		t.Run("float32", func(t *testing.T) {
			err := os.Setenv(testEnvName, "3")
			assert2.NoError(t, err)
			env := Env(testEnvName, float32(0))
			assert2.Equal(t, float32(3), env)
		})

		t.Run("bool", func(t *testing.T) {
			err := os.Setenv(testEnvName, "1")
			assert2.NoError(t, err)
			env := Env(testEnvName, false)
			assert2.Equal(t, true, env)
		})
	})
}

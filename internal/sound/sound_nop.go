//go:build !windows && !darwin

package sound

// Init is a no-op on platforms where the realtime sound backend
// isn't enabled in this project build.
func Init() error {
	return nil
}

// PlayClick is intentionally a no-op on unsupported platforms.
func PlayClick() {}

// PlayError is intentionally a no-op on unsupported platforms.
func PlayError() {}

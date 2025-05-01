package types

type ProtectV1 interface {
	// About application
	Info() (*ProtectInfo, error)

	// FIXME: Fill-in rest.
}

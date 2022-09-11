package obfuscate

type Obfuscator interface {
	Scan(string) (string, error)
}

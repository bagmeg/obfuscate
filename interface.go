package obfuscate

type Obfuscator interface {
	Tokenize(string)
	Parse()
}

package ciphers

type CipherResponse struct {
	Text  string `json:"text,omitempty"`
	SecretKey  []byte `json:"secret_key,omitempty"`
}

type CipherRequest struct {
	TextDecrypt  string `json:"text_decrypt"`
}



package models

// asi se guarda el usuario en la base (la password es el hash de bcrypt)
type Usuario struct {
	Id               int    `json:"id,omitempty"`
	Nombre           string `json:"nombre"`
	Correo           string `json:"correo"`
	Password         string `json:"password,omitempty"`
	Rol              string `json:"rol"`
	CorreoVerificado bool   `json:"correoVerificado"`
	FechaRegistro    string `json:"fechaRegistro,omitempty"`
}

// lo que el cliente guarda de la sesion
type UsuarioAuth struct {
	Nombre string `json:"nombre"`
	Correo string `json:"correo"`
	Rol    string `json:"rol"`
}

type CambioRol struct {
	Rol string `json:"rol"`
}

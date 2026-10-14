package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"api_mid_the_house_fit/models"
	"api_mid_the_house_fit/services"

	beego "github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	beego.Controller
	authService *services.AuthService
}

func (c *AuthController) Prepare() {
	c.authService = services.NewAuthService()
}

// @Title Registro de usuario
// @Description Registra un nuevo usuario en el sistema
// @Param	body	body 	models.UsuarioRegistro	true	"Datos de registro"
// @Success 201 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 409 {object} models.RespuestaError
// @router /registro [post]
func (c *AuthController) Registro() {
	var req models.UsuarioRegistro
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	usuario, err := c.authService.RegistrarUsuario(req)
	if err != nil {
		if err == services.ErrCorreoYaRegistrado {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusConflict)
			c.Data["json"] = models.RespuestaErrorGeneral("El correo ya está registrado", err)
			c.ServeJSON()
			return
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al registrar usuario", err)
		c.ServeJSON()
		return
	}

	jwtService := services.GetJWTService()
	accessToken, refreshToken, _ := jwtService.GenerateTokenPair(usuario.ID, usuario.Correo, usuario.Rol)

	response := models.RespuestaExitosa("Usuario registrado exitosamente", map[string]interface{}{
		"usuario": models.UsuarioResponse{
			ID:               usuario.ID,
			Nombre:           usuario.Nombre,
			Correo:           usuario.Correo,
			Rol:              usuario.Rol,
			CorreoVerificado: usuario.CorreoVerificado,
			FechaRegistro:    usuario.FechaRegistro.Format(time.RFC3339),
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Enviar código de verificación
// @Description Envía un código de verificación o recuperación al correo
// @Param	body	body 	models.CodigoVerificacionRequest	true	"Correo y tipo"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @router /enviar-codigo [post]
func (c *AuthController) EnviarCodigo() {
	var req models.CodigoVerificacionRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	codigo, err := c.authService.GenerarCodigoVerificacion(req.Correo, req.Tipo)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al enviar código", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Código enviado correctamente", map[string]interface{}{
		"mensaje": "Código enviado al correo",
		"codigo":  codigo, // Solo en desarrollo
	})
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Verificar código
// @Description Verifica el código de verificación o recuperación
// @Param	body	body 	models.VerificarCodigoRequest	true	"Correo, código y tipo"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @router /verificar-codigo [post]
func (c *AuthController) VerificarCodigo() {
	var req models.VerificarCodigoRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	err := c.authService.VerificarCodigo(req.Correo, req.Codigo, req.Tipo)
	if err != nil {
		if err == services.ErrCodigoInvalido || err == services.ErrCodigoYaUsado {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			c.Data["json"] = models.RespuestaErrorGeneral("Código inválido o expirado", err)
			c.ServeJSON()
			return
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al verificar código", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Código verificado correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Iniciar sesión
// @Description Inicia sesión de usuario
// @Param	body	body 	models.UsuarioLogin	true	"Credenciales"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @router /login [post]
func (c *AuthController) Login() {
	var req models.UsuarioLogin
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	usuario, err := c.authService.Login(req)
	if err != nil {
		if err == services.ErrCredencialesInvalidas {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			c.Data["json"] = models.RespuestaErrorGeneral("Credenciales inválidas", err)
			c.ServeJSON()
			return
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al iniciar sesión", err)
		c.ServeJSON()
		return
	}

	jwtService := services.GetJWTService()
	accessToken, refreshToken, _ := jwtService.GenerateTokenPair(usuario.ID, usuario.Correo, usuario.Rol)

	response := models.RespuestaExitosa("Sesión iniciada correctamente", map[string]interface{}{
		"usuario": models.UsuarioResponse{
			ID:               usuario.ID,
			Nombre:           usuario.Nombre,
			Correo:           usuario.Correo,
			Rol:              usuario.Rol,
			CorreoVerificado: usuario.CorreoVerificado,
			FechaRegistro:    usuario.FechaRegistro.Format(time.RFC3339),
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Cerrar sesión
// @Description Cierra la sesión del usuario
// @Param	Authorization	header	string	true	"Bearer token"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @router /logout [post]
func (c *AuthController) Logout() {
	authHeader := c.Ctx.Input.Header("Authorization")
	if authHeader == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("Token requerido", nil)
		c.ServeJSON()
		return
	}

	// En una implementación completa, aquí invalidaríamos el token en BD
	response := models.RespuestaExitosa("Sesión cerrada correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Recuperar contraseña
// @Description Envía código de recuperación de contraseña
// @Param	body	body 	models.RecuperarPasswordRequest	true	"Correo"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @router /recuperar-password [post]
func (c *AuthController) RecuperarPassword() {
	var req models.RecuperarPasswordRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	_, err := c.authService.GenerarCodigoVerificacion(req.Correo, models.TipoCodigoRecuperacion)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al enviar código de recuperación", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Código de recuperación enviado al correo", nil)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Cambiar contraseña
// @Description Cambia la contraseña usando código de recuperación
// @Param	body	body 	models.CambiarPasswordRequest	true	"Datos para cambio"
// @Success 200 {object} models.RespuestaAPI
// @Failure 400 {object} models.RespuestaError
// @Failure 401 {object} models.RespuestaError
// @router /cambiar-password [post]
func (c *AuthController) CambiarPassword() {
	var req models.CambiarPasswordRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Datos inválidos", err)
		c.ServeJSON()
		return
	}

	if req.NuevaPassword != req.ConfirmarPassword {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = models.RespuestaErrorGeneral("Las contraseñas no coinciden", nil)
		c.ServeJSON()
		return
	}

	err := c.authService.CambiarPassword(req.Correo, req.Codigo, req.NuevaPassword)
	if err != nil {
		if err == services.ErrCodigoInvalido || err == services.ErrCodigoYaUsado {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			c.Data["json"] = models.RespuestaErrorGeneral("Código inválido o expirado", err)
			c.ServeJSON()
			return
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = models.RespuestaErrorGeneral("Error al cambiar contraseña", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Contraseña cambiada correctamente", nil)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Refrescar token
// @Description Refresca el token de acceso usando refresh token
// @Param	Authorization	header	string	true	"Bearer refresh_token"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @router /refresh [post]
func (c *AuthController) RefreshToken() {
	authHeader := c.Ctx.Input.Header("Authorization")
	if authHeader == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("Refresh token requerido", nil)
		c.ServeJSON()
		return
	}

	refreshToken := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		refreshToken = authHeader[7:]
	}

	jwtService := services.GetJWTService()
	accessToken, err := jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("Refresh token inválido o expirado", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Token refrescado", map[string]string{
		"access_token": accessToken,
	})
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title Obtener perfil
// @Description Obtiene el perfil del usuario autenticado
// @Param	Authorization	header	string	true	"Bearer token"
// @Success 200 {object} models.RespuestaAPI
// @Failure 401 {object} models.RespuestaError
// @router /perfil [get]
func (c *AuthController) Perfil() {
	usuarioID := c.Ctx.Input.GetData("usuario_id")
	if usuarioID == nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = models.RespuestaErrorGeneral("No autenticado", nil)
		c.ServeJSON()
		return
	}

	usuario, err := c.authService.ObtenerUsuarioPorID(usuarioID.(int))
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Data["json"] = models.RespuestaErrorGeneral("Usuario no encontrado", err)
		c.ServeJSON()
		return
	}

	response := models.RespuestaExitosa("Perfil obtenido", models.UsuarioResponse{
		ID:               usuario.ID,
		Nombre:           usuario.Nombre,
		Correo:           usuario.Correo,
		Rol:              usuario.Rol,
		CorreoVerificado: usuario.CorreoVerificado,
		FechaRegistro:    usuario.FechaRegistro.Format(time.RFC3339),
	})
	c.Data["json"] = response
	c.ServeJSON()
}
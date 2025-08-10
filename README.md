# Go API - Find My Friend

API REST desarrollada en Go con arquitectura por capas.

## Arquitectura

Este proyecto sigue una arquitectura por capas bien definida:

```
┌─────────────────┐
│     Routes      │  ← Manejo de rutas HTTP
├─────────────────┤
│   Controllers   │  ← Lógica de control de requests/responses
├─────────────────┤
│    Services     │  ← Lógica de negocio
├─────────────────┤
│  Repositories   │  ← Acceso a datos
└─────────────────┘
```

### Capas:

1. **Routes**: Define los endpoints de la API y conecta con los controllers
2. **Controllers**: Maneja las requests HTTP, valida datos y retorna responses
3. **Services**: Contiene la lógica de negocio de la aplicación
4. **Repositories**: Gestiona el acceso a la base de datos

## Estructura del Proyecto

```
├── cmd/server/          # Punto de entrada de la aplicación
├── internal/            # Código interno de la aplicación
│   ├── routes/         # Definición de rutas
│   ├── controllers/    # Controladores HTTP
│   ├── services/       # Lógica de negocio
│   ├── repositories/   # Acceso a datos
│   ├── models/         # Modelos de datos
│   └── middleware/     # Middleware personalizado
├── pkg/                # Paquetes reutilizables
│   ├── database/       # Configuración de base de datos
│   └── utils/          # Utilidades comunes
└── configs/            # Configuraciones
```

## Tecnologías

- **Go 1.21+**
- **Gin**: Framework web para routing
- **GORM**: ORM para base de datos
- **PostgreSQL**: Base de datos
- **godotenv**: Manejo de variables de entorno

## Instalación

1. Clonar el repositorio
2. Instalar dependencias: `go mod tidy`
3. Configurar variables de entorno en `.env`
4. Ejecutar: `go run cmd/server/main.go`

## Variables de Entorno

Crear un archivo `.env` con:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=find_my_friend
SERVER_PORT=8080
```

## Endpoints de la API

### Usuarios
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users` - Obtener todos los usuarios
- `GET /api/v1/users/:id` - Obtener usuario por ID
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario
- `GET /api/v1/users/search?name=john` - Buscar usuarios por nombre

### Mascotas
- `POST /api/v1/pets` - Crear mascota perdida
- `GET /api/v1/pets?sort&page&size` - Obtener mascotas con paginación y ordenamiento
- `GET /api/v1/pets/:id` - Obtener mascota por ID
- `PUT /api/v1/pets/:id` - Actualizar mascota
- `DELETE /api/v1/pets/:id` - Eliminar mascota
- `PUT /api/v1/pets/found` - Marcar mascota como encontrada
- `GET /api/v1/pets/search?q=query&page&size` - Buscar mascotas
- `GET /api/v1/pets/user/:user_id` - Obtener mascotas de un usuario

### Parámetros de Query
- `page`: Número de página (default: 1)
- `size`: Tamaño de página (default: 10, max: 100)
- `sort`: Ordenamiento (name, -name, type, -type, created_at, -created_at)
- `q`: Término de búsqueda para mascotas

### Ejemplos de Uso

#### Crear una mascota perdida:
```json
POST /api/v1/pets
{
  "name": "Luna",
  "type": "dog",
  "breed": "Golden Retriever",
  "user_id": "uuid-del-usuario",
  "last_seen": {
    "time": "2024-01-15T10:30:00Z",
    "place": "Parque Central, Calle 123"
  }
}
```

#### Marcar mascota como encontrada:
```json
PUT /api/v1/pets/found
{
  "pet_id": "uuid-de-la-mascota"
}
``` 
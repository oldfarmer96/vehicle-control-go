-- 1. Habilitar extensión para UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Función genérica para actualizar el campo updated_at (como lo hace Prisma)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 3. Crear ENUMs (Mapeo de las enumeraciones de Prisma)
CREATE TYPE rol_universitario AS ENUM ('DOCENTE', 'ALUMNO', 'ADMINISTRATIVO', 'VISITANTE');
CREATE TYPE tipo_evento_acceso AS ENUM ('ENTRADA', 'SALIDA');
CREATE TYPE rol_web AS ENUM ('ADMINISTRADOR', 'CONSULTOR');

-- 4. Crear Tablas

-- Tabla: usuarios_web (UsuarioWeb en Prisma)
CREATE TABLE usuarios_web (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nombre VARCHAR(50) NOT NULL,
    apellidos VARCHAR(50) NOT NULL,
    dni VARCHAR(15) UNIQUE NOT NULL,
    username VARCHAR(15) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    rol rol_web NOT NULL,
    activo BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_usuarios_web_created_at ON usuarios_web(created_at);

-- Tabla: vehiculos (Vehiculo en Prisma)
CREATE TABLE vehiculos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    placa VARCHAR(15) UNIQUE NOT NULL,
    marca VARCHAR(50),
    modelo VARCHAR(50),
    color VARCHAR(50),
    vin VARCHAR(50),
    motor VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_vehiculos_created_at ON vehiculos(created_at);

-- Tabla: personas (Persona en Prisma)
CREATE TABLE personas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dni VARCHAR(15) UNIQUE NOT NULL,
    nombre_completo VARCHAR(150) NOT NULL,
    rol rol_universitario NOT NULL,
    tiene_acceso_permitido BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_personas_nombre_completo ON personas(nombre_completo);
CREATE INDEX idx_personas_created_at ON personas(created_at);

-- Tabla Intermedia: vehiculos_personas (VehiculoPersona en Prisma)
CREATE TABLE vehiculos_personas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehiculo_id UUID NOT NULL,
    persona_id UUID NOT NULL,
    asignado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Llaves Foráneas con Cascade Delete (Como en Prisma)
    CONSTRAINT fk_vehiculo FOREIGN KEY (vehiculo_id) REFERENCES vehiculos(id) ON DELETE CASCADE,
    CONSTRAINT fk_persona FOREIGN KEY (persona_id) REFERENCES personas(id) ON DELETE CASCADE,
    
    -- Restricción Única: Un vehículo y una persona no pueden vincularse dos veces
    CONSTRAINT unq_vehiculo_persona UNIQUE (vehiculo_id, persona_id)
);
CREATE INDEX idx_vehiculos_personas_vehiculo_id ON vehiculos_personas(vehiculo_id);
CREATE INDEX idx_vehiculos_personas_persona_id ON vehiculos_personas(persona_id);
CREATE INDEX idx_vehiculos_personas_asignado_en ON vehiculos_personas(asignado_en);

-- Tabla: eventos_acceso (EventoAcceso en Prisma)
CREATE TABLE eventos_acceso (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehiculo_id UUID NOT NULL,
    tipo_evento tipo_evento_acceso NOT NULL,
    punto_control VARCHAR(50) NOT NULL,
    confianza_ocr DECIMAL(5, 2),
    fecha_hora TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_evento_vehiculo FOREIGN KEY (vehiculo_id) REFERENCES vehiculos(id) ON DELETE RESTRICT
);
CREATE INDEX idx_eventos_acceso_vehiculo_id ON eventos_acceso(vehiculo_id);
CREATE INDEX idx_eventos_acceso_fecha_hora ON eventos_acceso(fecha_hora);
CREATE INDEX idx_eventos_acceso_vehiculo_fecha ON eventos_acceso(vehiculo_id, fecha_hora);
CREATE INDEX idx_eventos_acceso_tipo_fecha ON eventos_acceso(tipo_evento, fecha_hora);


-- 5. Asignar Triggers para updated_at (Solo en tablas que lo requieren)
CREATE TRIGGER set_updated_at_usuarios_web
BEFORE UPDATE ON usuarios_web FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_vehiculos
BEFORE UPDATE ON vehiculos FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_personas
BEFORE UPDATE ON personas FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
# Makefile

# Directorio donde Go está instalado (GOROOT)
GOROOT=$(shell go env GOROOT)

# Directorio de destino para el archivo wasm_exec.js y el archivo compilado
OUT_DIR=.

# Nombre del archivo de salida WebAssembly
WASM_FILE=$(OUT_DIR)/rabbits.wasm

# Paquete de Go que deseas compilar
PKG=github.com/demonodojo/rabbits

compile:
    # Compila el paquete Go a WebAssembly
	env NODE_ENV=development GOOS=js GOARCH=wasm go build -o $(WASM_FILE) $(PKG)

    	# Copia wasm_exec.js al directorio actual
	cp $(GOROOT)/lib/wasm/wasm_exec.js $(OUT_DIR)

.PHONY: compile

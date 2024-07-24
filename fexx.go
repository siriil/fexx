package main

import (
    "flag"
    "net/http"
    "log"
    "strings"
    "os"
    "path/filepath"
    "mime"
)

// Start inicia un servidor HTTP que sirve archivos estaticos con sustitucion de parametros.
func Start(addr string, port string, dir string) {
    // Crear un manejador para todas las rutas
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Construir la ruta del archivo
        log.Printf("%s %s", r.Method, r.URL)
        filePath := filepath.Join(dir, r.URL.Path[1:])

        // Obtener informacion del archivo
        fileInfo, err := os.Stat(filePath)
        if err != nil {
            http.NotFound(w, r)
            return
        }

        // Leer el archivo
        fileContent, err := os.ReadFile(filePath)
        if err != nil {
            http.NotFound(w, r)
            return
        }

        // Realizar sustituciones basadas en parametros de consulta
        for key, values := range r.URL.Query() {
            for _, value := range values {
                // Reemplazar todas las ocurrencias
                placeholder := key // PATRON "${" + key + "}"
                fileContent = []byte(strings.ReplaceAll(string(fileContent), placeholder, value))
            }
        }

        // Establecer el tipo de contenido basado en la extension del archivo
        ext := filepath.Ext(filePath)
        contentType := mime.TypeByExtension(ext)
        if contentType == "" {
            contentType = "application/octet-stream"
        }
        w.Header().Set("Content-Type", contentType)
        lastModified := fileInfo.ModTime().Format(http.TimeFormat)
        w.Header().Set("Last-Modified", lastModified)
        w.Header().Set("Server", "FileServer FEXX")
        w.Write(fileContent)
    })

    // Imprimir un mensaje indicando que el servidor esta corriendo
    log.Printf("Directorio desde %s en http://%s:%s/", dir, addr, port)

    // Iniciar el servidor y escuchar en el puerto dado
    if err := http.ListenAndServe(addr+":"+port, nil); err != nil {
        log.Fatalf("Error al iniciar el servidor: %v", err)
    }
}

func main() {
    addr := flag.String("addr", "localhost", "La direccion en el que el servidor escucha")
    port := flag.String("port", "8000", "El puerto en el que el servidor escucha")
    dir := flag.String("dir", ".", "El directorio desde el cual se proveen los archivos")
    flag.Parse()

    Start(*addr, *port, *dir)
}

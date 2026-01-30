package config

import (
    "bufio"
    "os"
    "path/filepath"
    "strings"
)

const defaultEnv = "dev"

type Config struct {
    Env      string
    Greeting string
    APIKey   string
}

func Load() Config {
    env := getenv("APP_ENV", defaultEnv)
    moduleRoot := findModuleRoot()
    fileValues := readEnvFile(filepath.Join(moduleRoot, "config", env+".env"))
    greeting := firstNonEmpty(os.Getenv("GREETING"), fileValues["GREETING"], "hello")
    apiKey := firstNonEmpty(os.Getenv("API_KEY"), fileValues["API_KEY"], "")
    return Config{
        Env:      env,
        Greeting: greeting,
        APIKey:   apiKey,
    }
}

func readEnvFile(path string) map[string]string {
    file, err := os.Open(path)
    if err != nil {
        return map[string]string{}
    }
    defer file.Close()

    values := map[string]string{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        if !strings.Contains(line, "=") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        values[key] = value
    }
    return values
}

func getenv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func firstNonEmpty(values ...string) string {
    for _, value := range values {
        if value != "" {
            return value
        }
    }
    return ""
}

func findModuleRoot() string {
    dir, err := os.Getwd()
    if err != nil {
        return "."
    }
    for {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            return dir
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            return dir
        }
        dir = parent
    }
}

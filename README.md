# Feature-Flag-Service

Ein Feature-Flag-Service als REST-API in Go, ausschließlich auf Basis der
Standardbibliothek (`net/http`). Flags lassen sich anlegen, auflisten, abrufen,
ändern und löschen; ein Evaluate-Endpunkt trifft eine deterministische
Ja/Nein-Entscheidung pro Nutzer anhand eines stabilen Hashes gegen den
`rollout_percent`. Ein thread-sicherer In-Memory-Store (mit `sync.RWMutex`),
Eingabevalidierung mit sauberen Statuscodes und JSON-Fehlerobjekten sowie
Zugriffs-Logging als Middleware runden den Service ab.

## Tech-Stack

- **Sprache**: Go 1.22+
- **HTTP**: `net/http` (Standardbibliothek, kein externes Web-Framework)
- **Store**: In-Memory mit `sync.RWMutex`
- **Tests**: Go-Tests (`go test`)

## Installation

```bash
# Voraussetzung: Go 1.22 oder neuer
git clone <repository-url>
cd <repository>
```

Keine externen Abhängigkeiten — es muss nichts installiert werden.

## Ausführen

```bash
go run .
```

Der Server lauscht standardmäßig auf Port `8080`. Über die Umgebungsvariable
`PORT` lässt sich ein anderer Port wählen:

```bash
PORT=9090 go run .
```

Der HTTP-Server startet mit expliziten Timeouts (`ReadHeaderTimeout`,
`ReadTimeout`, `WriteTimeout` jeweils 5 s), um langsame Verbindungen zu beenden.

### Build

```bash
go build ./...
```

## Endpunkte

Alle Antworten sind JSON; jeder Fehler hat die Form `{"error":"<meldung>"}`.

| Methode | Pfad                        | Beschreibung                                            |
|---------|-----------------------------|---------------------------------------------------------|
| POST    | `/flags`                    | Legt ein Flag an (Body: `{key, enabled, description?, rollout_percent?}`) → `201` |
| GET     | `/flags`                    | Listet alle Flags → `200` (leer als `[]`)                |
| GET     | `/flags/{key}`              | Liefert ein Flag → `200` / `404`                         |
| PUT     | `/flags/{key}`              | Ändert `enabled`/`description`/`rollout_percent` → `200` / `404` |
| DELETE  | `/flags/{key}`              | Entfernt ein Flag → `204` / `404`                        |
| GET     | `/flags/{key}/evaluate?user={id}` | Deterministische Rollout-Entscheidung → `200`      |
| GET     | `/healthz`                  | Health-Check → `200` `{"status":"ok"}`                   |

### Flag-Objekt

```json
{
  "key": "my-feature",
  "enabled": true,
  "description": "optional",
  "rollout_percent": 100
}
```

`description` ist `omitempty`, `rollout_percent` liegt im Bereich 0–100.

## Features

- Thread-sicherer In-Memory-Store (Mutex-geschützt)
- Eingabevalidierung (`key`, `rollout_percent`, Body-Größenlimit 1 MiB)
- JSON-Fehlerobjekte mit sauberen Statuscodes
- Logging-Middleware (Methode, Pfad, Statuscode, Dauer — keine Nutzerdaten)
- Health-Endpunkt `/healthz`
- HTTP-Server mit expliziten Timeouts gegen langsame Verbindungen

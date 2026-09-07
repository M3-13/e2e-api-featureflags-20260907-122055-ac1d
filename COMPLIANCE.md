VERDICT: APPROVED

## 1. GDPR / Datenschutz

**Bewertung:** Für den Projekttyp `go-backend` entsteht das relevante personenbezogene Datum ausschließlich im Evaluate-Endpunkt: der Query-Parameter `user` in `GET /flags/{key}/evaluate?user={id}` kann einen Online-Identifikator und damit ein personenbezogenes Datum darstellen.

Die Datenschutz-Akzeptanzkriterien sind vollständig erfüllt.

- **AC-17 erfüllt**  
  `internal/middleware/logging.go` protokolliert ausschließlich:
  ```go
  log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.status, time.Since(start))
  ```
  `r.URL.Path` enthält keinen Query-String. Der `user`-Parameter erscheint nicht im Log. Das bestätigt auch `internal/middleware/logging_test.go` mit `TestLoggingOmitsQueryStringAndUser`.

- **AC-18 erfüllt**  
  `internal/store/store.go` speichert ausschließlich `Flag`-Datensätze mit `key`, `enabled`, `description`, `rollout_percent`. Kein `user`-Feld ist vorhanden. In `internal/api/evaluate.go` wird `user` nur als lokale Variable innerhalb der Anfrage verwendet und nicht im Store abgelegt.

- **AC-19 erfüllt**  
  `internal/api/helpers.go` zentralisiert Fehlerantworten über `writeError` mit festen, internen Meldungen. Die Evaluate-Fehlerpfade (`missing user query parameter`, `user must not exceed 128 characters`, `flag not found`) geben keine clientseitigen Inhalte zurück. Auch die übrigen Handler (`create.go`, `read.go`, `update.go`, `delete.go`) spiegeln keine übermittelten personenbezogenen Daten in Fehlerantworten.

**Hinweis (nicht blockierend):**  
Der Betreiber des Dienstes muss für die tatsächliche Nutzung eine Rechtsgrundlage nach Art. 6 DSGVO sicherstellen. Technisch verhält sich der Dienst datensparsam: Der `user`-Wert wird nicht gespeichert, nicht geloggt und nur transient verarbeitet. Eine Datenschutzerklärung/Privacy Policy ist für ein reines Backend ohne Endnutzer-UI nicht zwingend Teil dieses Repositories.

## 2. EU Cyber Resilience Act (CRA)

**Bewertung:** Für ein Produkt mit digitalen Elementen sind die sichtbaren Sicherheitsanforderungen gut umgesetzt:

- **AC-13 erfüllt**  
  `internal/api/helpers.go` begrenzt Request-Bodys über `http.MaxBytesReader` auf 1 MiB (`maxBodyBytes`). Überschreitungen führen zu `413`.

- **AC-14 erfüllt**  
  `main.go` setzt explizite Timeouts:
  ```go
  ReadHeaderTimeout: 5 * time.Second,
  ReadTimeout:       5 * time.Second,
  WriteTimeout:      5 * time.Second,
  ```

- **AC-15 erfüllt**  
  `validateKey` begrenzt den `key` auf 128 Zeichen und erlaubt nur `[A-Za-z0-9_-]`. Der `user`-Parameter wird in `EvaluateFlag` explizit auf maximal 128 Zeichen begrenzt.

- **AC-16 / AC-17 erfüllt**  
  Keine sensiblen Anfragedaten im Logging.

**Hinweise (nicht blockierend):**  
Eine SBOM, ein dokumentierter Update-/Patch-Prozess oder ein separates Security-Properties-Dokument sind im sichtbaren Quellcode nicht vorhanden. Dafür gibt es in der Spec jedoch kein eigenes `[CRA]`-Akzeptanzkriterium; die sichtbaren Abhängigkeiten beschränken sich auf die Go-Standardbibliothek (`go.mod` hat keinen externen Dependency-Block im sichtbaren Stand), was das Lieferkettenrisiko gering hält. Für eine spätere Marktreife sollte eine SBOM ergänzt werden.

## 3. EU AI Act

**Nicht einschlägig.** Es ist keine KI-Funktion vorhanden. Der Service führt lediglich eine deterministische Hash-basierte Rollout-Entscheidung aus.

## 4. Pflichttexte & UI

**Nicht einschlägig.** Als reiner Go-Backend-Dienst ohne öffentliche Web-UI entstehen keine Pflichten zu Impressum, Cookie-Banner, Nutzungsbedingungen oder verbraucherspezifischem Widerrufshinweis.

## 5. Barrierefreiheit

**Nicht einschlägig.** Es gibt keine öffentliche, von Endnutzern bediente Weboberfläche.

## Gesamtergebnis

Keine offenen rechtlichen Blocker. Die datenschutzrelevanten Akzeptanzkriterien AC-17, AC-18 und AC-19 sowie die Sicherheitskriterien AC-13 bis AC-16 sind im sichtbaren, gemergten Produktstand vollständig umgesetzt.
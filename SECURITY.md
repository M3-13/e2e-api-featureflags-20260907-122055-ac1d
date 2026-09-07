VERDICT: APPROVED

## Sicherheitsbericht

### Zusammenfassung
Der Feature-Flag-Service erfüllt die vorgesehenen Sicherheits- und Datenschutzkriterien AC-13 bis AC-19. Es wurden keine ausnutzbaren Schwachstellen mit hohem oder kritischem Risiko festgestellt. Die manuelle Codeanalyse ergab keine Verstöße gegen die vereinbarten Sicherheitsanforderungen.

### Prüfbereiche

**Secrets**  
Keine hartkodierten Schlüssel, Passwörter, Tokens oder Credentials im Code. Die Port-Konfiguration erfolgt über `PORT` aus der Umgebung mit sicherem Default `8080`. Keine Secret-Protokollierung.

**Injection & Eingaben**  
- JSON-Bodies werden in `internal/api/helpers.go` über `http.MaxBytesReader` auf 1 MiB begrenzt (`maxBodyBytes`) und bei Überschreitung mit 413 beantwortet. Das erfüllt AC-13.
- Der `key`-Parameter wird in `validateKey` auf `[A-Za-z0-9_-]` und maximal 128 Zeichen geprüft. Das verhindert Path-Traversal und erfüllt AC-15.
- Der `user`-Queryparameter wird auf 128 Zeichen begrenzt (AC-15).  
- Keine SQL-, Command- oder Pfadinjektion möglich, da ausschließlich In-Memory-Store und Standardbibliothek verwendet werden.

**Authentifizierung / Autorisierung**  
Der Dienst hat keine Authentifizierung oder Autorisierung. Da die Spec hierfür kein Kriterium enthält, wird dies als nicht blockierender Hinweis unter „Notes“ geführt, nicht als Befund.

**Abhängigkeiten**  
Das Projekt verwendet ausschließlich die Go-Standardbibliothek (`go.mod` ohne externe Dependencies). Es sind keine bekannten verwundbaren Pakete vorhanden. Es wurde kein Scanner ausgeführt; die Analyse beruht auf manueller Inspektion.

**Konfiguration & Transport**  
In `main.go` ist der HTTP-Server mit expliziten Timeouts gestartet:  
`ReadHeaderTimeout: 5s`, `ReadTimeout: 5s`, `WriteTimeout: 5s`.  
Das erfüllt AC-14. Es sind keine unsicheren Defaults, offenen Debug-Einstellungen oder übermäßig breiten CORS-Konfigurationen sichtbar.

**Datenschutz & Logging**  
Die Logging-Middleware in `internal/middleware/logging.go` protokolliert ausschließlich Methode, Pfad ohne Query-String, Statuscode und Dauer. Query-Parameter wie `user` sowie Request-Bodys werden nicht geloggt. Das erfüllt AC-16 und AC-17.  
Der In-Memory-Store enthält nur Flag-Datensätze (`key`, `enabled`, `description`, `rollout_percent`); der `user`-Parameter wird nicht gespeichert (AC-18).  
Fehlerantworten verwenden feste, interne Meldungen ohne Clientdaten; insbesondere wird der `user`-Parameter nicht in Fehlerobjekten reflektiert (AC-19).

### Findings (blockierend oder mittel)
Keine. Es wurden keine Verstöße gegen die vereinbarten Sicherheitskriterien festgestellt.

### Notes (non-blocking)
- **Fehlende Authentifizierung/Autorisierung**: Jeder, der den Service erreicht, kann Flags anlegen, ändern und löschen. Dies ist nicht durch ein AC-Kriterium abgedeckt, sollte aber vor einem Produktivbetrieb ergänzt werden (z. B. API-Key, mTLS oder ein vorgelagertes Gateway).
- **Kein TLS**: Der Server lauscht auf HTTP ohne TLS. Ebenfalls kein AC-Kriterium; im Produktivbetrieb sollte TLS terminiert werden.
- **Kein Rate-Limiting / Begrenzung der Flag-Anzahl**: Ein Angreifer könnte beliebig viele Flags mit großen `description`-Feldern anlegen und so den In-Memory-Speicher füllen. Nicht durch die Spec abgedeckt.
- **Scanner-Hinweis**: Es wurde kein automatischer Sicherheitsscanner für das Go-Projekt ausgeführt; die Bewertung stützt sich auf manuelle Codeanalyse.
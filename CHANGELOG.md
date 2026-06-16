# Changelog

Todos los cambios relevantes de este proyecto se documentan en este archivo.

El formato se basa en [Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/)
y el proyecto sigue [Versionado Semántico](https://semver.org/lang/es/).

## [Unreleased]

## [0.1.0] - 2026-06-16

Primera versión pública.

### Añadido

- TUI de resultados de fútbol construida con Bubble Tea, inspirada en promiedos.com.ar.
- **Resultados por fecha**: partidos del día agrupados por liga, con colores de equipo,
  marcador y estado. Navegación de fechas (`←→`), `t` para volver a hoy.
- Contador `● N EN VIVO` en la cabecera y auto-recarga cada 15 s para partidos en juego.
- **Tabla de posiciones** con zonas de clasificación por color y soporte de grupos/fases.
- **Fixture** por liga: listado completo de partidos por fecha, navegable por ronda.
- **Detalle de partido**: marcador, estado, estadio/árbitro y canales de TV.
- Barra lateral con las ligas del día; `enter` abre liga (cabecera/sidebar) o partido.
- Flag `--version` y `--help`.
- Distribución: CI con GitHub Actions, binarios multiplataforma con GoReleaser,
  script de instalación y cask de Homebrew.

[Unreleased]: https://github.com/ianaya89/termiedos/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/ianaya89/termiedos/releases/tag/v0.1.0

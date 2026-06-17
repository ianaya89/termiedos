# Changelog

Todos los cambios relevantes de este proyecto se documentan en este archivo.

El formato se basa en [Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/)
y el proyecto sigue [Versionado Semántico](https://semver.org/lang/es/).

## [Unreleased]

## [0.2.0] - 2026-06-16

### Añadido

- **Goleadores**: lista de autores de gol por equipo con el minuto en el detalle de partido.
- **Tarjetas**: cronología de amonestaciones (amarilla, roja, doble amarilla) con minuto y equipo.
- **Estadísticas**: posesión, remates, faltas y demás métricas como barras proporcionales.
- **Forma reciente**: últimos resultados de cada equipo como fichas G/E/P.
- **Formaciones**: onces titulares de ambos equipos lado a lado, con dorsal, capitán (C),
  director técnico y marcadores de gol/tarjeta/cambio.
- **Historial**: balance head-to-head y enfrentamientos previos con su marcador.
- Scroll en el detalle de partido (`↑↓`, `g`/`G`, `pgup`/`pgdn`) con indicador de desbordamiento.

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

[Unreleased]: https://github.com/ianaya89/termiedos/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/ianaya89/termiedos/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/ianaya89/termiedos/releases/tag/v0.1.0

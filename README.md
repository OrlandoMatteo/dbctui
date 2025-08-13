# DBCtui

**DBCtui** is a terminal-based user interface (TUI) for exploring and searching DBC (CAN database) files. It allows you to quickly search, browse, and inspect CAN messages and signals defined in a DBC file, all from your terminal.

## Features

- **Fast Search:** Instantly filter signals and messages by name as you type.
- **Split View:** See a list of results and detailed information side-by-side.
- **Pagination:** Navigate large DBC files with page controls.
- **Keyboard Navigation:** Use arrow keys, tab, and shortcuts for efficient browsing.
- **Signal & Message Details:** View all key properties of signals and messages, including IDs, bit positions, value ranges, and more.

## Usage

### Build & Run

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/dbctui.git
   cd dbctui
   ```

2. **Build the project:**
   ```sh
   go build
   ```

3. **Run with a DBC file:**
   ```sh
   ./dbctui path/to/your/file.dbc
   ```

### Keyboard Shortcuts

- `↑` / `↓` or `k` / `j`: Move selection up/down
- `←` / `→`: Change page
- `Tab`: Switch between Signals and Messages
- `q` or `Ctrl+C`: Quit
- `Enter`: Select (future use)

## Project Structure

- [`main.go`](main.go): Entry point, loads the DBC file and starts the TUI.
- [`dbc/dbc.go`](dbc/dbc.go): DBC file parsing logic.
- [`can/`](can/): CAN domain models:
  - [`message.go`](can/message.go): CAN message structure.
  - [`signal.go`](can/signal.go): CAN signal structure.
  - [`state.go`](can/state.go): Signal state definitions.
  - [`problems.go`](can/problems.go): Problem/warning definitions.
- [`ui/model.go`](ui/model.go): TUI logic, state management, and rendering.

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling

## License

MIT License

---

*Made with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for CAN bus engineers and enthusiasts.*
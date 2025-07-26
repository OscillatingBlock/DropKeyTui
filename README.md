# DropKeyTui

**DropKeyTui** is a terminal-based TUI application for securely fetching, decrypting, and viewing end-to-end encrypted notes from a DropKey backend. Built with Go and BubbleTea, it leverages AES encryption and Ed25519 digital signatures for robust security and data integrity.

---

## Features

- End-to-end encryption: Securely retrieve and decrypt notes.
- Signed content verification: Ensure data authenticity with Ed25519 signatures.
- Syntax Highlighting : View notes with clear syntax Highlighting.
- Intuitive navigation: Search for pastes by ID with a user-friendly interface.
- Error handling: Gracefully manages expired, invalid, or tampered content.
- Client-side verification: All cryptographic operations occur locally.
- Configurable key pairs: Can use custom public/private keys for registration.

---

## Cryptography

- **Encryption**: AES-GCM (Galois/Counter Mode) using Go’s `crypto/aes` and `crypto/cipher` packages for secure data encryption.
- **Digital Signatures**: Ed25519 signatures via Go’s `crypto/ed25519` for authenticity and integrity.
- **Data Format**: Encrypted and signed JSON blobs containing `title` and `paste` fields, served by the DropKey backend.
- **Security**: All cryptographic operations are performed client-side, ensuring no unencrypted data is exposed to the backend.

---

## Getting Started

### Prerequisites

- **Go**: Version 1.21 or higher.
- **Terminal**: A terminal emulator with ANSI color support (e.g., iTerm2, Alacritty, or Windows Terminal).
- **Optional**: A `config.json` file with base64-encoded public/private key pairs for custom registration.

### Installation

1. **Clone the repository**:
   ```bash
    git clone git@github.com:OscillatingBlock/DropKeyTui.git

   cd DropKeyTui
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the TUI**:
   ```bash
   go run main.go
   ```

### Optional Configuration

To use custom key pairs for registration:

1. Generate an Ed25519 key pair (e.g., using a tool like `openssl` or programmatically).
2. Create a `config.json` file:
   ```json
   {
     "public_key": "base64-encoded-public-key",
     "private_key": "base64-encoded-private-key"
   }
   ```
3. During registration, the application will prompt you to specify the path to your `config.json` file, which it will load automatically.

---

## Backend Integration

DropKeyTui connects to the **DropKey backend**, a RESTful Go-based API for storing and retrieving encrypted pastes. The backend ensures secure storage, while the TUI handles all decryption and verification.

**Backend Repository**: [DropKey Backend](https://github.com/OscillatingBlock/DropKey) 

---

## Screenshots

Below are example interfaces of DropKeyTui:

- **Landing page**:
- 
  <img width="1894" height="957" alt="image" src="https://github.com/user-attachments/assets/347764df-c159-4168-9695-ad4bd7d8d86c" />

- **Syntax Highlight**: Syntax highlighting while viewing notes.
- 
 <img width="1883" height="957" alt="image_2025-07-27_00-37-07" src="https://github.com/user-attachments/assets/ec7bedae-079a-449c-ad98-67d5de540911" />


---

## Project Structure

```
dropkey-tui/
├── api
│   ├── client.go      # HTTP client for backend communication
│   └── models.go      # Data models for API responses
├── config
│   ├── config.go      # Configuration loading logic
│   └── session.go     # Session management
├── crypt
│   ├── cipher.go      # AES-GCM encryption/decryption logic
│   └── keys.go        # Ed25519 key handling
├── go.mod             # Go module dependencies
├── go.sum             # Dependency checksums
├── main.go            # Application entry point
├── README.md          # Project documentation
└── tui
    ├── model.go       # BubbleTea models for TUI state management
    ├── styles
    │   └── styles.go  # Lip Gloss styles for TUI rendering
    └── views
        ├── dashboard.go  # Main dashboard view
        ├── landing.go    # Landing page view
        ├── login.go      # Login view
        ├── paste_form.go # Form for paste interaction
        ├── paste_list.go # List of retrieved pastes
        ├── register.go   # Registration view
        └── search.go     # Search view for paste IDs
```

---

## Planned Features

- Edit pastes: Update existing notes directly in the TUI.

---

## Tech Stack

- **[Go](https://golang.org/)**: Core programming language for performance and simplicity.
- **[BubbleTea](https://github.com/charmbracelet/bubbletea)**: Framework for building terminal-based UIs.
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)**: Styling library for TUI components.
- **[Go Crypto](https://pkg.go.dev/crypto)**: Standard library for AES-GCM and Ed25519 operations.

---

## Contributing

Contributions are welcome! To get started:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m "Add your feature"`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

Please follow the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/).

---

## License

This project is licensed under the [MIT License].

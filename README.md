▄▄▄█████▓ ██▀███   ▒█████   ███▄ ▄███▓ ▄▄▄▄    ▒█████   ███▄    █ ▓█████ 
▓  ██▒ ▓▒▓██ ▒ ██▒▒██▒  ██▒▓██▒▀█▀ ██▒▓█████▄ ▒██▒  ██▒ ██ ▀█   █ ▓█   ▀ 
▒ ▓██░ ▒░▓██ ░▄█ ▒▒██░  ██▒▓██    ▓██░▒██▒ ▄██▒██░  ██▒▓██  ▀█ ██▒▒███   
░ ▓██▓ ░ ▒██▀▀█▄  ▒██   ██░▒██    ▒██ ▒██░█▀  ▒██   ██░▓██▒  ▐▌██▒▒▓█  ▄ 
  ▒██▒ ░ ░██▓ ▒██▒░ ████▓▒░▒██▒   ░██▒░▓█  ▀█▓░ ████▓▒░▒██░   ▓██░░▒████▒
  ▒ ░░   ░ ▒▓ ░▒▓░░ ▒░▒░▒░ ░ ▒░   ░  ░░▒▓███▀▒░ ▒░▒░▒░ ░ ▒░   ▒ ▒ ░░ ▒░ ░
    ░      ░▒ ░ ▒░  ░ ▒ ▒░ ░  ░      ░▒░▒   ░   ░ ▒ ▒░ ░ ░░   ░ ▒░ ░ ░  ░
  ░        ░░   ░ ░ ░ ░ ▒  ░      ░    ░    ░ ░ ░ ░ ▒     ░   ░ ░    ░   
            ░         ░ ░         ░    ░          ░ ░           ░    ░  ░
                                            ░                            

(c) 2025 Zayn Otley
https://github.com/intuitionamiga/Trombone

# Trombone

A Go-based command-line tool for detecting Rust-compiled binaries on Linux systems.

## Overview

Trombone uses multiple detection techniques to determine whether a given binary or shared library was compiled from Rust source code. It combines static analysis methods to provide a confidence score for its findings.

## Features

- Uses five distinct detection techniques
- Provides detailed evidence for each detection method
- Shows overall confidence rating based on collected evidence
- Works with both executables and shared libraries
- Simple command-line interface

## Detection Techniques

1. **Rust-specific Strings** - Searches for characteristic string patterns found in Rust binaries
2. **Symbol Name Analysis** - Examines symbol names for Rust's distinctive name mangling patterns
3. **ELF Header Inspection** - Analyses ELF headers for Rust-specific sections and metadata
4. **Memory Allocator Detection** - Identifies Rust's custom memory allocation functions
5. **Panic Handler Analysis** - Looks for Rust's panic handling code patterns

## Requirements

- Go 1.23 or higher
- Linux environment
- The following command-line tools:
  - `strings`
  - `nm`
  - `readelf`
  - `objdump`

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/trombone.git
cd trombone

# Build the tool
go build -o trombone

# Optional: Install to your PATH
sudo mv trombone /usr/local/bin/
```

## Usage

```bash
trombone /path/to/binary
```

### Example Output

```
✅ The binary '/usr/bin/ripgrep' appears to be compiled from Rust (confidence: very high)

Detailed analysis:
=================
✅ Rust-specific strings:
  • core::fmt::Formatter::pad_integral
  • core::result::Result
  • alloc::slice::merge_sort
  • thread 'main' panicked
  • ... (more findings omitted)

✅ Rust symbol names:
  • 000000000001f790 T _ZN4core6result13unwrap_failed17h8720a4769a244c85E
  • 0000000000018450 T _ZN4core3fmt9Formatter12pad_integral17hc895323d47ddec59E
  • 0000000000030830 T _ZN4core3fmt3num52_$LT$impl$u20$core..fmt..Debug$u20$for$u20$usize$GT$3fmt17hfd6f9c7751ca9f51E
  • ... (more findings omitted)

✅ ELF headers and sections:
  • 000000000000000000000000 0000000000000000  0                 rustc .debug_str
  • 000000000000000000000000 0000000000000000  0                 core::panicking::panic
  • ... (more findings omitted)

✅ Rust memory allocator:
  • 0000000000012f30 T __rust_alloc
  • 0000000000012f70 T __rust_dealloc
  • 0000000000012fa0 T __rust_realloc

✅ Rust panic handlers:
  • Found panic-related code:
    00000000000354b0 <_ZN4core9panicking5panic17hc5919d38115dffb6E>:
      354b0:       41 57                   push   %r15
      354b1:       41 56                   push   %r14
  • ... (more findings omitted)
```

## Limitations

- Heavily stripped binaries may yield false negatives
- Some detection techniques may fail on certain binaries
- The tool makes a best-effort determination based on available evidence
- Not designed for non-Linux platforms

## How It Works

Trombone runs a series of analyses on the target binary:

1. Extracts strings from the binary using the `strings` utility and looks for Rust-specific patterns
2. Uses `nm` to extract symbol names and identifies Rust's name mangling scheme
3. Examines ELF headers and sections with `readelf` for Rust-specific attributes
4. Searches for Rust's memory allocator function signatures
5. Uses `objdump` to look for panic handler code in the disassembly

Based on how many techniques yield positive results, it provides a confidence rating:

- 1 detection: Low confidence
- 2 detections: Medium confidence
- 3 detections: High confidence
- 4-5 detections: Very high confidence

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GPL v3 or higher License because the Rustaliban don't like the GPL - see the LICENSE file for details.

## Modern Slavery Statement
100% vibe coded from a single prompt in Claude Sonnet 3.7.

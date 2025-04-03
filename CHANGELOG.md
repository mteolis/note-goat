# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

---

## [1.0.1] - 2025-04-03
### Added
- Added `CHANGELOG.md` to track project changes and release history.

### Changed
- Renamed `ExtractSummary` function in gemini package to more general `ExtractAnswer`.

### Fixed
- Fix rate limit error exiting script. ([#1](https://github.com/mteolis/note-goat/issues/1)), ([PR #4](https://github.com/mteolis/note-goat/pull/4))
    - Implemented `WaitForRateLimit` function to handle Gemini API rate limit errors with retries.
    - Added exponential backoff logic for handling rate limit errors.
    - Added logging to notify users when retries are triggered due to rate limits.

## [1.0.0] - 2025-03-30
### Added
- Initial release of the project.
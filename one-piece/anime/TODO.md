Based on my analysis of the anime downloader code, I've identified several potential bugs and issues:

## Critical Bugs Found

1. **Incorrect Directory Creation Logic**:
   - In `main.go`, line ~28-30, the code checks if a file exists and then calls `mkdirAll(dirName)` when it doesn't exist
   - This is backwards logic - it should create directories before checking for file existence

2. **Race Condition in Download Limit**:
   - In `main.go`, line ~39, the `currentDownloads` counter is incremented before checking if the file exists
   - This means downloads are counted even when files already exist

3. **Incorrect Error Handling in Download Function**:
   - In `internal/anime/anime.go`, line ~53, the code attempts to download from an iframe URL but doesn't properly handle cases where the iframe source is not available

## Medium Severity Issues

1. **Hardcoded Destination Path in Makefile**:
   - In `makefile`, line ~20, the destination path is hardcoded to `/Volumes/NAS/Videos/One Piece/Anime`
   - This makes the service non-portable

2. **Inefficient Sleep Logic**:
   - In `main.go`, lines ~50-65, the sleep logic uses a loop that continuously checks time and sleeps for long durations
   - This is inefficient and could be simplified

3. **Resource Leak Potential**:
   - In `internal/anime/anime.go`, line ~70, the code creates a new collector but doesn't ensure proper cleanup in all error cases

## Suggested Fixes

1. **Fix directory creation logic** by properly creating directories before file existence checks
2. **Reorganize download counting** to avoid counting existing files
3. **Make the destination path configurable** instead of hardcoded
4. **Simplify sleep logic** to eliminate redundant loops

Would you like me to provide specific code fixes for any of these issues?

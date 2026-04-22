# 🏗️ Helix Architecture

## Repository Structure

```
.helix/
├── objects/      # Content-addressable storage
├── refs/
│   ├── heads/   # Branches
│   └── tags/    # Tags
├── HEAD         # Current branch pointer
└── config       # Configuration
```

## Object Types

### Blob
```
blob <size>\0<content>
```

### Tree
```
tree <size>\0
<mode> <name>\0<20-byte-sha>
```

### Commit
```
commit <size>\0
tree <sha>
parent <sha>
author <name> <email> <timestamp>

Message
```

---

*Building Git from first principles*

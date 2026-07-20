# Migration: docs → klamhq/facter

This file tracks what needs to happen to migrate the documentation to a dedicated `klamhq/facter` repository.

## What to do

```bash
# 1. Create the new repository on GitHub: klamhq/facter (public, no README)

# 2. Clone facter-oss and extract the docs history
git clone https://github.com/klamhq/facter-oss facter-oss-migrate
cd facter-oss-migrate
git filter-repo --path docs/ --path mkdocs.yml --path requirements.txt --path .github/workflows/docs.yml --path .github/workflows/release-doc.yml

# 3. Push to the new repo
git remote set-url origin https://github.com/klamhq/facter
git push origin main

# 4. Update mkdocs.yml in the new repo:
#    - repo_name: klamhq/facter
#    - repo_url: https://github.com/klamhq/facter
#    - site_url: https://klamhq.github.io/facter/

# 5. In facter-oss: remove docs/, mkdocs.yml, requirements.txt
#    Remove .github/workflows/docs.yml and release-doc.yml
#    Update all READMEs to point to https://klamhq.github.io/facter/

# 6. Update GitHub Pages source in klamhq/facter settings
```

## Files to move

| Source (facter-oss)                 | Destination (facter)                |
| ----------------------------------- | ----------------------------------- |
| `docs/`                             | `docs/`                             |
| `mkdocs.yml`                        | `mkdocs.yml`                        |
| `requirements.txt`                  | `requirements.txt`                  |
| `.github/workflows/docs.yml`        | `.github/workflows/docs.yml`        |
| `.github/workflows/release-doc.yml` | `.github/workflows/release-doc.yml` |

## Files to delete from facter-oss after migration

- `docs/`
- `mkdocs.yml`
- `requirements.txt`
- `.github/workflows/docs.yml`
- `.github/workflows/release-doc.yml`

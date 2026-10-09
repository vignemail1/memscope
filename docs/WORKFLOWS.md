# GitHub Actions Workflows Guide

Ce document explique comment utiliser les workflows GitHub Actions configurés pour memscope.

## 🚀 Workflows Disponibles

### 1. Build Workflow (`build.yml`)

**Déclenché automatiquement sur :**

- Push sur la branche `main`
- Pull Requests vers `main`

**Ce qu'il fait :**

- ✅ Compile pour Windows, Linux, macOS (Intel + ARM)
- ✅ Exécute tous les tests (unitaires + intégration)
- ✅ Génère les SBOM (CycloneDX + SPDX)
- ✅ Scan de sécurité avec Trivy
- ✅ Archive les artefacts (binaires, SBOM, checksums)

### 2. Release Workflow (`release.yml`)

**Déclenché sur :**

- Push de tags `v*` (exemple: `v1.0.0`)

**Ce qu'il fait :**

- 🔨 Build professionnel avec GoReleaser
- 📦 Packages multi-plateformes
- 📋 SBOM et rapports de sécurité
- 🚀 Création automatique de GitHub Release
- 📝 Changelog automatique

## 📋 Comment Utiliser

### Pour une Release

1. **Finaliser le code**

   ```bash
   # S'assurer que tout est prêt
   go test ./...
   go build ./cmd/memscope
   ```

2. **Créer et pousser un tag**

   ```bash
   # Version majeure (breaking changes)
   git tag v2.0.0

   # Version mineure (nouvelles fonctionnalités)
   git tag v1.1.0

   # Version patch (corrections de bugs)
   git tag v1.0.1

   # Pousser le tag
   git push origin v1.0.0
   ```

3. **Le workflow fait le reste !**
   - Build automatique multi-plateforme
   - Tests de validation
   - Création de la GitHub Release
   - Publication des binaires

### Pour les Tests Automatiques

Les tests s'exécutent automatiquement sur chaque :

- Push sur `main`
- Pull Request

Aucune action manuelle requise.

## 📦 Artefacts Générés

### Build Workflow

Les artefacts sont disponibles dans l'onglet "Actions" :

- `memscope-build-{sha}` : Tous les binaires et fichiers SBOM
- `memscope-sbom-{sha}` : SBOM uniquement (pour compliance)

### Release Workflow

Les releases incluent :

- **Binaires** : `memscope_Windows_x86_64.zip`, `memscope_Darwin_arm64.tar.gz`, etc.
- **SBOM** : `cyclonedx-sbom.json`, `spdx-sbom.json`, etc.
- **Sécurité** : `trivy-security-report.json`
- **Checksums** : `checksums.txt`
- **Métadonnées** : `release-info.json`

## 🔧 Configuration

### Variables d'Environnement

Les workflows utilisent des variables automatiques GitHub :

- `GITHUB_TOKEN` : Authentification (automatique)
- `GITHUB_REPOSITORY` : Nom du repo (automatique)

### Secrets Optionnels

Pour la signature (non requis) :

- `GPG_FINGERPRINT` : Empreinte de clé GPG pour signer les releases

### Personnalisation

#### Modifier les Plateformes de Build

Dans `.goreleaser.yml` :

```yaml
builds:
  - goos:
      - linux
      - windows
      - darwin
      # - freebsd  # Ajouter d'autres OS
    goarch:
      - amd64
      - arm64
      # - 386     # Ajouter d'autres architectures
```

#### Modifier les Tests CI

Dans `.github/workflows/build.yml` :

```yaml
- name: Run tests
  run: |
    go test -v ./...
    go test -v -short ./test/integration
    # go test -v -race ./...  # Ajouter race detection
```

## 🛠️ Dépannage

### Build Failures

**Problème** : Les tests échouent dans CI

```bash
# Tester localement d'abord
go test ./...
go test -short ./test/integration
```

**Problème** : Erreur de compilation GoReleaser

```bash
# Vérifier la config GoReleaser
goreleaser check

# Test en mode snapshot
goreleaser build --snapshot --clean
```

### Release Issues

**Problème** : Le tag ne déclenche pas de release

- Vérifier que le tag suit le format `v*` (exemple: `v1.0.0`)
- Vérifier les permissions du workflow dans Settings > Actions

**Problème** : Échec de publication

- Vérifier que `GITHUB_TOKEN` a les permissions `contents: write`
- Vérifier qu'il n'y a pas de conflit de nom de release

### SBOM Generation

**Problème** : Erreur de génération SBOM

```bash
# Vérifier les dépendances
go mod tidy
go mod verify

# Test local des outils SBOM
cyclonedx-gomod mod -json -output test-sbom.json
```

## 📚 Resources

### Documentation

- [GoReleaser](https://goreleaser.com/quick-start/)
- [GitHub Actions](https://docs.github.com/actions/quickstart)
- [CycloneDX](https://cyclonedx.org/docs/)
- [SPDX](https://spdx.dev/learn/overview/)

### Outils Utiles

```bash
# Installation locale des outils
go install github.com/goreleaser/goreleaser@latest
go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest

# Test local d'une release
goreleaser release --snapshot --clean --skip-publish
```

### Bonnes Pratiques

1. **Toujours tester localement** avant de créer un tag
2. **Utiliser le versioning sémantique** (semver.org)
3. **Écrire des messages de commit descriptifs** pour le changelog automatique
4. **Vérifier les SBOM** pour la compliance de sécurité
5. **Monitorer les workflows** pour détecter les problèmes rapidement

## 🎯 Exemples de Workflow

### Développement Normal

```bash
# Développement sur une feature branch
git checkout -b feature/new-analysis

# Développement et tests
go test ./...
git add .
git commit -m "feat: add new memory analysis feature"

# Push pour tests automatiques
git push origin feature/new-analysis

# Créer PR → tests automatiques s'exécutent
# Merger PR → build sur main s'exécute
```

### Release Process

```bash
# Sur main après merge
git checkout main
git pull origin main

# Vérifier que tout fonctionne
go test ./...
./memscope version

# Créer et pousser le tag de release
git tag v1.2.0
git push origin v1.2.0

# Le workflow release s'exécute automatiquement
# Release publiée sur GitHub avec tous les artefacts
```

### Hotfix Release

```bash
# Pour une correction urgente
git checkout main
git pull origin main

# Appliquer le fix
git add .
git commit -m "fix: critical security vulnerability"

# Tag de patch
git tag v1.2.1
git push origin v1.2.1

# Release automatique disponible en quelques minutes
```

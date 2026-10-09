# CI/CD Workflows Documentation

Ce document décrit les workflows GitHub Actions configurés pour le projet memscope.

## 📦 Build Workflow (build.yml)

**Déclenché sur :** Push sur `main` et Pull Requests

### Fonctionnalités

1. **Tests automatisés**
   - Tests unitaires complets
   - Tests d'intégration (mode court)
   
2. **Build multi-plateforme**
   - Windows (amd64)
   - Linux (amd64)
   - macOS (amd64 + arm64)

3. **Génération SBOM (Software Bill of Materials)**
   - **CycloneDX** : Format JSON et XML
   - **SPDX** : Format JSON et tag-value
   
4. **Sécurité**
   - Scan de sécurité avec Trivy
   - Upload des résultats vers GitHub Security

5. **Artefacts générés**
   - Binaires pour toutes les plateformes
   - Fichiers SBOM pour la compliance
   - Checksums SHA256
   - Métadonnées de package
   - Archives compressées

### Utilisation

```bash
# Le workflow se déclenche automatiquement
# Les artefacts sont disponibles dans l'onglet Actions > Build
```

### Artefacts disponibles

- `memscope-build-{sha}` : Tous les binaires et SBOM
- `memscope-sbom-{sha}` : SBOM uniquement (rétention 90 jours)

## 🚀 Release Workflow (release.yml)

**Déclenché sur :** Push de tags `v*` (ex: `v1.0.0`)

### Fonctionnalités

1. **Tests pré-release**
   - Validation complète avant publication
   
2. **GoReleaser**
   - Build automatisé multi-plateforme
   - Génération de checksums
   - Archives formatées par OS
   - Changelog automatique

3. **SBOM et sécurité**
   - SBOM CycloneDX et SPDX
   - Scan de sécurité Trivy
   - Métadonnées de release

4. **Publication GitHub Release**
   - Release notes formatées
   - Binaires attachés
   - SBOM inclus
   - Rapport de sécurité

### Utilisation

```bash
# Créer et pousser un tag
git tag v1.0.0
git push origin v1.0.0

# Le workflow crée automatiquement la release GitHub
```

### Contenu des releases

- **Binaires** : Windows, Linux, macOS (Intel + ARM)
- **Archives** : tar.gz (Unix), zip (Windows)  
- **SBOM** : CycloneDX JSON/XML, SPDX JSON/tag
- **Sécurité** : Rapport Trivy
- **Checksums** : SHA256 pour vérification

## 🔧 Configuration GoReleaser

Le fichier `.goreleaser.yml` configure :

### Builds
```yaml
builds:
  - main: ./cmd/memscope
    env: [CGO_ENABLED=0]
    goos: [linux, windows, darwin]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}
```

### Archives
- Nommage automatique basé sur l'OS et l'architecture
- Format tar.gz pour Unix, zip pour Windows
- Inclusion automatique de README, LICENSE, docs

### Changelog
- Génération automatique depuis les commits
- Groupement par type (feat, fix, perf)
- Filtrage des commits de maintenance

## 📋 SBOM (Software Bill of Materials)

### CycloneDX Format
```json
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
  "components": [...]
}
```

### SPDX Format  
```json
{
  "spdxVersion": "SPDX-2.3",
  "name": "memscope",
  "packages": [...]
}
```

### Utilisations
- **Compliance** : Traçabilité des dépendances
- **Sécurité** : Identification des vulnérabilités
- **Audit** : Validation des licences
- **Supply Chain** : Sécurisation de la chaîne logicielle

## 🔒 Sécurité

### Scans Trivy
- **Vulnérabilités** : Dépendances et code
- **Secrets** : Détection de clés exposées  
- **Misconfigurations** : Fichiers de config
- **Licences** : Validation de compatibilité

### Résultats
- Format SARIF pour GitHub Security
- Format JSON pour analyse externe
- Upload automatique vers Security tab

## 📊 Monitoring et Métriques

### GitHub Actions
- **Durée des builds** : ~3-5 minutes
- **Taille des artefacts** : ~50MB total
- **Tests coverage** : Validé automatiquement

### Optimisations
- Cache des modules Go
- Build parallèle des plateformes
- Artifacts avec rétention configurée

## 🎯 Bonnes Pratiques

### Versioning
```bash
# Versions sémantiques
v1.0.0    # Release majeure
v1.1.0    # Nouvelles fonctionnalités
v1.1.1    # Corrections de bugs
```

### Tags
```bash
# Créer une release
git tag -a v1.0.0 -m "Release v1.0.0 - Initial stable release"
git push origin v1.0.0

# Release candidate
git tag v1.1.0-rc1
git push origin v1.1.0-rc1
```

### Commits
```bash
# Messages structurés pour changelog automatique
feat: add new memory analysis feature
fix: resolve timing calculation bug  
perf: optimize SPD parsing performance
docs: update installation instructions
```

## 🚨 Troubleshooting

### Build failures
```bash
# Vérifier localement avant push
go test ./...
go build ./cmd/memscope

# Tester le workflow GoReleaser
goreleaser check
goreleaser build --snapshot --clean
```

### SBOM generation issues
```bash
# Vérifier les dépendances
go mod tidy
go mod verify

# Test SBOM generation local
cyclonedx-gomod mod -json -output sbom.json
```

### Release problems
```bash
# Vérifier les permissions
# - contents: write
# - packages: write (si Docker activé)

# Vérifier les secrets
# - GITHUB_TOKEN (automatique)
# - GPG_FINGERPRINT (pour signature, optionnel)
```

## 📚 Resources

- [GoReleaser Documentation](https://goreleaser.com/)
- [CycloneDX Specification](https://cyclonedx.org/)
- [SPDX Specification](https://spdx.dev/)
- [GitHub Actions Documentation](https://docs.github.com/actions)
- [Trivy Security Scanner](https://trivy.dev/)